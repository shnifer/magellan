package main

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
	"time"
)

var server *network.Server

func main() {

	roomServ := newRoomServer()

	startState := commons.State{
		StateID: commons.STATE_login,
	}

	opts := network.ServerOpts{
		Addr:        DEFVAL.Port,
		RoomServ:    roomServ,
		StartState:  startState.Encode(),
		NeededRoles: DEFVAL.NeededRoles,
	}

	server = network.NewServer(opts)
	defer server.Close()

	//waiting for enter to stop server
	for {
		time.Sleep(time.Second)
	}
}
