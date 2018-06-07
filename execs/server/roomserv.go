package main

import (
	"errors"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/network"
	"github.com/Shnifer/magellan/storage"
	"strconv"
	"strings"
	"sync"
	"time"
)

type roomServer struct {
	storage *storage.Storage

	stateMu  sync.RWMutex
	curState map[string]State

	stateDataMu sync.RWMutex
	stateData   map[string]StateData

	//map[roomName]subscribe
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

func daemonUpdateOtherShips(rs *roomServer, updatePeriodMs int) {
	for {
		doUpdateOtherShips(rs)
		time.Sleep(time.Duration(updatePeriodMs) * time.Millisecond)
	}
}

func doUpdateOtherShips(rs *roomServer){
	//map[galaxyName][]rooms
	rs.stateMu.RLock()
	defer rs.stateMu.RUnlock()

	m:=make(map[string][]string,len(rs.curState))
	for room,state:=range rs.curState{
		if (state.StateID!=STATE_cosmo && state.StateID!=STATE_warp) || state.GalaxyID=="" {
			continue
		}
		m[state.GalaxyID] = append(m[state.GalaxyID], room)
	}

	rs.stateDataMu.RLock()
	defer rs.stateDataMu.RUnlock()

	rs.commonMu.Lock()
	defer rs.commonMu.Unlock()

	var otherShip OtherShip
	for _,rooms:=range m{
		if len(rooms)<2 {
			continue
		}
		for i,room:=range rooms {
			CD:=rs.commonData[room]
			CD.ServerData.MsgID++
			CD.ServerData.OtherShips = CD.ServerData.OtherShips[:0]
			for j,otherRoom := range rooms{
				if i!=j {
					otherShip = OtherShip{
						Id: rs.curState[room].ShipID,
//						Name: rs.stateData[room].BSP.ShipName, check here!
						Ship: rs.commonData[otherRoom].PilotData.Ship,
					}
					CD.ServerData.OtherShips =append(CD.ServerData.OtherShips, otherShip)
				}
			}
			rs.commonData[room] = CD
		}
	}
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
	innerloop:
		for {
			select {
			case event := <-subscribe:
				server.AddCommand(roomName, EventToCommand(event))
			default:
				break innerloop
			}
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

	Log(LVL_DEBUG, "RdyStateData: try rd.stateMu.Lock()")
	rd.stateMu.Lock()
	prevState := rd.curState[room]
	state := State{}.Decode(stateStr)
	rd.curState[room] = state
	rd.stateMu.Unlock()

	stateData, subscribe := rd.loadStateData(state)

	Log(LVL_DEBUG, "RdyStateData: try rd.stateDataMu.Lock()")
	rd.stateDataMu.Lock()
	rd.stateData[room] = stateData
	rd.stateDataMu.Unlock()

	Log(LVL_DEBUG, "RdyStateData: try rd.subsMu.Lock()")
	rd.subsMu.Lock()
	if rd.subscribes[room] != nil {
		rd.storage.Unsubscribe(rd.subscribes[room])
	}
	rd.subscribes[room] = subscribe
	rd.subsMu.Unlock()

	Log(LVL_DEBUG, "RdyStateData: try rd.commonMu.RLock()")
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

	Log(LVL_DEBUG, "RdyStateData: try rd.commonMu.Lock()")
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

	switch {
	case strings.HasPrefix(command, CMD_ADDBUILDREQ):
		buildStr := command[len(CMD_ADDBUILDREQ):]
		b, err := Building{}.Decode([]byte(buildStr))
		if err != nil {
			Log(LVL_ERROR, "Command CMD_ADDBUILDREQ sent strange building: "+buildStr)
			return
		}
		key := b.Type + strconv.Itoa(rd.storage.NextID())

		err = rd.storage.Add(b.GalaxyID, key, buildStr)
		if err != nil {
			Log(LVL_ERROR, "OnCommand: room", room, "role", role, "command", command, ":", err)
			return
		}

		//duplicate in warp
		if b.Type == BUILDING_BEACON {
			b.PlanetID = b.GalaxyID
			b.GalaxyID = WARP_Galaxy_ID
			b.Period = 0
			buildStr = string(b.Encode())
			err = rd.storage.Add(b.GalaxyID, key, buildStr)
			if err != nil {
				Log(LVL_ERROR, "OnCommand: room", room, "role", role, "command", command, ":", err)
			}
		}
	case strings.HasPrefix(command, CMD_DELBUILDREQ):
		fullKey := command[len(CMD_DELBUILDREQ):]
		objKey, err := storage.ReadKey(fullKey)
		if err != nil {
			Log(LVL_ERROR, err)
		}
		err = rd.storage.Remove(objKey)
		if err != nil {
			Log(LVL_ERROR, err)
		}

		//try to delete in warp, if never was -- okey, Keys are unique, we marked as deleted something that will never spawn
		objKey.Area = WARP_Galaxy_ID
		err = rd.storage.Remove(objKey)
		if err != nil {
			Log(LVL_ERROR, err)
		}
	}
}
