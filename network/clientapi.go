package network

import (
	"log"
	"strings"
	"time"
)

type ClientOpts struct {
	//default ClientDefaultTimeout
	Timeout time.Duration

	//default ClientDefaultPingPeriod
	PingPeriod time.Duration

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
	OnCommonRecv func(data []byte, readOwnPart bool)

	OnStateChanged func(wanted string)

	//async, must close result chan then done
	OnGetStateData func([]byte)

	OnCommand func(command string)
}

type PauseReason struct {
	PingLost   bool
	IsFull     bool
	IsCoherent bool
	CurState   string
	WantState  string
}

func (pr PauseReason) String() string {
	msg := make([]string, 0)
	add := func(s string) {
		msg = append(msg, s)
	}

	if pr.PingLost {
		add("ping lost")
	}
	if pr.IsFull {
		add("full")
	} else {
		add("not full")
	}
	if pr.IsCoherent {
		add("coherent")
	} else {
		add("non-coherent")
	}
	if pr.CurState == pr.WantState {
		add("state " + pr.CurState)
	} else {
		add("current state " + pr.CurState)
		add("wanted state " + pr.WantState)
	}
	return strings.Join(msg, ", ")
}

func (c *Client) recalcPauseReason() {
	c.mu.RLock()
	pr := PauseReason{
		PingLost:   c.pingLost,
		IsFull:     c.isFull,
		IsCoherent: c.isCoherent,
		CurState:   c.curState,
		WantState:  c.wantState,
	}
	c.mu.RUnlock()

	c.prmu.Lock()
	c.pr = pr
	c.prmu.Unlock()
}

func (c *Client) PauseReason() PauseReason {
	c.prmu.RLock()
	defer c.prmu.RUnlock()

	return c.pr
}

func (c *Client) RequestNewState(wanted string) {
	if c.wantState != c.curState {
		log.Println("client is already changing state")
	}
	//_, err := c.doReq(POST, statePattern, []byte(wanted))

	c.sendCommand(COMMAND_REQUESTSTATE, wanted)
}

func (c *Client) SendRoomBroadcast(command string) {
	c.sendCommand(COMMAND_ROOMBROADCAST, command)
}
func (c *Client) SendRequest(command string) {
	c.sendCommand(COMMAND_CLIENTREQUEST, command)
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
