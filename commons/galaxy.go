package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

const (
	//GalaxyPoint.Type
	//also includes outerBuilds (BUILDING_BLACKBOX, BUILDING_BEACON)
	GPT_STAR       = "STAR"
	GPT_WARP       = "WARP"
	GPT_HARDPLANET = "HARDPLANET"
	GPT_GASPLANET  = "GASPLANET"
	GPT_ASTEROID   = "ASTEROID"
)

type Galaxy struct {
	//for systems - range of "system borders"
	SpawnDistance float64

	Points map[string]*GalaxyPoint

	//recalculated on Decode
	Ordered []*GalaxyPoint `json:"-"`
	maxLvl  int

	//used by update
	fixedTimeRest float64
}

type GalaxyPoint struct {
	//Id setted on load from file
	ID        string `json:"id,omitempty"`
	ParentID  string `json:"pid,omitempty"`
	IsVirtual bool   `json:"iv,omitempty"`

	//found on recalc
	//phys order level
	Level int `json:"lv,omitempty"`
	//graph order level, ignore
	GLevel int `json:"gl,omitempty"`

	Pos v2.V2

	Orbit  float64 `json:"orb,omitempty"`
	Period float64 `json:"per,omitempty"`

	Type     string  `json:"t,omitempty"`
	SpriteAN string  `json:"sp,omitempty"`
	Size     float64 `json:"s,omitempty"`

	Mass   float64 `json:"m,omitempty"`
	GDepth float64 `json:"gd,omitempty"`

	//for warp points
	WarpSpawnDistance float64 `json:"wsd,omitempty"`
	WarpYellowOutDist float64 `json:"wyo,omitempty"`
	WarpGreenOutDist float64 `json:"wgo,omitempty"`
	WarpGreenInDist float64 `json:"wgi,omitempty"`
	WarpRedOutDist float64 `json:"wro,omitempty"`

	ScanData string `json:"sd,omitempty"`

	Minerals   []string    `json:"mi,omitempty"`
	Emissions  []Emission  `json:"emm,omitempty"`
	Signatures []Signature `json:"sig,omitempty"`
	Color      color.RGBA  `json:"clr"`

	//updated on Decode or add|del building
	//map[ownerName]fullkey
	Mines      map[string]string `json:"mns,omitempty"`
	FishHouses map[string]string `json:"fhs,omitempty"`

	//for warp points
	//map[fullKey]message
	Beacons    map[string]string `json:"bcs,omitempty"`
	BlackBoxes map[string]string `json:"bbs,omitempty"`

	/*	HasMine     bool   `json:"hm,omitempty"`
		MineOwner   string `json:"mo,omitempty"`
		MineFullKey string `json:"mk,omitempty"`

		HasFishHouse     bool   `json:"fm,omitempty"`
		FishHouseOwner   string `json:"fo,omitempty"`
		FishHouseFullKey string `json:"fk,omitempty"`
	*/
}

func (gp GalaxyPoint) MarshalJSON() ([]byte, error) {
	//Marshal just as standard
	//to avoid recursive GalaxyPoint.MarshalJSON()
	type just GalaxyPoint
	buf, err := json.Marshal(just(gp))

	if err != nil {
		return buf, err
	}
	buf = bytes.Replace(buf, []byte(`"Pos":{},`), []byte{}, -1)
	return buf, nil
}
