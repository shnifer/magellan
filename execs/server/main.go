package main

import (
	"github.com/Shnifer/magellan/network"
	"time"
	"github.com/Shnifer/magellan/commons"
)

func main() {

	roomServ := newRoomServer()

	startState:=commons.State{
		Special: commons.STATE_login,
	}

	opts := network.ServerOpts{
		Addr:     DEFVAL.Port,
		RoomServ: roomServ,
		StartState: startState.Encode(),
	}

	server, err := network.NewServer(opts)
	if err != nil {
		panic(err)
	}
	defer server.Close()

	//waiting for enter to stop server
	for {
		time.Sleep(time.Second)
	}
}
