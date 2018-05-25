package main

import (
	"errors"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
	"github.com/Shnifer/magellan/network/storage"
	"sync"
	"time"
)

type roomServer struct {
	storage *storage.Storage

	stateMu  sync.RWMutex
	curState map[string]State

	stateDataMu sync.RWMutex
	stateData   map[string]StateData

	subsMu     sync.RWMutex
	subscribes map[string]chan storage.Event

	//map[roomName]commonData
	commonMu   sync.RWMutex
	commonData map[string]CommonData
}

func newRoomServer(disk *storage.Storage) *roomServer {
	stateData := make(map[string]StateData)
	commonData := make(map[string]CommonData)
	curState := make(map[string]State)
	subscribes := make(map[string]chan storage.Event)

	roomServ := &roomServer{
		stateData:  stateData,
		curState:   curState,
		commonData: commonData,
		subscribes: subscribes,
		storage:    disk,
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
	rs.subsMu.RLock()
	defer rs.subsMu.RUnlock()

	for roomName, subscribe := range rs.subscribes {
		if subscribe == nil {
			continue
		}
		for event := range subscribe {
			server.AddCommand(roomName, EventToCommand(event))
		}
	}
}

func (rd *roomServer) GetRoomCommon(room string) ([]byte, error) {
	defer LogFunc("GetRoomCommon")()

	rd.commonMu.RLock()
	commonData, ok := rd.commonData[room]
	rd.commonMu.RUnlock()

	if !ok {
		commonData = CommonData{}.Empty()
		rd.commonData[room] = commonData
	}

	msg := commonData.Encode()
	return msg, nil
}

func (rd *roomServer) SetRoomCommon(room string, data []byte) error {
	defer LogFunc("SetRoomCommon")()

	cd, err := CommonData{}.Decode(data)
	if err != nil {
		err := errors.New("SetRoomCommon: Can't decode AS CommonData")
		Log(LVL_ERROR, err)
		return err
	}

	rd.commonMu.Lock()
	defer rd.commonMu.Unlock()

	dst, ok := rd.commonData[room]
	if !ok {
		dst = CommonData{}
	}
	cd.FillNotNil(&dst)
	rd.commonData[room] = dst

	return nil
}

func (rd *roomServer) RdyStateData(room string, stateStr string) {
	defer LogFunc("RdyStateData")()

	rd.stateMu.Lock()
	prevState := rd.curState[room]
	state := State{}.Decode(stateStr)
	rd.curState[room] = state
	rd.stateMu.Unlock()

	stateData, subscribe := rd.loadStateData(state)

	rd.stateDataMu.Lock()
	rd.stateData[room] = stateData
	rd.stateDataMu.Unlock()

	rd.subsMu.Lock()
	if rd.subscribes[room] != nil {
		rd.storage.Unsubscribe(rd.subscribes[room])
	}
	rd.subscribes[room] = subscribe
	rd.subsMu.Unlock()

	rd.commonMu.RLock()
	commonData, ok := rd.commonData[room]
	rd.commonMu.RUnlock()

	if !ok {
		commonData = CommonData{}.Empty()

		rd.commonMu.Lock()
		rd.commonData[room] = commonData
		rd.commonMu.Unlock()
	}
	genData := generateCommonData(commonData, stateData, state, prevState)

	rd.commonMu.Lock()
	rd.commonData[room] = genData
	rd.commonMu.Unlock()
}

func (rd *roomServer) GetStateData(room string) []byte {
	defer LogFunc("GetStateData")()

	rd.stateDataMu.RLock()
	stateData, ok := rd.stateData[room]
	rd.stateDataMu.RUnlock()

	if !ok {
		err := errors.New("GetStateData: Room " + room + " not found")
		Log(LVL_ERROR, err)
		return nil
	}

	msg := stateData.Encode()

	return msg
}

func (rd *roomServer) OnCommand(room, role, command string) {
	Log(LVL_DEBUG, "room", room, "role", role, "command", command)
}
