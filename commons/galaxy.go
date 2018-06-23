package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

const (
	//GalaxyPoint.Type
	//also includes outerBuilds (BUILDING_BLACKBOX, BUILDING_BEACKON)
	GPT_STAR = "STAR"
	GPT_WARP = "WARP"
	GPS_PLANET = "PLANET"
	GPS_ASTEROID = "ASTEROID"
	GPS_SATELLITE = "SATELLITE"
)

type Galaxy struct {
	//for systems - range of "system borders"
	SpawnDistance float64

	Points map[string]*GalaxyPoint

	//recalculated on Decode
	Ordered []*GalaxyPoint `json:"-"`
	maxLvl  int
}

type GalaxyPoint struct {
	ID       string `json:"id,omitempty"`
	ParentID string `json:"pid,omitempty"`

	Pos v2.V2

	Orbit  float64 `json:"orb,omitempty"`
	Period float64 `json:"per,omitempty"`


	Type string  `json:"t,omitempty"`
	SpriteAN string `json:"sp,omitempty"`
	Size float64 `json:"s,omitempty"`

	Mass   float64 `json:"m,omitempty"`
	GDepth float64 `json:"gd,omitempty"`

	//for warp points
	WarpSpawnDistance float64 `json:"wsd,omitempty"`
	WarpInDistance    float64 `json:"did,omitempty"`

	ScanData string `json:"sd,omitempty"`

	Emissions  []Emission  `json:"emm,omitempty"`
	Signatures []Signature `json:"sig,omitempty"`
	Color      color.RGBA  `json:"clr"`

	HasMine     bool `json:"hm,omitempty"`
	MineOwner   string `json:"mo,omitempty"`
	MineFullKey string `json:"mk,omitempty"`

	HasFishHouse     bool `json:"fm,omitempty"`
	FishHouseOwner   string `json:"fo,omitempty"`
	FishHouseFullKey string `json:"fk,omitempty"`
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
