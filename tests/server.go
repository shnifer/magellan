package main

import (
	"io/ioutil"
	"log"
	"github.com/Shnifer/magellan/network"
	"io"
	"fmt"
)

type roomDummy struct{}
var neededRoles  = []string{"Pilot","Navigator"}

func (rd *roomDummy) GetRoomCommon(room string) ([]byte, error) {
	log.Println("GetRoomCommon("+room+") ")
	return []byte("dummy common state"),nil
}

func (rd *roomDummy) SetRoomCommon(room string, r io.Reader) error {
	str,err:=ioutil.ReadAll(r)
	if err!=nil{
		log.Println("ERROR! roomDummy.SetRoomCommon cant read io.Reader")
	}
	log.Println("SetRoomCommon",room," ",str)
	return nil
}

func (rd *roomDummy) CheckRoomFull(members network.RoomMembers) bool{
	for _,role := range neededRoles{
		if !members[role] {
			return false
		}
	}
	return true
}

func main() {
	rooms:=&roomDummy{}

	opts:=network.ServerOpts{
		Addr:":8000",
		RoomServ: rooms,
	}
	server,err:=network.NewServer(opts)
	_ = server
	if err!=nil{
		panic(err)
	}

	//waiting for enter to stop server
	str:=""
	fmt.Scanln(&str)
}
