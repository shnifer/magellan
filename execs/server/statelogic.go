package main

import (
	. "github.com/Shnifer/magellan/commons"
	"os"
	"time"
)

func generateCommonData(prevCommon CommonData, newState State) CommonData {
	defer LogFunc("generateCommonData " + newState.StateID + " " + newState.GalaxyID + " " + newState.ShipID)()
	prevCommon.PilotData.SessionTime = time.Now().Sub(StartDateTime).Seconds()
	return prevCommon
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
