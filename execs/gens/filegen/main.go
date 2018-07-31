package main

import (
	"encoding/json"
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"io/ioutil"
	"sort"
)

type WarpStat struct {
	StarCount        int
	HardPlanetsCount int
	MineralsTypes    int
	RichPlanets      int
	MinesCount       int
	RangeFromSolar   float64
	MineralList      []int
	GasList          []int
	newElements      bool
	g7               bool
	eg               bool
	em               bool
	EMetals          []int
	EGas             []int
}

type Planet struct {
	ID          string
	IsGas       bool
	Spheres     [15]int
	Minerals    []int
	Grav        int //*100% of G
	Temperature int
	Radiation   int
}

func main() {
	buf, err := ioutil.ReadFile("warpnomi.json")
	if err != nil {
		panic(err)
	}
	var warpP map[string]WarpStat
	err = json.Unmarshal(buf, &warpP)
	if err != nil {
		panic(err)
	}

	buf, err = ioutil.ReadFile("filegen_ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(buf, &Opts)

	buf, err = ioutil.ReadFile("all_planets.json")
	if err != nil {
		panic(err)
	}
	var allPlanet map[string][]Planet
	err = json.Unmarshal(buf, &allPlanet)
	if err != nil {
		panic(err)
	}

	warpSignatures := make(map[string][]Signature)
	warpMinerals := make(map[string][]int)

	fmt.Println("files loaded")
	for sysName, stat := range warpP {
		pref := sysName + "-"
		planets := allPlanet[sysName]
		points := make(map[string]*GalaxyPoint)

		createStars(stat, points, pref)
		createPlanets(stat, points, pref, planets)

		sigs := make(map[string]Signature)
		minerals := make(map[int]struct{})
		var sd float64
		sd = 1500
		for _, p := range points {
			if sd < p.Orbit {
				sd = p.Orbit
			}
			for _, sig := range p.Signatures {
				if v, exist := sigs[sig.TypeName]; exist {
					v.Dev = okrV2(v2.RandomInCircle(1))
					sigs[sig.TypeName] = v
				} else {
					sigs[sig.TypeName] = sig
				}
			}
			for _, min := range p.Minerals {
				minerals[min] = struct{}{}
			}
		}

		for _, v := range sigs {
			warpSignatures[sysName] = append(warpSignatures[sysName], v)
		}
		for i := range minerals {
			warpMinerals[sysName] = append(warpMinerals[sysName], i)
		}
		sort.IntSlice(warpMinerals[sysName]).Sort()

		galaxy := Galaxy{
			Points:        points,
			SpawnDistance: sd * 1.1,
		}
		dat, err := json.Marshal(galaxy)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("galaxy_"+sysName+".json", dat, 0)
	}

	var warpGal Galaxy
	wgDat, err := ioutil.ReadFile("galaxy_warp.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(wgDat, &warpGal)
	if err != nil {
		panic(err)
	}

	for id, sigs := range warpSignatures {
		v, ok := warpGal.Points[id]
		if !ok {
			continue
		}
		v.Signatures = sigs
		warpGal.Points[id] = v
	}

	for id, mins := range warpMinerals {
		v, ok := warpGal.Points[id]
		if !ok {
			continue
		}
		v.Minerals = mins
		warpGal.Points[id] = v
	}

	wgDat, err = json.Marshal(warpGal)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("galaxy_warp.json", wgDat, 0)
}
