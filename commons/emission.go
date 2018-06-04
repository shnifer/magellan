package commons

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
)

const (
	EMISSION_HEAT = "Heat"
	EMISSION_SLOW = "Slow"
	EMISSION_DMG  = "Dmg"
)

type Emission struct {
	Type      string
	MainRange float64 `json:",omitempty"`
	MainValue float64 `json:",omitempty"`
	FarRange  float64 `json:",omitempty"`
	FarValue  float64 `json:",omitempty"`
}

func CalculateEmissions(galaxy *Galaxy, ship v2.V2) (res map[string]float64) {
	defer LogFunc("commonRecv")()

	res = make(map[string]float64)

	var dist float64

	for _, point := range galaxy.Points {
		if len(point.Emissions) > 0 {
			dist = point.Pos.Sub(ship).Len()
		}
		for _, emission := range point.Emissions {
			v := 0.0
			if dist < emission.MainRange {
				v = emission.MainValue
			} else if dist < emission.FarRange {
				v = calcMiddle(dist, emission.MainRange, emission.FarRange, emission.MainValue, emission.FarValue)
			}
			if v != 0 {
				res[emission.Type] += v
			}
		}
	}

	return res
}

func calcMiddle(x, min, max, vMin, vMax float64) float64 {
	k := (x - min) / (max - min)
	return vMin + (vMax-vMin)*k
}
