package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"io/ioutil"
)

const DBPath = "res/server/DB/"

func loadStateData(state State) StateData {
	defer LogFunc("loadStateData")()

	var sd StateData

	if state.ShipID != "" {
		sd.BSP = loadShipState(state.ShipID)
	}

	if state.GalaxyID != "" {
		sd.Galaxy = loadGalaxyState(state.GalaxyID)
	}

	return sd
}

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
