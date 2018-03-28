package main

import (
	"io"
	"io/ioutil"
	"github.com/Shnifer/magellan/network"
	"sync"
	. "github.com/Shnifer/magellan/commons"
	"errors"
)

type roomServer struct{
	mu          sync.RWMutex
	stateData   map[string]MapData

	//map[roomName]commonDataMap
	commonData  map[string]MapData
	neededRoles []string
}

func newRoomServer() *roomServer{
	stateData := make(map[string]MapData)
	commonData := make(map[string]MapData)

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
		commonData:=make(MapData)
		rd.commonData[room] = commonData
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
	cd,err:=MapData{}.Decode(b)
	if err!=nil{
		err:=errors.New("SetRoomCommon: Can't decode")
		Log(LVL_ERROR, err)
		return err
	}

	if _,ok:=rd.commonData[room];!ok {
		rd.commonData[room] = make(MapData)
	}

	for key,val:=range cd {
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

func (rd *roomServer) RdyStateData(room string, state string) {
	rd.mu.Lock()
	defer rd.mu.Unlock()
	rd.stateData[room] = loadStateData(state)
}

func loadStateData(str string) MapData{
	md:=make(MapData)
	state:=State{}.Decode(str)
	if state.ShipID!=""{
		md[PART_BSP] = loadShipState(state.ShipID)
	}
	if state.GalaxyID!=""{
		md[PART_Galaxy] = loadGalaxyState(state.GalaxyID)
	}
	return md
}

const DBPath = "res/server/DB/"

//TODO: look in DB
func loadShipState(shipID string) string{
	buf, err:= ioutil.ReadFile(DBPath+"BSP_"+shipID+".json")
	if err!=nil{
		Log(LVL_ERROR, "Can't open file for ShipID ", shipID)
		return ""
	}
	return string(buf)
}

//TODO: look in DB
func loadGalaxyState(GalaxyID string) string{
	buf, err:= ioutil.ReadFile(DBPath+"Galaxy_"+GalaxyID+".json")
	if err!=nil{
		Log(LVL_ERROR, "Can't open file for galaxyID ", GalaxyID)
		return ""
	}
	return string(buf)
}


func (rd *roomServer) GetStateData(room string) []byte {
	rd.mu.RLock()
	defer rd.mu.RUnlock()
	commonData,ok := rd.stateData[room]
	if !ok{
		err:=errors.New("GetStateData: Room "+room+" not found")
		Log(LVL_ERROR, err)
		return nil
	}
	msg,err:= commonData.Encode()
	if err!=nil{
		Log(LVL_ERROR, err)
		return nil
	}
	return msg
}

//save examples of DB data
func init(){
	SaveDataExamples(DBPath)
}