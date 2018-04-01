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

	ClientDefaultTimeout = time.Second / 100
	ClientPingPeriod     = time.Second / 10
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
	IsFull      bool
	IsCoherent  bool
	RdyServData bool
	Wanted      string
}
