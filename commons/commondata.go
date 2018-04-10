package commons

import (
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
)

//just Rigid Body data
type RBData struct {
	Pos    graph.Point
	Ang    float64
	Vel    graph.Point
	AngVel float64
}

type CommonData struct {
	PilotData *PilotData
	NaviData  *NaviData
	CargoData *CargoData
	EngiData  *EngiData
}

type PilotData struct {
	Ship RBData
}

type NaviData struct {
	ActiveMarker bool
	MarkerPos    graph.Point
}
type CargoData struct {
	TurboBoost bool
}
type EngiData struct {
	//[0.0 - 1.0]
	Systems [8]float64
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
	case ROLE_Cargo:
		res.CargoData = cd.CargoData
	default:
		panic("CommonData.Part: Unknown role " + roleName)
	}
	return res
}

func (cd CommonData) FillNotNil(dest *CommonData) {
	if cd.PilotData != nil {
		dest.PilotData = cd.PilotData
	}
	if cd.NaviData != nil {
		dest.NaviData = cd.NaviData
	}
	if cd.EngiData != nil {
		dest.EngiData = cd.EngiData
	}
	if cd.CargoData != nil {
		dest.CargoData = cd.CargoData
	}
}

func (cd CommonData) ClearRole(roleName string) CommonData {
	switch roleName {
	case ROLE_Pilot:
		cd.PilotData = nil
	case ROLE_Navi:
		cd.NaviData = nil
	case ROLE_Engi:
		cd.EngiData = nil
	case ROLE_Cargo:
		cd.CargoData = nil
	default:
		panic("CommonData.ClearRole: Unknown role " + roleName)
	}
	return cd
}

func (CommonData) Empty() CommonData {
	return CommonData{
		PilotData: &PilotData{},
		NaviData:  &NaviData{},
		CargoData: &CargoData{},
		EngiData:  &EngiData{},
	}
}