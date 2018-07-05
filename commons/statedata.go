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

//Rework CalcDegrade on change
type BSP struct {
	ShipName   string  `json:"ship_name"` //human name
	FlightCorp string  `json:"corp_name"`
	Mass       float64 `json:"mass"`

	KnownMinerals []string `json:"known_minerals"`

	//start drop items, current is stored in NaviData
	BeaconCount int `json:"beacon_count"`
	//[]corpName, i.e. ["gd","gd","pre"]
	//[]planetName, i.e. ["CV8-85","RD4-42-13"]
	Mines   []string `json:"mines"`
	Landing []string `json:"landing"`

	MineMass    float64
	LandingMass float64
	BeaconMass  float64

	//0...100
	March_engine struct {
		Thrust_max   float64 `json:"thrust_max"`
		Thrust_acc   float64 `json:"thrust_acc"`
		Thrust_slow  float64 `json:"thrust_slow"`
		Reverse_max  float64 `json:"reverse_max"`
		Reverse_acc  float64 `json:"reverse_acc"`
		Reverse_slow float64 `json:"reverse_slow"`
		Heat_prod    float64 `json:"heat_prod"`
	} `json:"march_engine"`

	Warp_engine struct {
		Distort_max            float64 `json:"distort_max"`
		Distort_acc            float64 `json:"distort_acc"`
		Distort_slow           float64 `json:"distort_slow"`
		Consumption            float64 `json:"consumption"`
		Warp_enter_consumption float64 `json:"warp_enter_consumption"`
		Turn_speed             float64 `json:"turn_speed"`
		Turn_consumption       float64 `json:"turn_consumption"`
	} `json:"warp_engine"`

	Shunter struct {
		Turn_max    float64 `json:"turn_max"`
		Turn_acc    float64 `json:"turn_acc"`
		Turn_slow   float64 `json:"turn_slow"`
		Strafe_max  float64 `json:"strafe_max"`
		Strafe_acc  float64 `json:"strafe_acc"`
		Strafe_slow float64 `json:"strafe_slow"`
		Heat_prod   float64 `json:"heat_prod"`
	} `json:"shunter"`

	Radar struct {
		Range      float64 `json:"range"`
		Angle      float64 `json:"angle"`
		Scan_range float64 `json:"scan_range"`
		Scan_speed float64 `json:"scan_speed"`
	} `json:"radar"`

	Sonar struct {
		Range_max    float64 `json:"range_max"`
		Angle_min    float64 `json:"angle_min"`
		Angle_max    float64 `json:"angle_max"`
		Angle_change float64 `json:"angle_change"`
		Range_change float64 `json:"range_change"`
		Rotate_speed float64 `json:"rotate_speed"`
	} `json:"sonar"`

	Fuel_tank struct {
		Fuel_volume   float64 `json:"fuel_volume"`
		Compact       float64 `json:"compact"`
		Radiation_def float64 `json:"radiation_def"`
	} `json:"fuel_tank"`

	Lss struct {
		Thermal_def       float64 `json:"thermal_def"`
		Co2_level         float64 `json:"co2_level"`
		Air_volume        float64 `json:"air_volume"`
		Air_prepare_speed float64 `json:"air_prepare_speed"`
		Lightness         float64 `json:"lightness"`
	} `json:"lss"`

	Shields struct {
		Radiation_def   float64 `json:"radiation_def"`
		Disinfect_level float64 `json:"disinfect_level"`
		Mechanical_def  float64 `json:"mechanical_def"`
		Heat_reflection float64 `json:"heat_reflection"`
		Heat_capacity   float64 `json:"heat_capacity"`
		Heat_sink       float64 `json:"heat_sink"`
	} `json:"shields"`
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
		sd.Galaxy.AddBuilding(b)
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

func (base *BSP) CalcDegrade(degrade BSPDegrade) (result *BSP) {
	if base == nil {
		return &BSP{}
	}
	res := *base

	res.March_engine.Thrust_max = base.March_engine.Thrust_max * (1 - degrade.Thrust)
	res.March_engine.Reverse_max = base.March_engine.Reverse_max * (1 - degrade.Thrust_rev)
	res.March_engine.Thrust_acc = base.March_engine.Thrust_acc * (1 - degrade.Thrust_acc)
	res.March_engine.Reverse_acc = base.March_engine.Reverse_acc * (1 - degrade.Thrust_rev_acc)
	res.March_engine.Thrust_slow = base.March_engine.Thrust_slow * (1 - degrade.Thrust_slow)
	res.March_engine.Reverse_slow = base.March_engine.Reverse_slow * (1 - degrade.Thrust_rev_slow)
	res.March_engine.Heat_prod = base.March_engine.Heat_prod * (1 + degrade.Thrust_heat_prod)

	res.Warp_engine.Distort_max = base.Warp_engine.Distort_max * (1 - degrade.Distort_level)
	res.Warp_engine.Warp_enter_consumption = base.Warp_engine.Warp_enter_consumption * (1 + degrade.Warp_enter_consumption)
	res.Warp_engine.Distort_acc = base.Warp_engine.Distort_acc * (1 - degrade.Distort_level_acc)
	res.Warp_engine.Distort_slow = base.Warp_engine.Distort_slow * (1 - degrade.Distort_level_slow)
	res.Warp_engine.Consumption = base.Warp_engine.Consumption * (1 + degrade.Distort_consumption)
	res.Warp_engine.Turn_speed = base.Warp_engine.Turn_speed * (1 - degrade.Distort_turn)
	res.Warp_engine.Turn_consumption = base.Warp_engine.Turn_consumption * (1 + degrade.Distort_turn_consumption)

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
