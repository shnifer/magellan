package main

import (
	"io"
	"io/ioutil"
	"math/rand"
	"github.com/Shnifer/magellan/network"
	"sync"
	"time"
	. "github.com/Shnifer/magellan/commons"
	"errors"
)

type roomServer struct{
	mu          sync.RWMutex
	stateData   map[string]string

	//map[roomName]commonDataMap
	commonData  map[string]RoomCommonData
	neededRoles []string
}

func newRoomServer() *roomServer{
	stateData := make(map[string]string)
	commonData := make(map[string]RoomCommonData)

	return &roomServer{
		stateData: stateData,
		commonData: commonData,
		neededRoles: DEFVAL.NeededRoles,
	}
}

func (rd *roomServer) GetRoomCommon(room string) ([]byte, error) {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	commonData,ok := rd.commonData[room]
	if !ok{
		err:=errors.New("GetRoomCommon: Room "+room+" not found")
		Log(LVL_ERROR, err)
		return nil, err
	}
	return commonData.Encode()
}

func (rd *roomServer) SetRoomCommon(room string, r io.Reader) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		Log(LVL_ERROR,"SetRoomCommon cant read io.Reader")
	}
	cd,err:=RoomCommonData{}.Decode(b)
	if err!=nil{
		err:=errors.New("SetRoomCommon: Can't decode")
		Log(LVL_ERROR, err)
		return err
	}
	for key,val:=range cd {
		rd.commonData[room][key] = val
	}
	return nil
}

func (rd *roomServer) CheckRoomFull(members network.RoomMembers) bool {

	for _, role := range rd.neededRoles {
		if !members[role] {
			return false
		}
	}
	return true
}

func (rd *roomServer) RdyStateData(room string, state string) {

	//random delay for io operation with side projects like engineer DB
	n := 1 + rand.Intn(3)
	for i := 0; i < n; i++ {
		time.Sleep(time.Second)
	}
	rd.stateData[room] = "dummy state = "+state
}

func (rd *roomServer) GetStateData(room string) []byte {
	return []byte(rd.stateData[room])
}
