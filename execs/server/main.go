package main

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/network"
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
	"os"
	"os/signal"
	"time"
)

var server *network.Server

//todo: gamemaster control
func main() {
	log.Start(time.Duration(DEFVAL.LogLogTimeoutMs)*time.Millisecond,
		time.Duration(DEFVAL.LogRetryMinMs)*time.Millisecond,
		time.Duration(DEFVAL.LogRetryMaxMs)*time.Millisecond,
		DEFVAL.LogIP, DEFVAL.LogHostName)

	if DEFVAL.DoProf {
		commons.StartProfile(roleName,DEFVAL.DebugPort)
		defer commons.StopProfile(roleName)
	}

	logDiskOpts := diskv.Options{
		BasePath:     DEFVAL.LocalLogPath,
		CacheSizeMax: 1024,
	}
	logDisk := storage.New(DEFVAL.NodeName, logDiskOpts, 0)

	if DEFVAL.LogExchPort != "" && DEFVAL.LogExchPeriodMs > 0 {
		storage.RunExchanger(logDisk, DEFVAL.LogExchPort, DEFVAL.LogExchAddrs, DEFVAL.LogExchPeriodMs)
	}
	log.SetStorage(logDisk)

	diskOpts := diskv.Options{
		BasePath:     DEFVAL.StoragePath,
		CacheSizeMax: 1024 * 1024,
	}
	disk := storage.New(DEFVAL.NodeName, diskOpts, DEFVAL.DiskRefreshPeriod)
	if DEFVAL.GameExchPort != "" && DEFVAL.GameExchPeriodMs > 0 {
		storage.RunExchanger(disk, DEFVAL.GameExchPort, DEFVAL.GameExchAddrs, DEFVAL.GameExchPeriodMs)
	}

	restoreOpts := diskv.Options{
		BasePath:     DEFVAL.RestorePath,
		CacheSizeMax: 1024 * 1024,
	}
	diskRestore := diskv.New(restoreOpts)

	roomServ := newRoomServer(disk, diskRestore)

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
	go daemonUpdateOtherShips(roomServ, DEFVAL.OtherShipsUpdatePeriod)

	//waiting for enter to stop server
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}
