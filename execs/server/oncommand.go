package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/storage"
	"strconv"
	"strings"
)

func (rs *roomServer) OnCommand(room, role, command string) {
	Log(LVL_DEBUG, "room", room, "role", role, "command", command)

	switch {
	case strings.HasPrefix(command, CMD_ADDBUILDREQ):
		go rs.AddBuildCommand(command, room, role)
	case strings.HasPrefix(command, CMD_DELBUILDREQ):
		go rs.DelBuildCommand(command)
	case strings.HasPrefix(command, CMD_LOGGAMEEVENT):
		buf := []byte(command[len(CMD_LOGGAMEEVENT):])
		lge, err := LogGameEvent{}.Decode(buf)
		if err != nil {
			Log(LVL_ERROR, err)
			return
		}
		go SaveToStorage(lge.Key, lge.Args, lge.StateFields)
	case command == CMD_GRACEENDDIE || command == CMD_GRACEENDRETURN:
		alive := command == CMD_GRACEENDRETURN
		go rs.GraceDie(room, alive)
	}
}

func (rs *roomServer) AddBuildCommand(command, room, role string) {
	buildStr := strings.TrimPrefix(command, CMD_ADDBUILDREQ)
	b, err := Building{}.Decode([]byte(buildStr))
	if err != nil {
		Log(LVL_ERROR, "Command CMD_ADDBUILDREQ sent strange building: "+buildStr)
		return
	}
	if b.GalaxyID == ZERO_Galaxy_ID {
		Log(LVL_INFO, "build in zero galaxy, idiot")
		return
	}

	key := b.Type + strconv.Itoa(rs.storage.NextID())

	err = rs.storage.Add(b.GalaxyID, key, buildStr)
	if err != nil {
		Log(LVL_ERROR, "OnCommand: room", room, "role", role, "command", command, ":", err)
		return
	}

	//duplicate in warp
	if b.Type == BUILDING_BEACON || b.Type == BUILDING_BLACKBOX {
		b.PlanetID = b.GalaxyID
		b.GalaxyID = WARP_Galaxy_ID
		b.Period = 0
		buildStr = string(b.Encode())
		err = rs.storage.Add(b.GalaxyID, key, buildStr)
		if err != nil {
			Log(LVL_ERROR, "OnCommand: room", room, "role", role, "command", command, ":", err)
			return
		}
	}
}

func (rs *roomServer) DelBuildCommand(command string) {
	fullKey := command[len(CMD_DELBUILDREQ):]
	objKey, err := storage.ReadKey(fullKey)
	if err != nil {
		Log(LVL_ERROR, err)
		return
	}
	err = rs.storage.Remove(objKey)
	if err != nil {
		Log(LVL_ERROR, err)
		return
	}

	//try to delete in warp, if never was -- okey, Keys are unique, we marked as deleted something that will never spawn
	objKey.Area = WARP_Galaxy_ID
	err = rs.storage.Remove(objKey)
	if err != nil {
		Log(LVL_ERROR, err)
		return
	}
}

func (rs *roomServer) reportGraceDieToHy(room string, alive bool) {
	rs.RLock()
	defer rs.RUnlock()

	holder := rs.getHolder(room)
	sd := holder.getStateData()
	cd := holder.getCommon()
	startAZ := getStartAZ(sd)
	flightId := sd.BSP.FlightID
	if !alive {
		reportHyDead(flightId)
	} else {
		reportHyAlive(flightId, cd.EngiData.Counters.FlightTime, startAZ, cd.EngiData.AZ)
	}
}

func (rs *roomServer) GraceDie(room string, alive bool) {
	rs.reportGraceDieToHy(room, alive)
	go server.KillRoom(room)
}
