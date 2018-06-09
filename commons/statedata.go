package commons

import (
	"bytes"
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

type SpaceObjects struct {
}

type StateData struct {
	//ServerTime time.Time
	BSP    *BSP
	Galaxy *Galaxy
	//map[fullKey]Building
	Buildings map[string]Building
}

//Rework CalcDegrade on change
type BSP struct {
	ShipName string //human name

	Mass float64
	//0...100
	March_engine struct {
		Thrust,
		Thrust_acc,
		Thrust_slow,
		Thrust_rev,
		Thrust_rev_acc,
		Thrust_rev_slow,
		Heat_prod float64
	}

	Warp_engine struct {
		Distort_level,
		Distort_level_acc,
		Distort_level_slow,
		Distort_consumption,
		Warp_enter_consumption,
		Turn_speed,
		Distort_turn_consumption float64
	}

	Shunter struct {
		Turn_max,
		Turn_acc,
		Turn_slow,
		Strafe_max,
		Strafe_acc,
		Strafe_slow,
		Heat_prod float64
	}

	Radar struct {
		Range,
		Angle,
		Scan_range,
		Scan_speed float64
	}

	Sonar struct {
		Range_max,
		Angle_min,
		Angle_max,
		Angle_change,
		Range_change,
		Rotate_speed float64
	}

	Fuel_tank struct {
		Fuel_volume,
		Compact,
		Radiation_def float64
	}

	Lss struct {
		Thermal_def,
		Co2_level,
		Air_volume,
		Air_prepare_speed,
		Lightness float64
	}

	Shields struct {
		Radiation_def,
		Disinfect_level,
		Mechanical_def,
		Heat_reflection,
		Heat_capacity,
		Heat_sink float64
	}
}

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
	Size float64 `json:"s,omitempty"`

	Mass float64 `json:"m,omitempty"`

	//for warp points
	WarpSpawnDistance float64 `json:"wsd,omitempty"`
	WarpInDistance    float64 `json:"did,omitempty"`

	ScanData string `json:"sd,omitempty"`

	Emissions  []Emission  `json:"emm,omitempty"`
	Signatures []Signature `json:"sig,omitempty"`
	Color      color.RGBA  `json:"clr"`

	HasMine   bool
	MineOwner string
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
		sd.Galaxy.addBuilding(b)
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

func (base *BSP) CalcDegrade(degrade BSPDegrade) (result *BSP) {
	if base == nil {
		return &BSP{}
	}
	res := *base

	res.March_engine.Thrust = base.March_engine.Thrust * (1 - degrade.Thrust)
	res.March_engine.Thrust_rev = base.March_engine.Thrust_rev * (1 - degrade.Thrust_rev)
	res.March_engine.Thrust_acc = base.March_engine.Thrust_acc * (1 - degrade.Thrust_acc)
	res.March_engine.Thrust_rev_acc = base.March_engine.Thrust_rev_acc * (1 - degrade.Thrust_rev_acc)
	res.March_engine.Thrust_slow = base.March_engine.Thrust_slow * (1 - degrade.Thrust_slow)
	res.March_engine.Thrust_rev_slow = base.March_engine.Thrust_rev_slow * (1 - degrade.Thrust_rev_slow)
	res.March_engine.Heat_prod = base.March_engine.Heat_prod * (1 + degrade.Thrust_heat_prod)

	res.Warp_engine.Distort_level = base.Warp_engine.Distort_level * (1 - degrade.Distort_level)
	res.Warp_engine.Warp_enter_consumption = base.Warp_engine.Warp_enter_consumption * (1 + degrade.Warp_enter_consumption)
	res.Warp_engine.Distort_level_acc = base.Warp_engine.Distort_level_acc * (1 - degrade.Distort_level_acc)
	res.Warp_engine.Distort_level_slow = base.Warp_engine.Distort_level_slow * (1 - degrade.Distort_level_slow)
	res.Warp_engine.Distort_consumption = base.Warp_engine.Distort_consumption * (1 + degrade.Distort_consumption)
	res.Warp_engine.Turn_speed = base.Warp_engine.Turn_speed * (1 - degrade.Distort_turn)
	res.Warp_engine.Distort_turn_consumption = base.Warp_engine.Distort_turn_consumption * (1 + degrade.Distort_turn_consumption)

	res.Shunter.Turn_max = base.Shunter.Turn_max * (1 - degrade.Turn_max)
	res.Shunter.Turn_acc = base.Shunter.Turn_acc * (1 - degrade.Turn_acc)
	res.Shunter.Turn_slow = base.Shunter.Turn_slow * (1 - degrade.Turn_slow)
	res.Shunter.Strafe_max = base.Shunter.Strafe_max * (1 - degrade.Strafe_max)
	res.Shunter.Strafe_acc = base.Shunter.Strafe_acc * (1 - degrade.Strafe_acc)
	res.Shunter.Strafe_slow = base.Shunter.Strafe_slow * (1 - degrade.Strafe_slow)
	res.Shunter.Heat_prod = base.Shunter.Heat_prod * (1 + degrade.Maneur_heat_prod)

	res.Radar.Range = base.Radar.Range * (1 - degrade.Radar_range)
	res.Radar.Angle = base.Radar.Angle * (1 - degrade.Radar_angle)
	res.Radar.Scan_range = base.Radar.Scan_range * (1 - degrade.Scan_range)
	res.Radar.Scan_speed = base.Radar.Scan_speed * (1 - degrade.Scan_speed)

	res.Sonar.Range_max = base.Sonar.Range_max * (1 - degrade.Sonar_range_max)
	midAng := (base.Sonar.Angle_max + base.Sonar.Angle_min) / 2
	dAng := (base.Sonar.Angle_max - base.Sonar.Angle_min) / 2
	res.Sonar.Angle_min = midAng - dAng*(1-degrade.Sonar_angle_min)
	res.Sonar.Angle_max = midAng + dAng*(1-degrade.Sonar_angle_max)
	res.Sonar.Angle_change = base.Sonar.Angle_change * (1 - degrade.Sonar_angle_change)
	res.Sonar.Range_change = base.Sonar.Range_change * (1 - degrade.Sonar_range_change)
	res.Sonar.Rotate_speed = base.Sonar.Rotate_speed * (1 - degrade.Sonar_rotate_speed)

	return &res
}
