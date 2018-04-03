package main

import (
	"errors"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/network"
	"io"
	"io/ioutil"
	"sync"
	"os"
	"log"
)

type roomServer struct {
	mu        sync.RWMutex
	stateData map[string]CMapData
	curState map[string]State

	//map[roomName]commonDataMap
	commonData  map[string]CMapData
	neededRoles []string
}

func newRoomServer() *roomServer {
	stateData := make(map[string]CMapData)
	commonData := make(map[string]CMapData)
	curState := make(map[string]State)

	return &roomServer{
		stateData:   stateData,
		curState:    curState,
		commonData:  commonData,
		neededRoles: DEFVAL.NeededRoles,
	}
}

func (rd *roomServer) GetRoomCommon(room string) ([]byte, error) {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	commonData, ok := rd.commonData[room]
	if !ok {
		commonData := make(CMapData)
		rd.commonData[room] = commonData
	}

	msg, err:= commonData.Encode()
	log.Println(string(msg))
	return msg, err
}

func (rd *roomServer) SetRoomCommon(room string, r io.Reader) error {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		Log(LVL_ERROR, "SetRoomCommon cant read io.Reader")
	}
	cd, err := CMapData{}.Decode(b)
	if err != nil {
		err := errors.New("SetRoomCommon: Can't decode AS cMapData")
		Log(LVL_ERROR, err)
		return err
	}

	if _, ok := rd.commonData[room]; !ok {
		rd.commonData[room] = make(CMapData)
	}

	for key, val := range cd {
		rd.commonData[room][key] = val
	}
	return nil
}

func (rd *roomServer) CheckRoomFull(members network.RoomMembers) bool {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	for _, role := range rd.neededRoles {
		if !members[role] {
			return false
		}
	}
	return true
}

func (rd *roomServer) RdyStateData(room string, stateStr string) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	state:=State{}.Decode(stateStr)
	rd.curState[room] =  state
	rd.stateData[room] = loadStateData(state)
}

func loadStateData(state State) CMapData {
	md := make(CMapData)

	if state.ShipID != "" {
		md[PARTSTATE_BSP] = loadShipState(state.ShipID)
	}

	if state.GalaxyID != "" {
		md[PARTSTATE_Galaxy] = loadGalaxyState(state.GalaxyID)
	}

	return md
}

const DBPath = "res/server/DB/"

//TODO: look in DB
func loadShipState(shipID string) string {
	buf, err := ioutil.ReadFile(DBPath + "BSP_" + shipID + ".json")
	if err != nil {
		Log(LVL_ERROR, "Can't open file for ShipID ", shipID)
		return ""
	}
	return string(buf)
}

//TODO: look in DB
func loadGalaxyState(GalaxyID string) string {
	buf, err := ioutil.ReadFile(DBPath + "Galaxy_" + GalaxyID + ".json")
	if err != nil {
		Log(LVL_ERROR, "Can't open file for galaxyID ", GalaxyID)
		return ""
	}
	return string(buf)
}

func (rd *roomServer) GetStateData(room string) []byte {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	commonData, ok := rd.stateData[room]
	if !ok {
		err := errors.New("GetStateData: Room " + room + " not found")
		Log(LVL_ERROR, err)
		return nil
	}
	msg, err := commonData.Encode()
	if err != nil {
		Log(LVL_ERROR, err)
		return nil
	}

	return msg
}

func (rd *roomServer) IsValidState(roomName string, stateStr string) bool {
	rd.mu.RLock()
	defer rd.mu.RUnlock()

	state := State{}.Decode(stateStr)
	switch state.Special {
	case STATE_login:
		return state.GalaxyID == "" && state.ShipID == ""
	case STATE_cosmo:
		return rd.isValidFlyShip(roomName, state.ShipID) && rd.isValidFlyGalaxy(state.GalaxyID)
	case STATE_warp:
		return rd.isValidFlyShip(roomName, state.ShipID) && rd.isValidFlyGalaxy(state.GalaxyID)
	}
	return false
}

//run internal mutex call
func (rd *roomServer) isValidFlyShip(roomName string, shipID string) bool{
	if roomName == "" || shipID =="" {
		return false
	}

	for room, state:= range rd.curState {
		if room!=roomName && state.ShipID == shipID {
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


//save examples of DB data
func init() {
	SaveDataExamples(DBPath)
}
