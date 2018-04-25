package commons

import (
	"encoding/json"
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

type StateData struct {
	//ServerTime time.Time
	BSP    *BSP
	Galaxy *Galaxy
}

type BSP struct {
	//0...100
	Thrust,
	Thrust_rev,
	Thrust_acc,
	Thrust_rev_acc,
	Thrust_slow,
	Thrust_rev_slow,
	Thrust_heat_capacity,
	Thrust_heat_prod,
	Thrust_heat_sink float64

	Distort_level,
	Warp_enter_consumption,
	Distort_level_acc,
	Distort_level_slow,
	Distort_consumption,
	Distort_turn,
	Distort_turn_consumption float64

	Turn_max,
	Strafe_max,
	Turn_acc,
	Strafe_acc,
	Turn_slow,
	Strafe_slow,
	Maneur_heat_capacity,
	Maneur_heat_prod,
	Maneur_heat_sink float64

	Radar_range_min,
	Radar_range_max,
	Radar_angle_min,
	Radar_angle_max,
	Radar_angle_change,
	Radar_range_change,
	Scan_range,
	Scan_speed float64

	Sonar_range_min,
	Sonar_range_max,
	Sonar_angle_min,
	Sonar_angle_max,
	Sonar_angle_change,
	Sonar_range_change,
	Sonar_rotate_speed float64
}

type Galaxy struct {
	//for systems - range of "system borders"
	SpawnDistance float64

	Points []GalaxyPoint
}

type GalaxyPoint struct {
	ID       string
	ParentID string

	Pos v2.V2

	Orbit  float64
	Period float64

	Type  string
	Size  float64
	Color color.RGBA

	Mass float64

	//for warp points
	WarpSpawnDistance float64
	WarpInDistance    float64

	ScanData string
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
		res.Galaxy = &val
	}

	return res
}
