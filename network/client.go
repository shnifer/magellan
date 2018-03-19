package network

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"bytes"
)

type ClientOpts struct {
	//default ClientDefaultTimeout
	Timeout time.Duration

	Addr string

	Room, Role string

	OnReconnect  func()
	OnDisconnect func()
	OnPause      func()
	OnUnpause    func()

	OnCommonSend func() []byte
	OnCommonRecv func([]byte)
}

type Client struct {
	httpCli  http.Client
	opts     ClientOpts
	pingLost bool
	roomFull bool
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
		roomFull: false,
	}

	go clientPing(res)

	return res, nil
}

func (c *Client) procLostPing() {
	if !c.pingLost {
		if c.opts.OnDisconnect != nil {
			c.opts.OnDisconnect()
		}
	}
	c.pingLost = true
	c.roomFull = false
}

func (c *Client) procGoodPing() {
	if c.pingLost {
		if c.opts.OnReconnect != nil {
			c.opts.OnReconnect()
		}
	}
	c.pingLost = false
}

func (c *Client) procFullRoom() {
	c.procGoodPing()
	if !c.roomFull {
		if c.opts.OnUnpause != nil {
			c.opts.OnUnpause()
		}
	}
	c.roomFull = true
}

func (c *Client) procHalfRoom() {
	c.procGoodPing()
	if c.roomFull {
		if c.opts.OnPause != nil {
			c.opts.OnPause()
		}
	}
	c.roomFull = false
}

func (c *Client) readyToGetCommon()bool{
	return c.roomFull && !c.pingLost
}

func clientPing(c *Client) {
	tick := time.Tick(ClientPingPeriod)
	for {
		<-tick

		//do Ping to check online and State
		{
			resp, err := c.doReq(GET, pingPattern, nil)
			if err != nil {
				//Anyway connection is not good!
				c.procLostPing()

				urlErr, ok := err.(*url.Error)
				if !ok {
					log.Println("network.clientPing: Strange non-URL error client ping", err)
				} else if !urlErr.Timeout() {
					log.Println("network.clientPing: Strange non-timeout error client ping", err)
				}
				continue
			}

			switch string(resp) {
			case MSG_FullRoom:
				c.procFullRoom()
			case MSG_HalfRoom:
				c.procHalfRoom()
			default:
				log.Println("network.clientPing: strange ping resp!", string(resp))
			}
		}

		//Maybe it is better to run GetCommonData loop as other routine but YAGNI
		if c.readyToGetCommon(){
			var sentData []byte
			if c.opts.OnCommonSend!=nil {
				sentData = c.opts.OnCommonSend()
			}

			method:=GET
			var sentBuf io.Reader
			if sentData != nil || len(sentData)==0{
				method =POST
				sentBuf=bytes.NewBuffer(sentData)
			}

			resp, err:= c.doReq(method, roomPattern, sentBuf)
			if err!=nil{
				log.Println("CANT SEND common room data request")
				continue
			}
			if c.opts.OnCommonRecv!=nil{
				c.opts.OnCommonRecv(resp)
			}
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

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
