package main

import (
	"github.com/Shnifer/magellan/network"
	."github.com/Shnifer/magellan/commons"
	"log"
	"time"
)

func main(){
	opts:=network.ClientOpts{
		Addr: DEFVAL.Port,
		Room:         DEFVAL.Room,
		Role:         ROLE_Pilot,
		OnReconnect:  recon,
		OnDisconnect: discon,
		OnPause:      pause,
		OnUnpause:    unpause,
		OnCommonRecv: commonRecv,
		OnCommonSend: commonSend,
		OnStateChanged: stateChanged,
		OnGetStateData: getStateData,
	}

	client,err:=network.NewClient(opts)
	if err!=nil{
		panic(err)
	}
	_ = client

	for {
		time.Sleep(time.Second)
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

func stateChanged(wanted string){
	log.Println("new state wanted", wanted)
}

func getStateData(data []byte){
	log.Println("Loaded State Data", string(data))
}
