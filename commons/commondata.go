package commons

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
)

type CommonData struct {
	PilotData *PilotData
	NaviData  *NaviData
	EngiData  *EngiData
}

type PilotData struct {
	Ship           RBData
	SessionTime    float64
	ThrustVector   v2.V2
	HeatProduction float64

	//do not reload same Msg, cz of extrapolate
	MsgID int
}

type NaviData struct {
	//cosmo
	ActiveMarker bool
	MarkerPos    v2.V2

	//warp
	SonarDir   float64
	SonarRange float64
	SonarWide  float64
}
type CargoData struct {
	TurboBoost bool
}

type EngiData struct {
	//[0.0 - 1.0]
	//0 for fully OKEY, 1 - for totally DEGRADED
	BSPDegrade    *BSP
	HeatCumulated float64
	DmgCumulated  [8]float64
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
}

func (cd CommonData) WithoutRole(roleName string) CommonData {
	switch roleName {
	case ROLE_Pilot:
		cd.PilotData = nil
	case ROLE_Navi:
		cd.NaviData = nil
	case ROLE_Engi:
		cd.EngiData = nil
	default:
		panic("CommonData.WithoutRole: Unknown role " + roleName)
	}
	return cd
}

func (CommonData) Empty() CommonData {
	return CommonData{
		PilotData: &PilotData{},
		NaviData:  &NaviData{},
		EngiData:  &EngiData{BSPDegrade: &BSP{}},
	}
}

func (cd CommonData) Copy() (res CommonData) {
	res = cd
	if cd.PilotData != nil {
		val := *cd.PilotData
		res.PilotData = &val
	}
	if cd.NaviData != nil {
		val := *cd.NaviData
		res.NaviData = &val
	}
	if cd.EngiData != nil {
		val := *cd.EngiData
		res.EngiData = &val
	}
	return res
}
