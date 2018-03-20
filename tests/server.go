package main

import (
	"fmt"
	"github.com/Shnifer/magellan/network"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

type roomDummy struct{}

var neededRoles = []string{"Pilot", "Navigator"}

type RoomCommonState map[string]string

var RoomMu sync.RWMutex
var DummyRoomState RoomCommonState

func (rd *roomDummy) GetRoomCommon(room string) ([]byte, error) {
	RoomMu.RLock()
	defer RoomMu.RUnlock()
	str := fmt.Sprintln(DummyRoomState)
	log.Println("GetRoomCommon(", room, ") ", str)
	return []byte(str), nil
}

func (rd *roomDummy) SetRoomCommon(room string, r io.Reader) error {
	RoomMu.Lock()
	defer RoomMu.Unlock()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("ERROR! roomDummy.SetRoomCommon cant read io.Reader")
	}
	str := string(b)
	log.Println("SetRoomCommon", room, " ", str)
	parts := strings.Split(str, " ")
	if len(parts) > 1 {
		DummyRoomState[parts[0]] = parts[1]
	}
	return nil
}

func (rd *roomDummy) CheckRoomFull(members network.RoomMembers) bool {
	for _, role := range neededRoles {
		if !members[role] {
			return false
		}
	}
	return true
}

var DummyStateData string

func (rd *roomDummy) RdyStateData(room string, state string) {

	//random delay for io operation with side projects like engineer DB
	n := 1 + rand.Intn(3)
	for i := 0; i < n; i++ {
		time.Sleep(time.Second)
	}
	DummyStateData = "DummyStateData for state=" + state
}

func (rd *roomDummy) GetStateData(room string) []byte {
	return []byte(DummyStateData)
}

func main() {
	rooms := &roomDummy{}
	DummyRoomState = make(map[string]string)

	opts := network.ServerOpts{
		Addr:     ":8000",
		RoomServ: rooms,
	}

	server, err := network.NewServer(opts)
	_ = server
	if err != nil {
		panic(err)
	}

	//waiting for enter to stop server
	str := ""
	fmt.Scanln(&str)
}
