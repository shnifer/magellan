package main

import (
	"errors"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/network"
	"github.com/Shnifer/magellan/storage"
	"github.com/peterbourgon/diskv"
	"sync"
	"time"
)

type roomServer struct {
	storage *storage.Storage
	restore *diskv.Diskv

	//write lock only on add holder
	sync.RWMutex
	holders map[string]*roomHolder
}

func newRoomServer(disk *storage.Storage, restore *diskv.Diskv) *roomServer {
	roomServ := &roomServer{
		holders: make(map[string]*roomHolder),
		storage: disk,
		restore: restore,
	}
	return roomServ
}

func daemonUpdateSubscribes(rs *roomServer, server *network.Server, updatePeriodMs int) {
	for {
		doUpdateSubscribes(rs, server)
		time.Sleep(time.Duration(updatePeriodMs) * time.Millisecond)
	}
}

func doUpdateSubscribes(rs *roomServer, server *network.Server) {
	rs.RLock()
	defer rs.RUnlock()

	for _, holder := range rs.holders {
		holder.doUpdateSubscribes()
	}
}

//under w/out rs.RLock
func (rs *roomServer) getHolder(roomName string) *roomHolder {
	rs.RLock()
	if holder, ok := rs.holders[roomName]; ok {
		rs.RUnlock()
		return holder
	}
	rs.RUnlock()

	rs.Lock()
	defer rs.Unlock()

	if holder, ok := rs.holders[roomName]; ok {
		return holder
	}
	holder := newRoomHolder(roomName)
	rs.holders[roomName] = holder
	return holder
}

func (rs *roomServer) GetRoomCommon(room string) ([]byte, error) {
	defer LogFunc("GetRoomCommon")()
	commonData := rs.getHolder(room).getCommon()
	msg := commonData.Encode()
	return msg, nil

}

func (rs *roomServer) SetRoomCommon(room string, data []byte) error {
	defer LogFunc("SetRoomCommon")()

	cd, err := CommonData{}.Decode(data)
	if err != nil {
		err := errors.New("SetRoomCommon: Can't decode AS CommonData")
		Log(LVL_ERROR, err)
		return err
	}

	rs.getHolder(room).setCommon(cd)

	return nil
}

func (rs *roomServer) RdyStateData(room string, stateStr string) {
	defer LogFunc("RdyStateData")()

	holder := rs.getHolder(room)
	state := State{}.Decode(stateStr)
	if oldSub := holder.getSubscribe(); oldSub != nil {
		rs.storage.Unsubscribe(oldSub)
	}
	stateData, subscribe := rs.loadStateData(state)
	holder.rdyStateData(state, stateData, subscribe, true)
	holder.saveRestorePoint(rs.restore)
}

func (rs *roomServer) GetStateData(room string) []byte {
	defer LogFunc("GetStateData")()

	stateData := rs.getHolder(room).getStateData()
	msg := stateData.Encode()

	return msg
}
