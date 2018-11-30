package main

import (
	. "github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/network"
	"time"
)

var Client *network.Client

func initClient() {
	opts := network.ClientOpts{
		Addr:           DEFVAL.Port,
		Room:           DEFVAL.Room,
		Role:           DEFVAL.Role,
		Timeout:        time.Duration(DEFVAL.Timeout) * time.Millisecond,
		PingPeriod:     time.Duration(DEFVAL.PingPeriod) * time.Millisecond,
		OnReconnect:    recon,
		OnDisconnect:   discon,
		OnPause:        pause,
		OnUnpause:      unpause,
		OnCommonRecv:   commonRecv,
		OnCommonSend:   commonSend,
		OnStateChanged: stateChanged,
		OnGetStateData: getStateData,
		OnCommand:      onCommand,
	}

	var err error
	Client, err = network.NewClient(opts)
	if err != nil {
		panic(err)
	}
}

func getStateData(stateData []byte) {
	defer LogFunc("getStateData")()

	//get stateData
	sd, err := StateData{}.Decode(stateData)
	if err != nil {
		panic("Weird state stateData:" + err.Error())
	}
	Data.SetStateData(sd)
	Data.WaitDone()

	//first load data in scene, and only than count as done - so Client reports new state ready
	initSceneState()
}

func commonSend() []byte {
	defer LogFunc("commonSend")()
	return Data.MyPartToSend()
}

func commonRecv(buf []byte, readOwnPart bool) {
	defer LogFunc("commonRecv")()
	cd, err := CommonData{}.Decode(buf)
	if err != nil {
		panic("commonRecv: Can't decode CommonData " + err.Error())
	}
	if !readOwnPart {
		cd = cd.WithoutRole(DEFVAL.Role)
	}
	Data.LoadCommonData(cd)
}
