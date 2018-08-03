package commons

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
)

const (
	EMI_DMG_MECH       = "1"
	EMI_DMG_HEAT       = "2"
	EMI_DMG_RADI       = "3"
	EMI_DMG_GRAVI      = "4"
	EMI_VEL_UP         = "5"
	EMI_VEL_DOWN       = "6"
	EMI_DIST_UP        = "7"
	EMI_DIST_DOWN      = "8"
	EMI_WARP_TURN      = "9"
	EMI_ACCEL          = "10"
	EMI_REVERSE        = "11"
	EMI_ENGINE_HEAT    = "12"
	EMI_FUEL           = "13"
	EMI_TURN           = "14"
	EMI_STRAFE         = "15"
	EMI_SCAN_RADIUS    = "16"
	EMI_SCAN_SPEED     = "17"
	EMI_DROP_RADIUS    = "18"
	EMI_DROP_SPEED     = "19"
	EMI_RADAR_COSMOS   = "20"
	EMI_RADAR_WARP     = "21"
	EMI_RADAR_ANG_UP   = "22"
	EMI_RADAR_ANG_DOWN = "23"
	EMI_CO2            = "24"
	EMI_DEF_HEAT       = "25"
	EMI_DEF_RADI       = "26"
	EMI_DEF_MECH       = "27"
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
		if len(point.Emissions) == 0 {
			continue
		}

		dist = point.Pos.Sub(ship).Len()
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
