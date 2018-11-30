package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shnifer/magellan/network"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

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

var FrameN int

func commonSend() []byte {
	FrameN++
	return []byte(conf.Role + " " + strconv.Itoa(FrameN))
}

func commonRecv(buf []byte, readOwn bool) {
	//log.Println("commonRecv", string(buf))
}

type TConf struct {
	Room, Role string
}

var conf TConf

func stateChanged(wanted string) {
	log.Println("new state wanted", wanted)
}

func getStateData(data []byte) {
	log.Println("Loaded State Data", string(data))
}

func main() {
	buf, err := ioutil.ReadFile("conf.txt")
	if err != nil {
		conf.Room = "roomName"
		conf.Role = "roleName"
		b, err := json.Marshal(conf)
		if err != nil {
			panic(err)
		}
		var ib bytes.Buffer
		json.Indent(&ib, b, "", "\n    ")
		ioutil.WriteFile("conf.txt", ib.Bytes(), 0)
		log.Println("no conf file found, created new. restart")
		return
	}
	err = json.Unmarshal(buf, &conf)
	if err != nil {
		log.Println("ERROR unmarshal conf")
		panic(err)
	}

	Opts := network.ClientOpts{
		Addr:           "http://localhost:8000",
		Room:           conf.Room,
		Role:           conf.Role,
		OnReconnect:    recon,
		OnDisconnect:   discon,
		OnPause:        pause,
		OnUnpause:      unpause,
		OnCommonRecv:   commonRecv,
		OnCommonSend:   commonSend,
		OnStateChanged: stateChanged,
		OnGetStateData: getStateData,
	}

	client, err := network.NewClient(Opts)
	if err != nil {
		panic(err)
	}

	go func() {
		t := time.Tick(time.Second * 3)
		for {
			<-t
			pr := client.PauseReason()
			log.Println("pingLost", pr.PingLost, "Full", pr.IsFull, "Coherent", pr.IsCoherent,
				"Cur", pr.CurState, "Want", pr.WantState)
		}
	}()

	_ = client

	//waiting for enter to stop client
	str := ""
	for {
		fmt.Scanln(&str)
		if str == "quit" {
			return
		}
		client.RequestNewState(str, false)
	}
}
