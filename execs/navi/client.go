package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
)

var Client *network.Client
var NetData TData

func initClient() {
	opts := network.ClientOpts{
		Addr:           DEFVAL.Port,
		Room:           DEFVAL.Room,
		Role:           DEFVAL.Role,
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
	NetData.SetStateData(sd)

	//first load data in scene, and only than count as done - so Client reports new state ready
	initSceneState()
}

func commonSend() []byte {
	defer LogFunc("commonSend")()
	return NetData.CommonPartEncoded(DEFVAL.Role)
}

func commonRecv(buf []byte, readOwnPart bool) {
	defer LogFunc("commonRecv")()
	cd, err := CommonData{}.Decode(buf)
	if err != nil {
		panic("commonRecv: Can't decode CommonData " + err.Error())
	}
	NetData.LoadCommonData(cd, DEFVAL.Role, readOwnPart)
}
