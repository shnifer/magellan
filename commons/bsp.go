package commons

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
	} `json:"shunter"`

	Radar struct {
		Range_Max    float64 `json:"range_max"`
		Angle_Min    float64 `json:"angle_min"`
		Angle_Max    float64 `json:"angle_max"`
		Angle_Change float64 `json:"angle_change"`
		Range_Change float64 `json:"range_change"`
		Rotate_Speed float64 `json:"rotate_speed"`
		AZ           float64 `json:"az_level"`
	} `json:"radar"`

	Scanner struct {
		DropRange float64 `json:"drop_range"`
		DropSpeed float64 `json:"drop_speed"`
		ScanRange float64 `json:"scan_range"`
		ScanSpeed float64 `json:"scan_speed"`
		AZ        float64 `json:"az_level"`
	} `json:"scaner"`

	Fuel_tank struct {
		Fuel_volume   float64 `json:"fuel_volume"`
		Compact       float64 `json:"compact"`
		Radiation_def float64 `json:"radiation_def"`
		AZ            float64 `json:"az_level"`
	} `json:"fuel_tank"`

	Lss struct {
		Thermal_def       float64 `json:"thermal_def"`
		Co2_level         float64 `json:"co2_level"`
		Air_volume        float64 `json:"air_volume"`
		Air_prepare_speed float64 `json:"air_speed"`
		Lightness         float64 `json:"lightness"`
		AZ                float64 `json:"az_level"`
	} `json:"lss"`

	Shields struct {
		Radiation_def   float64 `json:"radiation_def"`
		Disinfect_level float64 `json:"disinfect_level"`
		Mechanical_def  float64 `json:"mechanical_def"`
		Heat_reflection float64 `json:"heat_reflection"`
		Heat_capacity   float64 `json:"heat_capacity"`
		Heat_sink       float64 `json:"heat_sink"`
		AZ              float64 `json:"az_level"`
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
		ID       int    `json:"id"`
		UserName string `json:"name"`
	} `json:"known_minerals"`
}

func (base *BSP) CalcDegrade(degrade BSPDegrade) (result *BSP) {
	if base == nil {
		return &BSP{}
	}
	res := *base

	//todo: calc degrade

	return &res
}
