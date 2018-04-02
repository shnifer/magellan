package commons

import (
	"bytes"
	"encoding/json"
	"github.com/Shnifer/magellan/graph"
	"io/ioutil"
)

// map[PartName]json_PartStruct
type CMapData map[string]string

func (md CMapData) Encode() ([]byte, error) {
	buf, err := json.Marshal(md)
	if err != nil {
		Log(LVL_ERROR, "can't encode RoomCommonData")
		return nil, err
	}
	return buf, nil
}

//static method!
func (CMapData) Decode(data []byte) (CMapData, error) {
	rcd := CMapData{}
	err := json.Unmarshal(data, &rcd)
	if err != nil {
		Log(LVL_ERROR, "can't decode RoomCommonData")
		return nil, err
	}
	return rcd, nil
}

type CBSP struct {
	MaxSpeed float64
	Systems  []float64
}

func (CBSP) New() CBSP {
	return CBSP{
		Systems: make([]float64, 8),
	}
}

func (bsp CBSP) Encode() string {
	buf, err := json.Marshal(bsp)
	if err != nil {
		Log(LVL_ERROR, "can't encode CBSP")
		return ""
	}
	return string(buf)
}

func (CBSP) Decode(data []byte) CBSP {
	bsp := CBSP{}
	err := json.Unmarshal(data, &bsp)
	if err != nil {
		Log(LVL_ERROR, "CBSP.Decode can't Unmarshal")
		panic(err)
	}

	return bsp
}

type CGalaxyObj struct {
	ID string

	ParentID string
	Radius   float64
	AngSpeed float64
	AngStart float64

	ObjType string
	Size    float64

	ScienceData string
}

type CGalaxy struct {
	Objects []CGalaxyObj
}

func (CGalaxy) New() CGalaxy {
	return CGalaxy{
		Objects: make([]CGalaxyObj, 0),
	}
}

func (galaxy CGalaxy) Encode() string {
	buf, err := json.Marshal(galaxy)
	if err != nil {
		Log(LVL_ERROR, "can't encode CGalaxy")
		return ""
	}
	return string(buf)
}

func (CGalaxy) Decode(data []byte) CGalaxy {
	galaxy := CGalaxy{}
	err := json.Unmarshal(data, &galaxy)
	if err != nil {
		Log(LVL_ERROR, "CGalaxy.Decode can't Unmarshal")
		panic(err)
	}

	return galaxy
}

//creates examples of DB files
func SaveDataExamples(path string) {
	bsp := []byte(CBSP{}.New().Encode())
	bufBsp := bytes.Buffer{}
	json.Indent(&bufBsp, bsp, "", "    ")
	ioutil.WriteFile(path+"example_bsp.json", bufBsp.Bytes(), 0)

	galaxy := CGalaxy{}.New()
	galaxy.Objects = append(galaxy.Objects, CGalaxyObj{})
	bufGalaxy := bytes.Buffer{}
	json.Indent(&bufGalaxy, []byte(galaxy.Encode()), "", "    ")
	ioutil.WriteFile(path+"example_galaxy.json", bufGalaxy.Bytes(), 0)
}

type CShipPos struct {
	pos    graph.Point
	ang    float64
	vel    graph.Point
	angVel float64
}

func (shipPos CShipPos) Encode() string {
	buf, err := json.Marshal(shipPos)
	if err != nil {
		Log(LVL_ERROR, "can't encode CShipPos")
		return ""
	}
	return string(buf)
}

func (CShipPos) Decode(data []byte) CShipPos {
	shipPos := CShipPos{}
	err := json.Unmarshal(data, &shipPos)
	if err != nil {
		Log(LVL_ERROR, "CShipPos.Decode can't Unmarshal")
		panic(err)
	}

	return shipPos
}
