package commons

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
)

type StateData struct {
	//ServerTime time.Time
	BSP    *BSP
	Galaxy *Galaxy
	//map[fullKey]Building
	Buildings map[string]Building
}

func (sd StateData) Encode() []byte {
	buf, err := json.Marshal(sd)
	if err != nil {
		Log(LVL_ERROR, "can't marshal stateData", err)
		return nil
	}
	return buf
}

func (StateData) Decode(buf []byte) (sd StateData, err error) {
	err = json.Unmarshal(buf, &sd)
	if err != nil {
		return StateData{}, err
	}
	sd.Galaxy.RecalcLvls()
	for _, b := range sd.Buildings {
		if b.GalaxyID != WARP_Galaxy_ID {
			sd.Galaxy.AddBuilding(b)
		} else {
			sd.Galaxy.AddWarpBuilding(b)
		}
	}
	return sd, nil
}

func (sd StateData) Copy() (res StateData) {
	res = sd
	if sd.BSP != nil {
		val := *sd.BSP
		res.BSP = &val
	}
	if sd.Galaxy != nil {

		val := *sd.Galaxy
		val.Points = make(map[string]*GalaxyPoint, len(sd.Galaxy.Points))
		val.Ordered = nil
		for k, v := range sd.Galaxy.Points {
			val.Points[k] = v
		}
		val.RecalcLvls()
		res.Galaxy = &val
	}

	return res
}
