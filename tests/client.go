package main

import (
	"github.com/Shnifer/magellan/network"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

func discon(){
	log.Println("lost connect")
}

func recon(){
	log.Println("recon!")
}

func pause(){
	log.Println("pause...")
}

func unpause(){
	log.Println("...unpause!")
}

type TConf struct {
	Room, Role string
}

func main() {
	conf := TConf{}
	buf, err:= ioutil.ReadFile("conf.txt")
	if err!=nil{
		conf.Room = "roomName"
		conf.Role = "roleName"
		b,err:=json.Marshal(conf)
		if err!=nil{
			panic(err)
		}
		var ib bytes.Buffer
		json.Indent(&ib, b, "","\n    ")
		ioutil.WriteFile("conf.txt",ib.Bytes(),0)
		log.Println("no conf file found, created new. restart")
		return
	}
	err=json.Unmarshal(buf,&conf)
	if err!=nil{
		log.Println("ERROR unmarshal conf")
		panic(err)
	}

	Opts:=network.ClientOpts{
		Addr: "http://localhost:8000",
		Room: conf.Room,
		Role: conf.Role,
		OnReconnect:  recon,
		OnDisconnect: discon,
		OnPause: pause,
		OnUnpause: unpause,
	}

	client,err:=network.NewClient(Opts)
	if err!=nil{
		panic(err)
	}

	_ = client

	//waiting for enter to stop client
	str:=""
	fmt.Scanln(&str)
}
