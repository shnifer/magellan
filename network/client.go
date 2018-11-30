package network

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/wrnt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	mu sync.RWMutex

	httpCli     http.Client
	httpCliLong http.Client
	opts        ClientOpts

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

	//do we need to SEND our part of common
	isMyPartActual       bool
	recvGoroutineStarted bool

	//commands, mutex inside
	send *wrnt.Send
	recv *wrnt.Recv

	//mutex for PauseReason only
	prmu sync.RWMutex
	pr   PauseReason
}

func NewClient(opts ClientOpts) (*Client, error) {
	defer LogFunc("network.NewClient")()

	if opts.Timeout == 0 {
		opts.Timeout = ClientDefaultTimeout
	}

	if opts.PingPeriod == 0 {
		opts.PingPeriod = ClientDefaultPingPeriod
	}

	httpCli := http.Client{
		Timeout: opts.Timeout,
	}
	httpCliLong := http.Client{
		Timeout: ClientLargeTimeout,
	}

	res := &Client{
		httpCli:     httpCli,
		httpCliLong: httpCliLong,
		opts:        opts,

		//starts from unconnected states,
		//so opt.OnReconnect and opt.OnUnpause will be called on first connection
		pingLost: true,
		onPause:  true,
		send:     wrnt.NewSend(),
		recv:     wrnt.NewRecv(),
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
	defer LogFunc("network.doPingReq")()

	c.mu.Lock()
	defer c.mu.Unlock()

	resp, err := c.doReq(GET, pingPattern, nil, false)
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
			Log(LVL_WARN, "network.doPingReq: Strange non-URL error client ping", err)
		} else if !urlErr.Timeout() {
			Log(LVL_WARN, "network.doPingReq: Strange non-timeout error client ping", err)
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
	defer LogFunc("network.checkWantedState")()

	c.mu.Lock()
	defer c.mu.Unlock()

	//state changed
	wanted := pingResp.Room.Wanted
	if wanted != c.wantState {
		Log(LVL_DEBUG, "wanted != c.wantState:", wanted, "!=", c.wantState)
		c.wantState = wanted
		c.isMyPartActual = false
		//aware client about new state
		if c.opts.OnStateChanged != nil {
			c.opts.OnStateChanged(wanted)
		}
	}

	//run GetStateData goroutine if needed
	if c.curState != c.wantState && !c.recvGoroutineStarted && pingResp.Room.RdyServData {
		Log(LVL_DEBUG, "started go getStateData(c)")
		c.recvGoroutineStarted = true
		go getStateData(c)
	}
}

//runned in goroutine
func getStateData(c *Client) {
	defer LogFunc("network.getStateData")()

	resp, err := c.doReq(GET, statePattern, nil, true)

	if err != nil {
		//weird, but will try next ping circle
		Log(LVL_WARN, "can't get new ServData", err)
		c.mu.Lock()
		c.recvGoroutineStarted = false
		c.mu.Unlock()
		return
	}

	if c.opts.OnGetStateData == nil {
		//set wanted state now
		c.mu.Lock()
		c.curState = c.wantState
		c.isMyPartActual = true
		c.recvGoroutineStarted = false
		c.mu.Unlock()
		return
	}

	//run hook and wait for done chan close
	buf := bytes.NewBuffer(resp)
	dec := gob.NewDecoder(buf)
	var DataResp StateDataResp
	dec.Decode(&DataResp)
	c.opts.OnGetStateData(DataResp.StateData)
	if c.opts.OnCommonRecv != nil {
		c.opts.OnCommonRecv(DataResp.StartCommon, true)
	}
	c.mu.Lock()
	c.curState = c.wantState
	c.isMyPartActual = true
	c.recvGoroutineStarted = false
	c.mu.Unlock()
}

func clientReceiveCommands(c *Client, resp CommonResp) {
	defer LogFunc("network.clientReceiveCommands")()

	if c.opts.OnCommand == nil {
		return
	}

	commands := c.recv.Unpack(resp.Message)

	for _, command := range commands {
		c.opts.OnCommand(command)
	}
}

func doCommonReq(c *Client) {
	defer LogFunc("network.doCommonReq")()

	var req CommonReq
	var sentData []byte

	if c.opts.OnCommonSend != nil {
		c.mu.RLock()
		isMyPartActual := c.isMyPartActual
		c.mu.RUnlock()

		if isMyPartActual {
			sentData = c.opts.OnCommonSend()
		}
	}

	if sentData != nil && len(sentData) > 0 {
		req.DataSent = true
		req.Data = string(sentData)
	}
	//	}
	message, err := c.send.Pack()
	if err == nil {
		req.Message = message
	}
	req.ClientConfirmN = c.recv.LastRecv()

	buf, err := json.Marshal(req)
	if err != nil {
		Log(LVL_ERROR, "can't marshal commonReq")
		return
	}

	respBytes, err := c.doReq(POST, roomPattern, buf, false)
	if err != nil {
		Log(LVL_WARN, "CANT SEND common room data request", err)
		return
	}
	var resp CommonResp
	err = json.Unmarshal(respBytes, &resp)
	if err != nil {
		Log(LVL_ERROR, "Can't unmarshal common resp")
	}

	clientReceiveCommands(c, resp)

	if c.opts.OnCommonRecv != nil {
		//		c.opts.OnCommonRecv([]byte(resp.Data), !c.isMyPartActual)
		c.opts.OnCommonRecv([]byte(resp.Data), false)
	}

	//c.isMyPartActual = true

}

func clientPing(c *Client) {
	defer LogFunc("network.clientPing")()

	tick := time.Tick(c.opts.PingPeriod)
	for {
		<-tick

		//do Ping to check online and State
		pingResp, err := doPingReq(c)
		if err != nil {
			c.recalcPauseReason()
			continue
		}
		checkWantedState(c, pingResp)
		c.send.Confirm(pingResp.ServerConfirmN)

		if !c.onPause {
			doCommonReq(c)
		}
		c.recalcPauseReason()
	}
}

func (c *Client) doReq(method, path string, reqBody []byte, largeTimeout bool) (respBody []byte, er error) {

	bodyBuf := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(method, c.opts.Addr+path, bodyBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set(roomAttr, c.opts.Room)
	req.Header.Set(roleAttr, c.opts.Role)
	req.Header.Set(stateAttr, c.curState)

	var resp *http.Response
	if largeTimeout {
		resp, err = c.httpCliLong.Do(req)
	} else {
		resp, err = c.httpCli.Do(req)
	}
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
		Log(LVL_ERROR, errStr)
		return nil, errors.New(errStr)
	}

	return buf, nil
}

func (c *Client) sendCommand(prefix string, command string) {
	if len(prefix) != 1 {
		panic("sendCommand wrong prefix!")
	}

	c.send.AddItems(prefix + command)
}
