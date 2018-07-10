package main

import (
	"encoding/json"
	."github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/static"
	."github.com/Shnifer/magellan/log"
)

//COPYPASTE Server.OutConns
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

func loadNewGalaxy(GalaxyID string) {
	CurGalaxy = loadGalaxyState(GalaxyID)
}