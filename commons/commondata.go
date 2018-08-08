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
	Ship        RBData  `json:"sh,omitempty"`
	SessionTime float64 `json:"ss,omitempty"`
	FlightTime  float64 `json:"ft,omitempty"`
	//for cosmo
	ThrustVector v2.V2 `json:"tv,omitempty"`
	//for warp
	Distortion float64 `json:"wd,omitempty"`
	Dir        float64 `json:"dr,omitempty"`
	//warp position for return from zero system
	WarpPos v2.V2 `json:"wp,omitempty"`

	//to Engi

	HeatProduction float64 `json:"hp,omitempty"`

	//do not reload same Msg, cz of ship.Pos extrapolate and SessionTime+=dt
	MsgID int `json:"id"`
}

type NaviData struct {
	//drop items
	BeaconCount int `json:"bc,omitempty"`
	//[]corpName, i.e. ["gd","gd","pre"]
	//[]planetName, i.e. ["CV8-85","RD4-42-13"]
	Mines   []string `json:"mn,omitempty"`
	Landing []string `json:"ld,omitempty"`

	//cosmo
	IsScanning    bool   `json:"is,omitempty"`
	IsDrop        bool   `json:"st,omitempty"`
	ScanObjectID  string `json:"so,omitempty"`
	IsOrbiting    bool   `json:"io,omitempty"`
	OrbitObjectID string `json:"oo,omitempty"`
	ActiveMarker  bool   `json:"ma,omitempty"`
	MarkerPos     v2.V2  `json:"mp,omitempty"`

	//warp
	SonarDir   float64 `json:"sd,omitempty"`
	SonarRange float64 `json:"sr,omitempty"`
	SonarWide  float64 `json:"sw,omitempty"`
}

type EngiCounters struct {
	Fuel       float64 `json:"f,omitempty"`
	HoleSize   float64 `json:"h,omitempty"`
	Pressure   float64 `json:"p,omitempty"`
	Air        float64 `json:"a,omitempty"`
	Calories   float64 `json:"t,omitempty"`
	CO2        float64 `json:"co2,omitempty"`
	FlightTime float64 `json:"ft,omitempty"`
}

type EngiData struct {
	//[0.0 - 1.0]
	//0 for fully OKEY, 1 - for totally DEGRADED
	BSPDegrade BSPDegrade         `json:"deg,omitempty"`
	AZ         [8]float64         `json:"az,omitempty"`
	InV        [8]uint16          `json:"inv,omitempty"`
	Emissions  map[string]float64 `json:"emm,omitempty"`

	//Counters
	Counters EngiCounters `json:"c,omitempty"`
}

//Rework CalcDegrade on change
type BSPDegrade BSPParams

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
		EngiData:   &EngiData{Emissions: make(map[string]float64), BSPDegrade: emptyDegrade()},
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
