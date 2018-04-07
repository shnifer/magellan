package network

import (
	"net/http"
	"time"
)

const (
	GET  = http.MethodGet
	POST = http.MethodPost

	roomAttr  = "room"
	roleAttr  = "role"
	stateAttr = "state"

	roomPattern  = "/room/"
	pingPattern  = "/ping/"
	statePattern = "/state/"

	ClientDefaultTimeout  = time.Second / 100
	ClientPingPeriod      = time.Second / 10
	ClientLostPingsNumber = 3

	ServerRoomUpdatePeriod = time.Second / 10
	ServerLastSeenTimeout  = 3 * ServerRoomUpdatePeriod
)

//network.Server - network.Client ping response
//IsFull - all needed roles are online
//IsCoherent - state of room is confirmed by all roles
//RdyServData - Server downloaded data for new state (needed while not coherent)
//Wanted - Wanted state
type RoomState struct {
	//room specific
	IsFull      bool
	IsCoherent  bool
	RdyServData bool
	Wanted      string
}

type PingResp struct {
	Room                RoomState
	LastCommandReceived int
}

type CommonReq struct {
	DataSent bool
	Data     string

	//Command string first rune is service flag
	//Client->Server
	CommandsBaseN int
	Commands      []string

	//for Server<-Client
	LastReceivedCommandN int
}

type CommonResp struct {
	Data string

	//Server->Client
	CommandsBaseN int
	Commands      []string
}

const (
	COMMAND_CLIENT       = "C"
	COMMAND_REQUESTSTATE = "S"
)
