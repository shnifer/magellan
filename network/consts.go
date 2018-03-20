package network

import (
	"time"
	"net/http"
)

const (
	GET = http.MethodGet
	POST = http.MethodPost

	roomAttr = "room"
	roleAttr = "attr"

	roomPattern = "/room/"
	pingPattern = "/ping/"

	ClientDefaultTimeout = time.Second/10
	ClientPingPeriod = time.Second/10

	ServerRoomUpdatePeriod = time.Second/10
	ServerLastSeenTimeout = 3*ServerRoomUpdatePeriod

	MSG_FullRoom = "FullRoom"
	MSG_HalfRoom = "HalfRoom"
)

//network.Server - network.Client ping response
type pingMSG struct {
	isRoomFull bool

}