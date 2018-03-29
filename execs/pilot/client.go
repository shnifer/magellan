package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
	"log"
)

var Client *network.Client

func startClient() {
	opts := network.ClientOpts{
		Addr:           DEFVAL.Port,
		Room:           DEFVAL.Room,
		Role:           ROLE_Pilot,
		OnReconnect:    recon,
		OnDisconnect:   discon,
		OnPause:        pause,
		OnUnpause:      unpause,
		OnCommonRecv:   commonRecv,
		OnCommonSend:   commonSend,
		OnStateChanged: stateChanged,
		OnGetStateData: getStateData,
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

func pause() {
	log.Println("pause...")

}

func unpause() {
	log.Println("...unpause!")
}

func commonSend() []byte {
	return nil
}

func commonRecv(buf []byte) {
}