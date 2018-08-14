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

type loadPlan struct {
	timeout  time.Time
	state    State
	shipId   string
	restoreN int
}

type roomServer struct {
	storage *storage.Storage
	restore *diskv.Diskv

	//write lock only on add holder
	sync.RWMutex
	holders map[string]*roomHolder

	loadMu sync.Mutex
	//map[roomname]loadPlan
	loadPlans map[string]loadPlan
}

func newRoomServer(disk *storage.Storage, restore *diskv.Diskv) *roomServer {
	roomServ := &roomServer{
		holders:   make(map[string]*roomHolder),
		storage:   disk,
		restore:   restore,
		loadPlans: make(map[string]loadPlan),
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

//under rs.RLock
func (rs *roomServer) getHolder(roomName string) *roomHolder {
	for {
		if holder, ok := rs.holders[roomName]; ok && holder != nil {
			return holder
		}

		rs.RUnlock()
		rs.Lock()
		if h, ok := rs.holders[roomName]; ok && h != nil {
			rs.Unlock()
			rs.RLock()
			continue
		}
		holder := newRoomHolder(roomName)
		rs.holders[roomName] = holder
		rs.Unlock()
		rs.RLock()
	}
}

func (rs *roomServer) GetRoomCommon(room string) ([]byte, error) {
	defer LogFunc("GetRoomCommon")()
	rs.RLock()
	defer rs.RUnlock()

	commonData := rs.getHolder(room).getCommon()
	msg := commonData.Encode()
	return msg, nil

}

func (rs *roomServer) SetRoomCommon(room string, data []byte) error {
	defer LogFunc("SetRoomCommon")()
	rs.RLock()
	defer rs.RUnlock()

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
	rs.RLock()
	defer rs.RUnlock()

	holder := rs.getHolder(room)
	rs.loadMu.Lock()
	plan, ok := rs.loadPlans[room]
	rs.loadMu.Unlock()

	if ok && time.Now().Before(plan.timeout) && plan.state.Encode() == stateStr {
		//do restore
		rs.loadRestorePoint(room, plan.shipId, plan.restoreN)
	} else {
		//do usual load
		state := State{}.Decode(stateStr)
		if oldSub := holder.getSubscribe(); oldSub != nil {
			rs.storage.Unsubscribe(oldSub)
		}
		stateData, subscribe := rs.loadStateData(state)
		holder.rdyStateData(state, stateData, subscribe, true)
		holder.saveRestorePoint(rs.restore)
	}
}

func (rs *roomServer) GetStateData(room string) []byte {
	defer LogFunc("GetStateData")()
	rs.RLock()
	defer rs.RUnlock()

	stateData := rs.getHolder(room).getStateData()
	msg := stateData.Encode()

	return msg
}

func (rs *roomServer) OnKillRoom(roomName string) {
	rs.Lock()
	defer rs.Unlock()

	if holder, ok := rs.holders[roomName]; ok {
		if oldSub := holder.getSubscribe(); oldSub != nil {
			rs.storage.Unsubscribe(oldSub)
		}
	}
	delete(rs.holders, roomName)
}
