package network

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	mu sync.RWMutex

	httpCli http.Client
	opts    ClientOpts

	started         bool
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

	//do we need to RECEIVE our part of common
	isMyPartActual bool

	//commands to Send
	//own mutex
	scmu              sync.Mutex
	sendCommandsBaseN int
	sendCommands      []string

	lastReceivedCommandN int

	//mutex for PauseReason only
	prmu sync.RWMutex
	pr   PauseReason
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
		pingLost:     true,
		onPause:      true,
		sendCommands: make([]string, 0),
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
}

func doPingReq(c *Client) (PingResp, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.doReq(GET, pingPattern, nil)
	if err != nil {
		//Connection is not good if ClientLostPingsNumber in row
		if !c.pingLost {
			c.pingLostCounter++
			if c.pingLostCounter >= ClientLostPingsNumber {
				c.pingLostCounter = 0
				c.setPingLost(true)
				c.setOnPause(true)
			}
		}

		urlErr, ok := err.(*url.Error)
		if !ok {
			log.Println("network.doPingReq: Strange non-URL error client ping", err)
		} else if !urlErr.Timeout() {
			log.Println("network.doPingReq: Strange non-timeout error client ping", err)
		}
		return PingResp{}, err
	}

	c.setPingLost(false)

	var pingResp PingResp
	err = json.Unmarshal(resp, &pingResp)
	if err != nil {
		return PingResp{}, err
	}

	c.isFull = pingResp.Room.IsFull
	c.isCoherent = pingResp.Room.IsCoherent

	//check for pause
	needPause := c.pingLost || !c.isFull || !c.isCoherent
	c.setOnPause(needPause)

	return pingResp, nil
}

func checkWantedState(c *Client, pingResp PingResp) {
	c.mu.Lock()
	defer c.mu.Unlock()

	//state changed
	wanted := pingResp.Room.Wanted
	if wanted != c.wantState {
		//Drop commands
		c.sendCommandsBaseN += len(c.sendCommands)
		c.sendCommands = c.sendCommands[:0]

		c.wantState = wanted
		c.isMyPartActual = false
		//aware client about new state
		if c.opts.OnStateChanged != nil {
			c.opts.OnStateChanged(wanted)
		}
	}

	if c.wantState != c.curState {
		//rdy to grab new state Data
		if pingResp.Room.RdyServData {
			resp, err := c.doReq(GET, statePattern, nil)
			if err != nil {
				//weird, but will try next ping circle
				log.Println("can't get new ServData", err)
				return
			}

			//After successfully got and passed new StateData change cur state
			if c.opts.OnGetStateData == nil {
				//set wanted state now
				c.curState = c.wantState
			} else {
				//run hook and wait for done chan close
				go func() {
					c.opts.OnGetStateData(resp)
					c.mu.Lock()
					c.curState = c.wantState
					c.mu.Unlock()
				}()
			}
		}
	}
}

func clientCheckSentCommandsReceived(c *Client, pingResp PingResp) {
	c.scmu.Lock()
	defer c.scmu.Unlock()

	//fresh just started client. Continue from servers last position
	if c.sendCommandsBaseN == 0 {
		c.sendCommandsBaseN = pingResp.LastCommandReceived + 1
		return
	}
	delta := pingResp.LastCommandReceived - c.sendCommandsBaseN + 1
	if delta < 0 {
		log.Println("strange LastCommandReceived<sendCommandsBaseN",
			pingResp.LastCommandReceived, c.sendCommandsBaseN)
	}
	if delta == 0 {
		return
	}
	if delta > len(c.sendCommands) {
		log.Println("strange delta>sendCommandsBaseN+len",
			delta, len(c.sendCommands))
		delta = len(c.sendCommands)
	}
	c.sendCommands = c.sendCommands[delta:]
	c.sendCommandsBaseN += delta
}

func clientReceiveCommands(c *Client, resp CommonResp) {
	if c.opts.OnCommand == nil {
		return
	}

	for i, command := range resp.Commands {
		commandN := resp.CommandsBaseN + i
		if commandN <= c.lastReceivedCommandN {
			continue
		}
		c.opts.OnCommand(command)
	}

	c.lastReceivedCommandN = resp.CommandsBaseN + len(resp.Commands) - 1
}

func doCommonReq(c *Client) {
	c.scmu.Lock()
	defer c.scmu.Unlock()

	var req CommonReq
	var sentData []byte

	if c.isMyPartActual {
		if c.opts.OnCommonSend != nil {
			sentData = c.opts.OnCommonSend()
		}

		if sentData != nil && len(sentData) > 0 {
			req.DataSent = true
			req.Data = string(sentData)
		}
	}

	req.Commands = c.sendCommands
	req.CommandsBaseN = c.sendCommandsBaseN
	req.LastReceivedCommandN = c.lastReceivedCommandN

	buf, err := json.Marshal(req)
	if err != nil {
		log.Println("can't marshal commonReq")
		return
	}

	respBytes, err := c.doReq(POST, roomPattern, buf)
	if err != nil {
		log.Println("CANT SEND common room data request", err)
		return
	}
	var resp CommonResp
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		log.Println("Can't unmarshal common resp")
	}

	if c.isMyPartActual {
		clientReceiveCommands(c, resp)
	}

	if c.opts.OnCommonRecv != nil {
		c.opts.OnCommonRecv([]byte(resp.Data), !c.isMyPartActual)
	}

	c.isMyPartActual = true

}

func clientPing(c *Client) {
	tick := time.Tick(ClientPingPeriod)
	for {
		<-tick

		//do Ping to check online and State
		pingResp, err := doPingReq(c)
		if err != nil {
			c.recalcPauseReason()
			continue
		}
		checkWantedState(c, pingResp)
		clientCheckSentCommandsReceived(c, pingResp)

		if !c.onPause {
			doCommonReq(c)
		}
		c.recalcPauseReason()
	}
}

func (c *Client) doReq(method, path string, reqBody []byte) (respBody []byte, er error) {
	bodyBuf := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(method, c.opts.Addr+path, bodyBuf)
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

func (c *Client) sendCommand(prefix string, command string) {
	if len(prefix) != 1 {
		panic("sendCommand wrong prefix!")
	}

	c.scmu.Lock()
	c.sendCommands = append(c.sendCommands, prefix+command)
	c.scmu.Unlock()
}
