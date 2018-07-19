package main

import (
	"io/ioutil"
	"encoding/json"
	"math/rand"
)

const (
	NONE = 0
	WEAK = 1
	NORM = 2
	STRONG = 3
	EARTH = 1
	NEW = 2
	LIQUID = 1
	HARD = 2
	GASOUS = 3
	WAS = 4
	RADICAL = 3
	MOVING = 1
	EXTINCT = 2
	PRESENT = 1
)

type WarpStat struct {
	StarCount        int
	HardPlanetsCount int
	MineralsCount    int
	RangeFromSolar   float64
	MineralList      []int
	GasList          []int
	EMetals          []int
	EGas             []int
}

type Options struct {
}

var Opts Options

type planet struct{
	isGas bool
	spheres [15]int
	minerals []int
	grav int //*100% of G
	temp int
	radi int
}

func main(){
	buf, err := ioutil.ReadFile("warpnomi.json")
	if err != nil {
		panic(err)
	}
	var warpP map[string]WarpStat
	err = json.Unmarshal(buf, &warpP)
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadFile("planetgen_ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &Opts)

	outData:=make (map[string][]planet)

	var maxCount = 10
	for id, stat:=range warpP{
		maxCount--
		if maxCount==0{
			break
		}
		sysData:=make([]planet,0)
		hpc:=stat.HardPlanetsCount
		for i:=0; i<stat.HardPlanetsCount; i++{
			sysData = append(sysData, genHardPlanet(stat))
		}

		var numGas int
		if hpc>0 {
			numGas = int(float64(hpc) * (1 + rand.Float64()))
		} else {
			numGas = rand.Intn(4)
		}
		for i:=0; i<numGas int {
			sysData = append(sysData, genGasPlanet(stat))
		}

		outData[id] = sysData
	}

	buf,err=json.Marshal(outData)
	if err!=nil{
		panic(err)
	}
	err=ioutil.WriteFile("all_planets.json",buf,0)
	if err!=nil{
		panic(err)
	}
}
