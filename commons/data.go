package commons

import (
	"encoding/json"
	"io/ioutil"
	"bytes"
)

// map[PartName]json_PartStruct
type MapData map[string]string

func (md MapData) Encode() ([]byte, error){
	buf, err:= json.Marshal(md)
	if err!=nil{
		Log(LVL_ERROR, "can't encode RoomCommonData")
		return nil, err
	}
	return buf, nil
}

//static method!
func (MapData) Decode(data []byte) (MapData, error){
	rcd:=MapData{}
	err:= json.Unmarshal(data,&rcd)
	if err!=nil{
		Log(LVL_ERROR, "can't decode RoomCommonData")
		return nil, err
	}
	return rcd, nil
}

type BSP struct{
	MaxSpeed float64
	Systems []float64
}

func (BSP) New() BSP{
	return BSP{
		Systems:make([]float64,8),
	}
}

func (bsp BSP) Encode() string{
	buf, err:= json.Marshal(bsp)
	if err!=nil{
		Log(LVL_ERROR, "can't encode RoomCommonData")
		return ""
	}
	return string(buf)
}

func (BSP) Decode(data []byte) BSP{
	bsp :=BSP{}
	err := json.Unmarshal(data, &bsp)
	if err!=nil{
		Log(LVL_ERROR, "BSP.Decode can't Unmarshal")
		panic(err)
	}

	return bsp
}

type GalaxyObj struct{
	ID string

	ParentID string
	Radius float64
	AngSpeed float64
	AngStart float64

	ObjType string
	Size float64

	ScienceData string
}

type Galaxy struct{
	Objects []GalaxyObj
}

func (Galaxy) New() Galaxy{
	return Galaxy{
		Objects:make([]GalaxyObj,0),
	}
}

func (galaxy Galaxy) Encode() string{
	buf, err:= json.Marshal(galaxy)
	if err!=nil{
		Log(LVL_ERROR, "can't encode RoomCommonData")
		return ""
	}
	return string(buf)
}


func (Galaxy) Decode(data []byte) Galaxy{
	galaxy:=Galaxy{}
	err := json.Unmarshal(data, &galaxy)
	if err!=nil{
		Log(LVL_ERROR, "Galaxy.Decode can't Unmarshal")
		panic(err)
	}

	return galaxy
}

//creates examples of DB files
func SaveDataExamples(path string) {
	bsp:= []byte(BSP{}.New().Encode())
	bufBsp := bytes.Buffer{}
	json.Indent(&bufBsp, bsp, "","    ")
	ioutil.WriteFile(path+"example_bsp.json", bufBsp.Bytes(), 0)

	galaxy := Galaxy{}.New()
	galaxy.Objects = append(galaxy.Objects, GalaxyObj{})
	bufGalaxy := bytes.Buffer{}
	json.Indent(&bufGalaxy, []byte(galaxy.Encode()), "","    ")
	ioutil.WriteFile(path+"example_galaxy.json", bufGalaxy.Bytes(), 0)
}