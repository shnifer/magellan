package commons

const (
	SYS_MARCH = iota
	SYS_SHUNTER
	SYS_WARP
	SYS_SHIELD
	SYS_RADAR
	SYS_SCANNER
	SYS_FUEL
	SYS_LSS
)

type BSPParams struct {
	//0...100
	March_engine struct {
		Thrust_max   float64 `json:"thrust"`
		Thrust_acc   float64 `json:"accel"`
		Thrust_slow  float64 `json:"slowdown"`
		Reverse_max  float64 `json:"thrust_rev"`
		Reverse_acc  float64 `json:"accel_rev"`
		Reverse_slow float64 `json:"slowdown_rev"`
		Heat_prod    float64 `json:"heat_prod"`
		AZ           float64 `json:"az_level"`
		Volume       float64 `json:"volume"`
	} `json:"march_engine"`

	Warp_engine struct {
		Distort_max            float64 `json:"distort"`
		Distort_acc            float64 `json:"distort_acc"`
		Distort_slow           float64 `json:"distort_slow"`
		Consumption            float64 `json:"consumption"`
		Warp_enter_consumption float64 `json:"warp_enter_consumption"`
		Turn_speed             float64 `json:"turn_speed"`
		Turn_consumption       float64 `json:"turn_consumption"`
		AZ                     float64 `json:"az_level"`
		Volume                 float64 `json:"volume"`
	} `json:"warp_engine"`

	Shunter struct {
		Turn_max    float64 `json:"turn"`
		Turn_acc    float64 `json:"turn_acc"`
		Turn_slow   float64 `json:"turn_slow"`
		Strafe_max  float64 `json:"strafe"`
		Strafe_acc  float64 `json:"strafe_acc"`
		Strafe_slow float64 `json:"strafe_slow"`
		Heat_prod   float64 `json:"heat_prod"`
		AZ          float64 `json:"az_level"`
		Volume      float64 `json:"volume"`
	} `json:"shunter"`

	Radar struct {
		Range_Max    float64 `json:"range_max"`
		Angle_Min    float64 `json:"angle_min"`
		Angle_Max    float64 `json:"angle_max"`
		Angle_Change float64 `json:"angle_change"`
		Range_Change float64 `json:"range_change"`
		Rotate_Speed float64 `json:"rotate_speed"`
		AZ           float64 `json:"az_level"`
		Volume       float64 `json:"volume"`
	} `json:"radar"`

	Scanner struct {
		DropRange float64 `json:"drop_range"`
		DropSpeed float64 `json:"drop_speed"`
		ScanRange float64 `json:"scan_range"`
		ScanSpeed float64 `json:"scan_speed"`
		AZ        float64 `json:"az_level"`
		Volume    float64 `json:"volume"`
	} `json:"scaner"`

	Fuel_tank struct {
		Fuel_volume     float64 `json:"fuel_volume"`
		Fuel_Protection float64 `json:"fuel_protection"`
		Radiation_def   float64 `json:"radiation_def"`
		AZ              float64 `json:"az_level"`
		Volume          float64 `json:"volume"`
	} `json:"fuel_tank"`

	Lss struct {
		Thermal_def       float64 `json:"thermal_def"`
		Co2_level         float64 `json:"co2_level"`
		Air_volume        float64 `json:"air_volume"`
		Air_prepare_speed float64 `json:"air_speed"`
		Lightness         float64 `json:"lightness"`
		AZ                float64 `json:"az_level"`
		Volume            float64 `json:"volume"`
	} `json:"lss"`

	Shields struct {
		Radiation_def   float64 `json:"radiation_def"`
		Disinfect_level float64 `json:"disinfect_level"`
		Mechanical_def  float64 `json:"mechanical_def"`
		Heat_reflection float64 `json:"heat_reflection"`
		Heat_capacity   float64 `json:"heat_capacity"`
		Heat_sink       float64 `json:"heat_sink"`
		AZ              float64 `json:"az_level"`
		Volume          float64 `json:"volume"`
	} `json:"shields"`
}

type BSPCargo struct {
	Beacons struct {
		Mass  float64 `json:"weight"`
		Count int     `json:"amount"`
	} `json:"beacons"`
	Mines []struct {
		Mass  float64 `json:"weight"`
		Owner string  `json:"company"`
	} `json:"mines"`
	Modules []struct {
		Mass   float64 `json:"weight"`
		Owner  string  `json:"company"`
		Planet string  `json:"planet_id"`
	} `json:"modules"`
}

//Rework CalcDegrade on change
type BSP struct {
	FlightID int `json:"flight_id"`
	Dock     int `json:"dock"`

	Ship struct {
		Name      string  `json:"name"` //human name
		NodesMass float64 `json:"nodes_weight"`
	} `json:"ship"`

	BSPParams `json:"params"`

	BSPCargo `json:"cargo"`

	KnownMinerals []struct {
		ID       string `json:"id"`
		UserName string `json:"name"`
	} `json:"known_minerals"`
}

func (base *BSP) CalcDegrade(degrade BSPDegrade) (result *BSP) {
	if base == nil {
		return &BSP{}
	}
	res := *base

	res.March_engine.Thrust_max *= degrade.March_engine.Thrust_max
	res.March_engine.Thrust_acc *= degrade.March_engine.Thrust_acc
	res.March_engine.Thrust_slow *= degrade.March_engine.Thrust_slow
	res.March_engine.Reverse_max *= degrade.March_engine.Reverse_max
	res.March_engine.Reverse_acc *= degrade.March_engine.Reverse_acc
	res.March_engine.Reverse_slow *= degrade.March_engine.Reverse_slow
	res.March_engine.Heat_prod *= degrade.March_engine.Heat_prod

	res.Shunter.Turn_max *= degrade.Shunter.Turn_max
	res.Shunter.Turn_acc *= degrade.Shunter.Turn_acc
	res.Shunter.Turn_slow *= degrade.Shunter.Turn_slow
	res.Shunter.Strafe_max *= degrade.Shunter.Strafe_max
	res.Shunter.Strafe_acc *= degrade.Shunter.Strafe_acc
	res.Shunter.Strafe_slow *= degrade.Shunter.Strafe_slow
	res.Shunter.Heat_prod *= degrade.Shunter.Heat_prod

	res.Warp_engine.Distort_max *= degrade.Warp_engine.Distort_max
	res.Warp_engine.Distort_acc *= degrade.Warp_engine.Distort_acc
	res.Warp_engine.Distort_slow *= degrade.Warp_engine.Distort_slow
	res.Warp_engine.Turn_speed *= degrade.Warp_engine.Turn_speed
	res.Warp_engine.Turn_consumption *= degrade.Warp_engine.Turn_consumption
	res.Warp_engine.Warp_enter_consumption *= degrade.Warp_engine.Warp_enter_consumption
	res.Warp_engine.Consumption *= degrade.Warp_engine.Consumption

	res.Shields.Heat_sink *= degrade.Shields.Heat_sink
	res.Shields.Heat_reflection *= degrade.Shields.Heat_reflection
	res.Shields.Heat_capacity *= degrade.Shields.Heat_capacity
	res.Shields.Disinfect_level *= degrade.Shields.Disinfect_level
	res.Shields.Mechanical_def *= degrade.Shields.Mechanical_def
	res.Shields.Radiation_def *= degrade.Shields.Radiation_def

	res.Fuel_tank.Fuel_Protection *= degrade.Fuel_tank.Fuel_Protection

	res.Scanner.DropRange *= degrade.Scanner.DropRange
	res.Scanner.DropSpeed *= degrade.Scanner.DropSpeed
	res.Scanner.ScanRange *= degrade.Scanner.ScanRange
	res.Scanner.ScanSpeed *= degrade.Scanner.ScanSpeed

	res.Radar.Angle_Change *= degrade.Radar.Angle_Change
	res.Radar.Angle_Min *= degrade.Radar.Angle_Min
	res.Radar.Angle_Max *= degrade.Radar.Angle_Max
	res.Radar.Range_Max *= degrade.Radar.Range_Max
	res.Radar.Range_Change *= degrade.Radar.Range_Change
	res.Radar.Rotate_Speed *= degrade.Radar.Rotate_Speed

	res.Lss.Thermal_def *= degrade.Lss.Thermal_def
	res.Lss.Air_prepare_speed *= degrade.Lss.Air_prepare_speed

	return &res
}

func emptyDegrade() BSPDegrade {
	var res BSPDegrade

	res.March_engine.Thrust_max = 1
	res.March_engine.Thrust_acc = 1
	res.March_engine.Thrust_slow = 1
	res.March_engine.Reverse_max = 1
	res.March_engine.Reverse_acc = 1
	res.March_engine.Reverse_slow = 1
	res.March_engine.Heat_prod = 1

	res.Shunter.Turn_max = 1
	res.Shunter.Turn_acc = 1
	res.Shunter.Turn_slow = 1
	res.Shunter.Strafe_max = 1
	res.Shunter.Strafe_acc = 1
	res.Shunter.Strafe_slow = 1
	res.Shunter.Heat_prod = 1

	res.Warp_engine.Distort_max = 1
	res.Warp_engine.Distort_acc = 1
	res.Warp_engine.Distort_slow = 1
	res.Warp_engine.Turn_speed = 1
	res.Warp_engine.Turn_consumption = 1
	res.Warp_engine.Warp_enter_consumption = 1
	res.Warp_engine.Consumption = 1

	res.Shields.Heat_sink = 1
	res.Shields.Heat_reflection = 1
	res.Shields.Heat_capacity = 1
	res.Shields.Disinfect_level = 1
	res.Shields.Mechanical_def = 1
	res.Shields.Radiation_def = 1

	res.Fuel_tank.Fuel_Protection = 1

	res.Scanner.DropRange = 1
	res.Scanner.DropSpeed = 1
	res.Scanner.ScanRange = 1
	res.Scanner.ScanSpeed = 1

	res.Radar.Angle_Change = 1
	res.Radar.Angle_Min = 1
	res.Radar.Angle_Max = 1
	res.Radar.Range_Max = 1
	res.Radar.Range_Change = 1
	res.Radar.Rotate_Speed = 1

	res.Lss.Thermal_def = 1
	res.Lss.Air_prepare_speed = 1
	res.Lss.Co2_level = 0

	return res
}
