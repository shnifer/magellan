package main

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
	"os"
	"os/signal"
	"time"
)

var server *network.Server

const (
	storagePath = "xstore"
)

func main() {
	if DEFVAL.DoProf {
		commons.StartProfile(roleName)
		defer commons.StopProfile(roleName)
	}

	diskOpts := diskv.Options{
		BasePath:     storagePath,
		CacheSizeMax: 1024 * 1024,
	}
	disk := storage.New(DEFVAL.NodeName, diskOpts)

	roomServ := newRoomServer(disk)

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

	go daemonUpdateSubscribes(roomServ, server, DEFVAL.SubscribeUpdatePeriod)

	//waiting for enter to stop server
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}
