package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/v2"
	"image/color"
	"reflect"
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

	Points  map[string]*GalaxyPoint
	Ordered []*GalaxyPoint `json:"-"`
	maxLvl  int
	//changed on ordered slice
	//lvlLists [][]string
}

type GalaxyPoint struct {
	ID       string `json:",omitempty"`
	ParentID string `json:",omitempty"`

	Pos v2.V2

	Orbit  float64 `json:",omitempty"`
	Period float64 `json:",omitempty"`

	Type string
	Size float64

	Mass float64 `json:",omitempty"`

	//for warp points
	WarpSpawnDistance float64 `json:",omitempty"`
	WarpInDistance    float64 `json:",omitempty"`

	ScanData string `json:",omitempty"`

	Emissions []Emission `json:",omitempty"`
	Color     color.RGBA
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
	type just GalaxyPoint
	buf, err := json.Marshal(just(gp))
	if err != nil {
		return buf, err
	}
	buf = bytes.Replace(buf, []byte(`"Pos":{},`), []byte{}, -1)
	return buf, nil
}

func (BSP) CalcDegrade(base, degrade *BSP) (res *BSP) {
	if base == nil || degrade == nil {
		return &BSP{}
	}
	res = new(BSP)
	vBase := reflect.ValueOf(base).Elem()
	vDegrade := reflect.ValueOf(degrade).Elem()
	vRes := reflect.ValueOf(res).Elem()
	t := vRes.Type()
	fc := t.NumField()

	for i := 0; i < fc; i++ {
		x := vBase.Field(i).Float() * (1.0 - vDegrade.Field(i).Float())
		vRes.Field(i).SetFloat(x)
	}
	return res
}
