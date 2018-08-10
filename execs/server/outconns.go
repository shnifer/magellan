package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/static"
	"github.com/Shnifer/magellan/storage"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (rs *roomServer) loadStateData(state State) (sd StateData, subscribe chan storage.Event) {
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
		sd.Buildings, subscribe = loadBuildingsAndSubscribe(rs.storage, state.GalaxyID)
	}

	return sd, subscribe
}

func loadShipState(shipID string) *BSP {
	var res BSP

	var buf []byte
	var err error
	if static.Exist("DB", "bsp_"+shipID+".json") {
		buf, err = static.Load("DB", "bsp_"+shipID+".json")

		if err != nil {
			Log(LVL_ERROR, "Can't open file for ShipID ", shipID)
			return nil
		}

		err = json.Unmarshal(buf, &res)
		if err != nil {
			Log(LVL_ERROR, "can't unmarshal data for ship", shipID, err)
			return nil
		}
	} else {
		var exist bool
		res, exist = RequestHyShip(shipID)
		if !exist || buf == nil {
			Log(LVL_ERROR, "Can't get Hy data for ShipID ", shipID)
			return nil
		}
	}
	log.Println(res)
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

func RequestHyShip(shipID string) (dat BSP, exist bool) {
	shipN, err := strconv.Atoi(shipID)
	if err != nil {
		return BSP{}, false
	}

	client := &http.Client{
		Timeout: time.Second,
	}
	body := fmt.Sprintf("{\"flight_id\":%v}", shipN)
	bodyBuf := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(http.MethodPost, DEFVAL.ShipsRequestHyServerAddr, bodyBuf)
	if err != nil {
		Log(LVL_ERROR, "can't request Hy flight data with request ", body, ":", err)
		return BSP{}, false
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		Log(LVL_ERROR, "can't request Hy flight data with request ", body, ":", err)
		return BSP{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return BSP{}, false
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log(LVL_ERROR, "can't read responce with request ", body, ":", err)
		return BSP{}, false
	}
	err = json.Unmarshal(data, &dat)
	if err != nil {
		Log(LVL_ERROR, "can't unmarshal hy data for ship ", shipID, err)
		return BSP{}, false
	}
	if dat.HyStatus != "freight" && !DEFVAL.DebugMode {
		Log(LVL_ERROR, "ship is not freight! ", shipID, err)
		return BSP{}, false
	}

	return dat, true
}
