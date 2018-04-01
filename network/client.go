package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type ClientOpts struct {
	//default ClientDefaultTimeout
	Timeout time.Duration

	Addr string

	Room, Role string

	//Specific self disconnect (server lost). may be needed later , but in general use Pause
	OnReconnect  func()
	OnDisconnect func()

	//Any reason's pause of game process (disconnect self, disconnect other, loading new state self or other
	//Specific reason may be getted by PauseReason()
	OnPause   func()
	OnUnpause func()

	OnCommonSend func() []byte
	OnCommonRecv func([]byte)

	OnStateChanged func(wanted string)

	//async, must close result chan then done
	OnGetStateData func([]byte) chan struct{}
}

type Client struct {
	mu sync.RWMutex

	httpCli http.Client
	opts    ClientOpts

	started bool
	pingLostCounter int

	//for hooks
	pingLost bool
	onPause  bool

	//copy of last ping state
	isFull     bool
	isCoherent bool

	//for state machine
	curState  string
	wantState string
}

func NewClient(opts ClientOpts) (*Client, error) {
	if opts.Timeout == 0 {
		opts.Timeout = ClientDefaultTimeout
	}

	httpCli := http.Client{
		Timeout: opts.Timeout,
	}

	res := &Client{
		httpCli: httpCli,
		opts:    opts,

		//starts from unconnected states,
		//so opt.OnReconnect and opt.OnUnpause will be called on first connection
		pingLost: true,
		onPause:  true,
	}

	return res, nil
}

func (c *Client) setPingLost(lost bool) {
	if lost && !c.pingLost {
		if c.opts.OnDisconnect != nil {
			c.opts.OnDisconnect()
		}
	}
	if !lost && c.pingLost {
		if c.opts.OnReconnect != nil {
			c.opts.OnReconnect()
		}
	}
	c.pingLost = lost
}

var pauseC int
func (c *Client) setOnPause(pause bool) {
	if pause && !c.onPause {
		if c.opts.OnPause != nil {
			c.opts.OnPause()
		}
	}
	if !pause && c.onPause {
		if c.opts.OnUnpause != nil {
			c.opts.OnUnpause()
		}
	}
	c.onPause = pause

	if pause {
		pauseC++
		log.Println("pause tick", pauseC)
	} else {
		pauseC = 0
	}
}

func doPingReq(c *Client) (RoomState, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.doReq(GET, pingPattern, nil)
	if err != nil {
		//Connection is not good if ClientLostPingsNumber in row
		if !c.pingLost {
			c.pingLostCounter++
			log.Println("pingLostCounter", c.pingLostCounter)
			if c.pingLostCounter>=ClientLostPingsNumber {
				c.pingLostCounter = 0
				c.setPingLost(true)
				c.setOnPause(true)
			}
		}

		urlErr, ok := err.(*url.Error)
		if !ok {
			log.Println("network.clientPing: Strange non-URL error client ping", err)
		} else if !urlErr.Timeout() {
			log.Println("network.clientPing: Strange non-timeout error client ping", err)
		}
		return RoomState{}, err
	}

	c.setPingLost(false)

	var pingResp RoomState
	err = json.Unmarshal(resp, &pingResp)
	if err != nil {
		return RoomState{}, err
	}

	c.isFull = pingResp.IsFull
	c.isCoherent = pingResp.IsCoherent

	//check for pause
	needPause := c.pingLost || !c.isFull || !c.isCoherent
	c.setOnPause(needPause)

	return pingResp, nil
}

func checkWantedState(c *Client, roomState RoomState) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//state changed
	wanted := roomState.Wanted
	if wanted != c.wantState {
		c.wantState = wanted
		//aware client about new state
		if c.opts.OnStateChanged != nil {
			c.opts.OnStateChanged(wanted)
		}
	}

	if c.wantState != c.curState {
		//rdy to grab new state Data
		if roomState.RdyServData {
			resp, err := c.doReq(GET, statePattern, nil)
			if err != nil {
				//weird, but will try next ping circle
				log.Println("can't get new Serv Data")
				return
			}

			//After successfully got and passed new StateData change cur state
			if c.opts.OnGetStateData == nil {
				//set wanted state now
				c.curState = c.wantState
				//Get commonState after reading StateData
				doCommonReq(c, true)
			} else {
				//run hook and wait for done chan close
				stateDataDone:=c.opts.OnGetStateData(resp)
				go func() {
					<-stateDataDone
					c.mu.Lock()
					c.curState = c.wantState
					//Get commonState after reading StateData
					doCommonReq(c, true)
					c.mu.Unlock()
				}()
			}
		}
	}
}

func doCommonReq(c *Client, onlyGet bool) {
	method := GET
	var sentBuf io.Reader

	if !onlyGet {
		var sentData []byte
		if c.opts.OnCommonSend != nil {
			sentData = c.opts.OnCommonSend()
		}

		if sentData != nil && len(sentData) > 0 {
			method = POST
			sentBuf = bytes.NewBuffer(sentData)
		}
	}

	resp, err := c.doReq(method, roomPattern, sentBuf)
	if err != nil {
		log.Println("CANT SEND common room data request", err)
		return
	}
	if c.opts.OnCommonRecv != nil {
		c.opts.OnCommonRecv(resp)
	}
}
func clientPing(c *Client) {
	tick := time.Tick(ClientPingPeriod)
	for {
		<-tick

		//do Ping to check online and State
		RoomState, err := doPingReq(c)
		if err != nil {
			//log.Println(err)
			continue
		}
		checkWantedState(c, RoomState)

		//Maybe it is better to run GetCommonData loop as other routine but YAGNI
		if !c.onPause {
			doCommonReq(c, false)
		}
	}
}

func (c *Client) doReq(method, path string, reqBody io.Reader) (respBody []byte, er error) {
	req, err := http.NewRequest(method, c.opts.Addr+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set(roomAttr, c.opts.Room)
	req.Header.Set(roleAttr, c.opts.Role)
	req.Header.Set(stateAttr, c.curState)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.Header.Get("error") == "1" {
		errStr := string(buf)
		log.Println(errStr)
		return nil, errors.New(errStr)
	}

	return buf, nil
}

type PauseReason struct {
	PingLost   bool
	IsFull     bool
	IsCoherent bool
	CurState   string
	WantState  string
}

func (c *Client) PauseReason() PauseReason {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return PauseReason{
		PingLost:   c.pingLost,
		IsFull:     c.isFull,
		IsCoherent: c.isCoherent,
		CurState:   c.curState,
		WantState:  c.wantState,
	}
}

func (c *Client) RequestNewState(wanted string) {
	buf := strings.NewReader(wanted)
	c.doReq(POST, statePattern, buf)
}

func (c *Client) Start() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return
	}

	c.started = true
	go clientPing(c)
}