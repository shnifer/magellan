package commons

import (
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
	"time"
)

type StateData struct {
	ServerTime time.Time
	BSP        *BSP
	Galaxy     *Galaxy
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
}

type Galaxy struct {
	Points []GalaxyPoint
}

type GalaxyPoint struct {
	ID       string
	ParentID string

	Pos graph.Point

	Orbit    float64
	Period   float64
	DegStart float64

	Type string
	Size float64
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
