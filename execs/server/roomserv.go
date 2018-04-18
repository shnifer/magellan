package main

import (
	"errors"
	. "github.com/Shnifer/magellan/commons"
	"sync"
)

type roomServer struct {
	mu        sync.RWMutex
	stateData map[string]StateData
	curState  map[string]State

	//map[roomName]commonDataMap
	commonData map[string]CommonData
}

func newRoomServer() *roomServer {
	stateData := make(map[string]StateData)
	commonData := make(map[string]CommonData)
	curState := make(map[string]State)

	return &roomServer{
		stateData:  stateData,
		curState:   curState,
		commonData: commonData,
	}
}

func (rd *roomServer) GetRoomCommon(room string) ([]byte, error) {
	defer LogFunc("GetRoomCommon")()

	rd.mu.RLock()
	defer rd.mu.RUnlock()

	commonData, ok := rd.commonData[room]
	if !ok {
		commonData = CommonData{}.Empty()
		rd.commonData[room] = commonData
	}

	msg := commonData.Encode()
	return msg, nil
}

func (rd *roomServer) SetRoomCommon(room string, data []byte) error {
	defer LogFunc("SetRoomCommon")()

	rd.mu.Lock()
	defer rd.mu.Unlock()

	cd, err := CommonData{}.Decode(data)
	if err != nil {
		err := errors.New("SetRoomCommon: Can't decode AS CommonData")
		Log(LVL_ERROR, err)
		return err
	}

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

	rd.mu.Lock()
	defer rd.mu.Unlock()

	prevState := rd.curState[room]
	state := State{}.Decode(stateStr)
	rd.curState[room] = state
	stateData := loadStateData(state)
	rd.stateData[room] = stateData

	commonData, ok := rd.commonData[room]
	if !ok {
		commonData = CommonData{}.Empty()
		rd.commonData[room] = commonData
	}

	rd.commonData[room] = generateCommonData(commonData, stateData, state, prevState)
}

func (rd *roomServer) GetStateData(room string) []byte {
	defer LogFunc("GetStateData")()

	rd.mu.RLock()
	defer rd.mu.RUnlock()

	stateData, ok := rd.stateData[room]
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
