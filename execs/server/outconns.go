package main

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/static"
	"github.com/Shnifer/magellan/storage"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"strconv"
	"log"
)

func (rd *roomServer) loadStateData(state State) (sd StateData, subscribe chan storage.Event) {
	defer LogFunc("loadStateData")()

	if state.ShipID != "" {
		sd.BSP = loadShipState(state.ShipID)
	}

	if state.GalaxyID != "" {
		if state.GalaxyID == ZERO_Galaxy_ID {
			sd.Galaxy = zeroGalaxy()
		} else {
			sd.Galaxy = loadGalaxyState(state.GalaxyID)
		}
		sd.Buildings, subscribe = loadBuildingsAndSubscribe(rd.storage, state.GalaxyID)
	}

	return sd, subscribe
}

func loadShipState(shipID string) *BSP {
	var res BSP

	buf, err := static.Load("DB", "bsp_"+shipID+".json")

	if err != nil {
		Log(LVL_ERROR, "Can't open file for ShipID ", shipID)
		return nil
	}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal file for ship", shipID, err)
		return nil
	}
	return &res
}

func loadGalaxyState(GalaxyID string) *Galaxy {
	var res Galaxy
	buf, err := static.Load("DB", "galaxy_"+GalaxyID+".json")
	if err != nil {
		Log(LVL_ERROR, "Can't open file for galaxyID ", GalaxyID)
		return nil
	}

	err = json.Unmarshal(buf, &res)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal file for galaxy", GalaxyID)
		return nil
	}

	//First restore ID's
	for id, v := range res.Points {
		if v.ID == "" {
			v.ID = id
			res.Points[id] = v
		}
	}
	//Second - recalc lvls!
	res.RecalcLvls()

	return &res
}

func zeroGalaxy() *Galaxy {
	res := Galaxy{
		Points: make(map[string]*GalaxyPoint),
	}
	res.RecalcLvls()

	return &res
}

func loadBuildingsAndSubscribe(storage *storage.Storage, GalaxyID string) (builds map[string]Building, subscribe chan storage.Event) {
	diskData, subscribe := storage.SubscribeAndData(GalaxyID)
	builds = make(map[string]Building, len(diskData))
	for objectKey, data := range diskData {
		b, err := Building{}.Decode([]byte(data))
		if err != nil {
			Log(LVL_ERROR, "Wrong diskData", string(data))
			continue
		}
		b.FullKey = objectKey.FullKey()
		builds[b.FullKey] = b
	}
	return builds, subscribe
}

//save examples of DB data
func init() {
	saveDataExamples("")
}

func saveDataExamples(path string) {
	bsp, _ := json.Marshal(BSP{})

	bufBsp := bytes.Buffer{}
	json.Indent(&bufBsp, bsp, "", "    ")
	ioutil.WriteFile(path+"example_bsp.json", bufBsp.Bytes(), 0)

	galaxy := Galaxy{Points: make(map[string]*GalaxyPoint)}
	galaxy.Points["samplePoint"] = &GalaxyPoint{
		ParentID:          "parentID",
		Orbit:             100,
		Period:            80,
		Mass:              10,
		WarpRedOutDist:    10,
		WarpGreenInDist:   20,
		WarpGreenOutDist:  30,
		WarpYellowOutDist: 40,
		WarpSpawnDistance: 80,
		ScanData:          "Eurika!",
		Emissions:         []Emission{{Type: "Heat", MainRange: 100, FarValue: 200, MainValue: 100, FarRange: 200}},
	}
	galaxyStr, _ := json.Marshal(galaxy)
	//bufGalaxy := bytes.Buffer{}
	//json.Indent(&bufGalaxy, galaxyStr, "", "    ")
	ioutil.WriteFile(path+"example_galaxy.json", galaxyStr, 0)
}

func RequestHyShip(shipID string) (data []byte, exist bool) {
	shipN, err := strconv.Atoi(shipID)
	if err != nil {
		return nil, false
	}

	client := &http.Client{
		Timeout: time.Second,
	}
	body := fmt.Sprintf("{\"flight_id\":%v}", shipN)
	bodyBuf := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(http.MethodPost, DEFVAL.ShipsRequestHyServerAddr, bodyBuf)
	if err != nil {
		Log(LVL_ERROR, "can't request Hy flight data with request ", body, ":", err)
		return nil, false
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		Log(LVL_ERROR, "can't request Hy flight data with request ", body, ":", err)
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode!=200{
		return nil, false
	}

	data,err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(LVL_ERROR, "can't read responce with request ", body, ":", err)
		return nil, false
	}
	log.Println(string(data))
	return data, true
}
