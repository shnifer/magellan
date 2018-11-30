package commons

import (
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/v2"
)

const (
	EMI_DMG_MECH       = "1"  //in Engi Counters
	EMI_DMG_HEAT       = "2"  //in Engi Counters
	EMI_DMG_RADI       = "3"  //in Engi Counters
	EMI_DMG_GRAVI      = "4"  //in pilot warpscene -> engi command GDmgHard GDmgMedium
	EMI_VEL_UP         = "5"  //Pilot cosmoscene procEmi
	EMI_VEL_DOWN       = "6"  //Pilot cosmoscene procEmi
	EMI_DIST_UP        = "7"  //in Engi into BSPDegrade
	EMI_DIST_DOWN      = "8"  //in Engi into BSPDegrade
	EMI_WARP_TURN      = "9"  //in Engi into BSPDegrade
	EMI_ACCEL          = "10" //in Engi into BSPDegrade
	EMI_REVERSE        = "11" //in Engi into BSPDegrade
	EMI_ENGINE_HEAT    = "12" //in Engi into BSPDegrade
	EMI_FUEL           = "13" //in Engi into BSPDegrade
	EMI_TURN           = "14" //in Engi into BSPDegrade
	EMI_STRAFE         = "15" //in Engi into BSPDegrade
	EMI_SCAN_RADIUS    = "16" //in Engi into BSPDegrade
	EMI_SCAN_SPEED     = "17" //in Engi into BSPDegrade
	EMI_DROP_RADIUS    = "18" //in Engi into BSPDegrade
	EMI_DROP_SPEED     = "19" //in Engi into BSPDegrade
	EMI_RADAR_COSMOS   = "20" //in Engi into BSPDegrade
	EMI_RADAR_WARP     = "21" //in Engi into BSPDegrade
	EMI_RADAR_ANG_UP   = "22" //in Engi into BSPDegrade
	EMI_RADAR_ANG_DOWN = "23" //in Engi into BSPDegrade
	EMI_CO2            = "24" //in Engi into BSPDegrade
	EMI_DEF_HEAT       = "25" //in Engi into BSPDegrade
	EMI_DEF_RADI       = "26" //in Engi into BSPDegrade
	EMI_DEF_MECH       = "27" //in Engi into BSPDegrade
	EMI_WORMHOLE       = "28" //Engi update checks
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

	for _, point := range galaxy.Ordered {
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

/*Повреждение (Механическое)
Повреждение (Термическое)
Повреждение (Радиационное)
Повреждение (Гравитационное)
Увеличение скорости
Снижение скорости
Увеличение искривления (в варпе)
Снижение искривления (в варпе)
Уменьшение скорости поворота (в варпе)
Снижение маршевой тяги
Снижение реверсивной тяги
Большее тепловыделение движков (ухудшение теплоотвода)
Увеличение расхода топлива (распад топлива)
Уменьшение скорости поворота
Уменьшение тяги стрейфа
Уменьшение радиусов сканирования
Уменьшение скорости сканирования
Уменьшение радиусов сброса
Уменьшение скорости сброса
Уменьшение дальности радара
Уменьшение дальности радара (в варпе)
Увеличение угла (в варпе)
Уменьшение угла (в варпе)
Ухудшение работы СЖО (повышение уровня СО2)
Снижение тепловой защиты (степени отражения)
Снижение радиационной защиты
Снижение механической защиты*/
