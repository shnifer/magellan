package main

import (
	"bytes"
	"encoding/json"
	"errors"
	. "github.com/Shnifer/magellan/commons"
	"io/ioutil"
	"os"
	"sync"
	"time"
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

	state := State{}.Decode(stateStr)
	rd.curState[room] = state
	rd.stateData[room] = loadStateData(state)

	commonData, ok := rd.commonData[room]
	if !ok {
		commonData = CommonData{}.Empty()
		rd.commonData[room] = commonData
	}

	rd.commonData[room] = generateCommonData(commonData, state)
}

func generateCommonData(prevCommon CommonData, newState State) CommonData {
	defer LogFunc("generateCommonData " + newState.StateID + " " + newState.GalaxyID + " " + newState.ShipID)()
	prevCommon.PilotData.SessionTime = time.Now().Sub(StartDateTime).Seconds()
	return prevCommon
}

func loadStateData(state State) StateData {
	defer LogFunc("loadStateData")()

	var sd StateData

	//sd.ServerTime = time.Now()

	if state.ShipID != "" {
		sd.BSP = loadShipState(state.ShipID)
	}

	if state.GalaxyID != "" {
		sd.Galaxy = loadGalaxyState(state.GalaxyID)
	}

	return sd
}

const DBPath = "res/server/DB/"

//TODO: look in DB
func loadShipState(shipID string) *BSP {
	var res BSP
	buf, err := ioutil.ReadFile(DBPath + "BSP_" + shipID + ".json")
	if err != nil {
		Log(LVL_ERROR, "Can't open file for ShipID ", shipID)
		return nil
	}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal file for ship", shipID)
		return nil
	}
	return &res
}

//TODO: look in DB
func loadGalaxyState(GalaxyID string) *Galaxy {
	var res Galaxy
	buf, err := ioutil.ReadFile(DBPath + "Galaxy_" + GalaxyID + ".json")
	if err != nil {
		Log(LVL_ERROR, "Can't open file for galaxyID ", GalaxyID)
		return nil
	}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal file for galaxy", GalaxyID)
		return nil
	}
	return &res
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

func (rd *roomServer) IsValidState(roomName string, stateStr string) bool {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	var res bool
	state := State{}.Decode(stateStr)
	switch state.StateID {
	case STATE_login:
		res = state.GalaxyID == "" && state.ShipID == ""
	case STATE_cosmo:
		res = rd.isValidFlyShip(roomName, state.ShipID) && rd.isValidFlyGalaxy(state.GalaxyID)
	case STATE_warp:
		res = rd.isValidFlyShip(roomName, state.ShipID) && rd.isValidFlyGalaxy(state.GalaxyID)
	}

	if !res {
		server.AddCommand(roomName, CMD_STATECHANGRFAIL)
	}

	return res
}

//run internal mutex call
func (rd *roomServer) isValidFlyShip(roomName string, shipID string) bool {
	if roomName == "" || shipID == "" {
		return false
	}

	for room, state := range rd.curState {
		if room != roomName && state.ShipID == shipID {
			return false
		}
	}

	if _, err := os.Stat(DBPath + "BSP_" + shipID + ".json"); os.IsNotExist(err) {
		return false
	}

	return true
}

//run internal mutex call
func (rd *roomServer) isValidFlyGalaxy(galaxyID string) bool {
	if _, err := os.Stat(DBPath + "Galaxy_" + galaxyID + ".json"); os.IsNotExist(err) {
		return false
	}

	return true
}

func (rd *roomServer) OnCommand(room, role, command string) {
	Log(LVL_DEBUG, "room", room, "role", role, "command", command)
}

//save examples of DB data
func init() {
	saveDataExamples(DBPath)
}

func saveDataExamples(path string) {
	bsp, _ := json.Marshal(BSP{})

	bufBsp := bytes.Buffer{}
	json.Indent(&bufBsp, bsp, "", "    ")
	ioutil.WriteFile(path+"example_bsp.json", bufBsp.Bytes(), 0)

	galaxy := Galaxy{}
	galaxy.Points = append(galaxy.Points, GalaxyPoint{})
	galaxyStr, _ := json.Marshal(galaxy)
	bufGalaxy := bytes.Buffer{}
	json.Indent(&bufGalaxy, galaxyStr, "", "    ")
	ioutil.WriteFile(path+"example_galaxy.json", bufGalaxy.Bytes(), 0)
}
