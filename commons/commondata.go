package commons

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
)

type CommonData struct {
	PilotData  *PilotData  `json:"pd"`
	NaviData   *NaviData   `json:"nd"`
	EngiData   *EngiData   `json:"ed"`
	ServerData *ServerData `json:"sd"`
}

type PilotData struct {
	Ship           RBData  `json:"sh"`
	SessionTime    float64 `json:"ss"`
	ThrustVector   v2.V2   `json:"tv"`
	HeatProduction float64 `json:"hp"`

	//do not reload same Msg, cz of ship.Pos extrapolate and SessionTime+=dt
	MsgID int `json:"id"`
}

type NaviData struct {
	//drop items
	BeaconCount int `json:"bc"`
	//[]corpName, i.e. ["gd","gd","pre"]
	//[]planetName, i.e. ["CV8-85","RD4-42-13"]
	Mines   []string `json:"mn"`
	Landing []string `json:"ld"`

	//cosmo
	IsScanning    bool
	ScanObjectID  string
	IsOrbiting    bool
	OrbitObjectID string
	ActiveMarker  bool  `json:"ma"`
	MarkerPos     v2.V2 `json:"mp"`

	//warp
	SonarDir   float64 `json:"sd"`
	SonarRange float64 `json:"sr"`
	SonarWide  float64 `json:"sw"`
}

type EngiData struct {
	//[0.0 - 1.0]
	//0 for fully OKEY, 1 - for totally DEGRADED
	BSPDegrade    BSPDegrade
	HeatCumulated float64
	DmgCumulated  [8]float64
}

//Rework CalcDegrade on change
type BSPDegrade struct {
	Thrust,
	Thrust_rev,
	Thrust_acc,
	Thrust_rev_acc,
	Thrust_slow,
	Thrust_rev_slow,
	Thrust_heat_prod float64

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

	Radar_range,
	Radar_angle,
	Scan_range,
	Scan_speed float64

	Sonar_range_max,
	Sonar_angle_min,
	Sonar_angle_max,
	Sonar_angle_change,
	Sonar_range_change,
	Sonar_rotate_speed float64
}

type OtherShipData struct {
	Id   string
	Name string
	Ship RBData
}

type ServerData struct {
	OtherShips []OtherShipData

	MsgID int
}

func (cd CommonData) Encode() []byte {
	buf, err := json.Marshal(cd)
	if err != nil {
		Log(LVL_ERROR, "Can't marshal CommonData", err)
		return nil
	}
	return buf
}

func (CommonData) Decode(buf []byte) (cd CommonData, err error) {
	err = json.Unmarshal(buf, &cd)
	if err != nil {
		return CommonData{}, err
	}
	return cd, nil
}

func (cd CommonData) Part(roleName string) CommonData {
	res := CommonData{}
	switch roleName {
	case ROLE_Pilot:
		res.PilotData = cd.PilotData
	case ROLE_Navi:
		res.NaviData = cd.NaviData
	case ROLE_Engi:
		res.EngiData = cd.EngiData
	case ROLE_Server:
		res.ServerData = cd.ServerData
	default:
		panic("CommonData.Part: Unknown role " + roleName)
	}
	return res
}

func (cd CommonData) FillNotNil(dest *CommonData) {
	if cd.PilotData != nil {
		if dest.PilotData == nil || dest.PilotData.MsgID != cd.PilotData.MsgID {
			dest.PilotData = cd.PilotData
		}
	}
	if cd.NaviData != nil {
		dest.NaviData = cd.NaviData
	}
	if cd.EngiData != nil {
		dest.EngiData = cd.EngiData
	}
	if cd.ServerData != nil {
		if dest.ServerData == nil || dest.ServerData.MsgID != cd.ServerData.MsgID {
			dest.ServerData = cd.ServerData
		}
	}
}

func (cd CommonData) WithoutRole(roleName string) CommonData {
	switch roleName {
	case ROLE_Pilot:
		cd.PilotData = nil
	case ROLE_Navi:
		cd.NaviData = nil
	case ROLE_Engi:
		cd.EngiData = nil
	case ROLE_Server:
		cd.ServerData = nil
	default:
		panic("CommonData.WithoutRole: Unknown role " + roleName)
	}
	return cd
}

func (CommonData) Empty() CommonData {
	return CommonData{
		PilotData:  &PilotData{},
		NaviData:   &NaviData{Mines: []string{}, Landing: []string{}},
		EngiData:   &EngiData{},
		ServerData: &ServerData{OtherShips: []OtherShipData{}},
	}
}

//deep copy
func (cd CommonData) Copy() (res CommonData) {
	res = cd
	if cd.PilotData != nil {
		val := *cd.PilotData
		res.PilotData = &val
	}
	if cd.NaviData != nil {
		val := *cd.NaviData
		mines := make([]string, len(cd.NaviData.Mines))
		landings := make([]string, len(cd.NaviData.Landing))
		copy(mines, cd.NaviData.Mines)
		copy(landings, cd.NaviData.Landing)
		val.Mines = mines
		val.Landing = landings
		res.NaviData = &val
	}
	if cd.EngiData != nil {
		val := *cd.EngiData
		res.EngiData = &val
	}
	if cd.ServerData != nil {
		val := *cd.ServerData
		otherShips := make([]OtherShipData, len(cd.ServerData.OtherShips))
		copy(otherShips, cd.ServerData.OtherShips)
		val.OtherShips = otherShips
		res.ServerData = &val
	}
	return res
}
