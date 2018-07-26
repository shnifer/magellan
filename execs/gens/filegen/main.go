package main

import (
	"encoding/json"
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	"io/ioutil"
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

	fmt.Println("files loaded")
	for sysName, stat := range warpP {
		fmt.Println("system " + sysName)
		pref := sysName + "-"
		planets := allPlanet[sysName]
		points := make(map[string]*GalaxyPoint)

		fmt.Println("got")
		createStars(stat, points, pref)
		fmt.Println("stars created")
		createPlanets(stat, points, pref, planets)
		fmt.Println("planets created")

		galaxy := Galaxy{
			Points:        points,
			SpawnDistance: 5000,
		}
		dat, err := json.Marshal(galaxy)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile("galaxy_"+sysName+".json", dat, 0)
	}
}
