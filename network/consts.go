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

	ClientDefaultTimeout = time.Second / 10
	ClientPingPeriod     = time.Second / 10

	ServerRoomUpdatePeriod = time.Second / 10
	ServerLastSeenTimeout  = 3 * ServerRoomUpdatePeriod
	)

//both room
//network.Server - network.Client ping response
type roomState struct {
	isFull      bool
	isCoherent  bool
	rdyServData bool
	wanted      string
}
