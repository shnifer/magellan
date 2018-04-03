package main

import (
	"github.com/Shnifer/magellan/network"
	"log"
)

var Client *network.Client

func initClient() {
	opts := network.ClientOpts{
		Addr:           DEFVAL.Port,
		Room:           DEFVAL.Room,
		Role:           DEFVAL.Role,
		OnReconnect:    recon,
		OnDisconnect:   discon,
		OnPause:        pause,
		OnUnpause:      unpause,
		OnCommonRecv:   Data.commonRecv,
		OnCommonSend:   Data.commonSend,
		OnStateChanged: stateChanged,
		OnGetStateData: Data.getStateData,
	}

	var err error
	Client, err = network.NewClient(opts)
	if err != nil {
		panic(err)
	}
}

func discon() {
	log.Println("lost connect")
}

func recon() {
	log.Println("recon!")
}