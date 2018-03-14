package main

import (
	"github.com/Shnifer/magellan/network"
	"net"
	"log"
	"time"
	"runtime"
)

func main(){
	for {
		listener, err := net.Listen("tcp", "localhost:6666")
		if err != nil {
			log.Println(err)
		}

		Server := network.NewServer(listener)
		log.Println("server created, num or routines: ", runtime.NumGoroutine())
		time.Sleep(time.Second/4)
		Server.Close()
		log.Println("server closed, num or routines: ", runtime.NumGoroutine())
		}
}