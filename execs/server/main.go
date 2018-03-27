package main

import (
	"github.com/Shnifer/magellan/network"
	"time"
)

func main() {
	roomServ := newRoomServer()

	opts := network.ServerOpts{
		Addr:     DEFVAL.Port,
		RoomServ: roomServ,
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
