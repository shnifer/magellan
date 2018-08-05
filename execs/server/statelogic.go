package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/static"
	"github.com/Shnifer/magellan/v2"
	"time"
)

func generateCommonData(common CommonData, stateData StateData, newState, prevState State) CommonData {
	defer LogFunc("generateCommonData " + newState.StateID + " " + newState.GalaxyID + " " + newState.ShipID)()

	//DROP inter clients params
	common.PilotData.HeatProduction = 0
	common.PilotData.ThrustVector = v2.ZV
	common.NaviData.SonarDir = common.PilotData.Ship.Ang
	common.NaviData.IsOrbiting = false
	common.NaviData.IsScanning = false
	common.ServerData.OtherShips = nil

	sessionTime := time.Now().Sub(StartDateTime).Seconds()
	common.PilotData.SessionTime = sessionTime
	stateData.Galaxy.Update(sessionTime)

	switch newState.StateID {
	case STATE_cosmo:
		//start of fly
		if prevState.StateID == STATE_login {
			prepareStartCommon(&common, stateData)
		}

		//from warp to cosmo
		if prevState.StateID == STATE_warp {
			common.PilotData.Ship.Pos =
				v2.InDir(180 + common.PilotData.Ship.Ang).Mul(stateData.Galaxy.SpawnDistance)
			common.PilotData.Ship.Vel = v2.InDir(common.PilotData.Ship.Ang)
		}
	case STATE_warp:
		//from cosmo to warp
		common = toWarpCommonData(common, stateData, newState, prevState)
	}

	return common
}

func toWarpCommonData(common CommonData, stateData StateData, newState, prevState State) CommonData {
	fromSystem := prevState.GalaxyID
	var pos v2.V2
	var spawnRange float64
	if fromSystem == ZERO_Galaxy_ID {
		pos = common.PilotData.WarpPos
		spawnRange = 0
	} else {
		var systemPoint *GalaxyPoint
		var found bool
		for _, v := range stateData.Galaxy.Points {
			if v.ID == fromSystem {
				systemPoint = v
				found = true
				break
			}
		}
		if !found {
			Log(LVL_ERROR, "toWarpCommonData: can't find system", fromSystem, "on warp map!")
			return common
		}

		pos = systemPoint.Pos
		spawnRange = systemPoint.WarpSpawnDistance
	}
	ang := common.PilotData.Ship.Ang

	ship := RBData{
		Pos:    pos.AddMul(v2.InDir(ang), spawnRange),
		Vel:    v2.InDir(ang).Mul(DEFVAL.StartWarpSpeed),
		AngVel: 0,
		Ang:    ang,
	}
	common.PilotData.Ship = ship
	common.PilotData.Distortion = DEFVAL.MinDistortion
	common.PilotData.Dir = ang
	common.NaviData.SonarDir = ship.Ang

	return common
}

func (rd *roomServer) IsValidState(roomName string, stateStr string) bool {
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
		server.AddCommand(roomName, CMD_STATECHANGEFAIL)
	}

	return res
}

func prepareStartCommon(common *CommonData, stateData StateData) {
	common.PilotData.Ship.Pos =
		stateData.Galaxy.Points[DEFVAL.SolarStartLocationName].Pos
	common.NaviData.Mines = make([]string, len(stateData.BSP.Mines))
	for i, v := range stateData.BSP.Mines {
		common.NaviData.Mines[i] = v.Owner
	}
	common.NaviData.Landing = make([]string, len(stateData.BSP.Modules))
	for i, v := range stateData.BSP.Modules {
		common.NaviData.Landing[i] = v.Owner
	}
	common.NaviData.BeaconCount = stateData.BSP.Beacons.Count
	common.EngiData.InV = [8]uint16{}
	common.EngiData.AZ = getStartAZ(stateData)
	common.EngiData.Counters.Fuel = stateData.BSP.Fuel_tank.Fuel_volume
	common.EngiData.Counters.Air = stateData.BSP.Lss.Air_volume
}

//run internal mutex call
func (rd *roomServer) isValidFlyShip(roomName string, shipID string) bool {
	if roomName == "" || shipID == "" {
		return false
	}

	rd.stateMu.RLock()
	defer rd.stateMu.RUnlock()

	for room, state := range rd.curState {
		if room != roomName && state.ShipID == shipID {
			return false
		}
	}

	return static.Exist("DB", "bsp_"+shipID+".json")
}

//run internal mutex call
func (rd *roomServer) isValidFlyGalaxy(galaxyID string) bool {
	if galaxyID == ZERO_Galaxy_ID {
		return true
	}
	return static.Exist("DB", "galaxy_"+galaxyID+".json")
}

func daemonUpdateOtherShips(rs *roomServer, updatePeriodMs int) {
	for {
		doUpdateOtherShips(rs)
		time.Sleep(time.Duration(updatePeriodMs) * time.Millisecond)
	}
}

func doUpdateOtherShips(rs *roomServer) {
	rs.stateMu.RLock()
	defer rs.stateMu.RUnlock()

	//map[galaxyName][]rooms
	m := make(map[string][]string, len(rs.curState))
	//map[room]galaxyName
	r := make(map[string]string, len(rs.curState))
	for room, state := range rs.curState {
		if (state.StateID != STATE_cosmo && state.StateID != STATE_warp) || state.GalaxyID == "" {
			continue
		}
		galaxy := state.GalaxyID
		m[galaxy] = append(m[state.GalaxyID], room)
		r[room] = galaxy
	}

	for galaxy, rooms := range m {
		if len(rooms) < 2 {
			delete(m, galaxy)
		}
	}

	rs.stateDataMu.RLock()
	defer rs.stateDataMu.RUnlock()

	rs.commonMu.Lock()
	defer rs.commonMu.Unlock()

	var otherShip OtherShipData
	for room, galaxy := range r {
		CD := rs.commonData[room]
		CD.ServerData.MsgID++
		CD.ServerData.OtherShips = CD.ServerData.OtherShips[:0]
		otherRooms, ok := m[galaxy]
		if !ok {
			rs.commonData[room] = CD
			continue
		}
		for _, otherRoom := range otherRooms {
			if otherRoom == room {
				continue
			}
			sd, ok := rs.stateData[otherRoom]
			if !ok {
				continue
			}
			if sd.BSP == nil {
				continue
			}
			ocd, ok := rs.commonData[otherRoom]
			if !ok {
				continue
			}
			if ocd.PilotData == nil {
				continue
			}
			otherShip = OtherShipData{
				Id:   rs.curState[otherRoom].ShipID,
				Name: sd.BSP.Ship.Name,
				Ship: ocd.PilotData.Ship,
			}
			CD.ServerData.OtherShips = append(CD.ServerData.OtherShips, otherShip)

		}
		rs.commonData[room] = CD
	}
}

func getStartAZ(sd StateData) [8]float64 {
	return [8]float64{
		sd.BSP.March_engine.AZ,
		sd.BSP.Shunter.AZ,
		sd.BSP.Warp_engine.AZ,
		sd.BSP.Shields.AZ,
		sd.BSP.Radar.AZ,
		sd.BSP.Scanner.AZ,
		sd.BSP.Fuel_tank.AZ,
		sd.BSP.Lss.AZ,
	}
}
