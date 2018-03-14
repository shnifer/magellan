package main

import (
	"github.com/Shnifer/magellan/network"
	"net"
	"log"
	"time"
	"runtime"
)

func main(){
	listener, err := net.Listen("tcp", "localhost:6666")
	if err != nil {
		log.Println(err)
	}
	Server := network.NewServer(listener)
	defer Server.Close()
	log.Println("server created, num or routines: ", runtime.NumGoroutine())

	for {
		time.Sleep(time.Second)
		log.Println("num or routines: ", runtime.NumGoroutine())
		}
}