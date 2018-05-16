package main

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
	"os"
	"os/signal"
	"time"
)

var server *network.Server

func main() {
	if DEFVAL.DoProf {
		commons.StartProfile(roleName)
		defer commons.StopProfile(roleName)
	}

	roomServ := newRoomServer()

	startState := commons.State{
		StateID: commons.STATE_login,
	}

	opts := network.ServerOpts{
		Addr:             DEFVAL.Port,
		RoomUpdatePeriod: time.Duration(DEFVAL.RoomUpdatePeriod) * time.Millisecond,
		LastSeenTimeout:  time.Duration(DEFVAL.LastSeenTimeout) * time.Millisecond,
		RoomServ:         roomServ,
		StartState:       startState.Encode(),
		NeededRoles:      DEFVAL.NeededRoles,
	}

	server = network.NewServer(opts)
	defer server.Close()

	//waiting for enter to stop server
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}
