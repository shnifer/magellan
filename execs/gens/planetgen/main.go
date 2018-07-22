package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	NONE                   = 0
	WEAK                   = 1
	NORM                   = 2
	STRONG                 = 3
	EARTH                  = 1
	NEW                    = 2
	EARTHANDNEW            = 3
	LIQUID                 = 1
	HARD                   = 2
	GASOUS                 = 3
	HARDANDGASOUS          = 5
	LIQUIDANDGASOUS        = 6
	HARDANDLIQUID          = 7
	HARDANDLIQUIDANDGASOUS = 8
	WAS                    = 4
	RADICAL                = 3
	MOVING                 = 1
	EXTINCT                = 2
	PRESENT                = 1
)

const (
	MAGNET = iota
	RADIATIONBELT
	OXYGEN
	GASES
	ATMOMETALS
	OZONE
	ION
	WATER
	MIXTURES
	PEDOMETALS
	COREMADE
	COREVEL
	VULCAN
	LITOMETAL
	BIO
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

type Options struct {
}

var Opts Options

type Planet struct {
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

	dat, err := ioutil.ReadFile("planetgen_ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &Opts)

	f, err := os.Create("all_planets.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.UseCRLF = true
	w.Write([]string{"система", "звёзд", "твёрдых планет", "кол-во минералов системы", "минералы системы", "газы системы", "богатых планет", "шахт",
		"номер планеты", "газовая", "металлы планеты", "температура"})

	outData := make(map[string][]Planet)

	for id, stat := range warpP {
		stat.newElements = len(stat.GasList)+len(stat.MineralList) > 0
		stat.g7 = has(stat.GasList, 7)
		stat.eg = len(stat.EGas) > 0
		stat.em = len(stat.EMetals) > 0
		sysData := make([]Planet, 0)
		hpc := stat.HardPlanetsCount
		for i := 0; i < stat.HardPlanetsCount; i++ {
			sysData = append(sysData, genHardPlanet(stat, i))
		}

		var numGas int
		if hpc > 0 {
			numGas = int(float64(hpc) * (1 + rand.Float64()))
		} else {
			numGas = rand.Intn(4)
		}
		for i := 0; i < numGas; i++ {
			sysData = append(sysData, genGasPlanet(stat))
		}

		//mineralsSpread
		rich := make([]int, 0)
		for i := 0; i < stat.RichPlanets; i++ {
			var ok bool
			var n int
			var z int
			for !ok {
				z++
				if z > 100 {
					log.Println(id, "can't ", stat.HardPlanetsCount, rich)
					break
				}
				n = rand.Intn(stat.HardPlanetsCount)
				ok = !has(rich, n)
			}
			rich = append(rich, n)
		}
		if len(stat.MineralList) == 1 {
			for _, plN := range rich {
				sysData[plN].Minerals = stat.MineralList
			}
		} else if len(stat.MineralList) == 2 {
			if stat.RichPlanets == 1 {
				sysData[rich[0]].Minerals = stat.MineralList
			} else if stat.MinesCount == 3 {
				sysData[rich[0]].Minerals = stat.MineralList
				sysData[rich[1]].Minerals = []int{stat.MineralList[rand.Intn(2)]}
			} else if stat.MinesCount == 4 {
				sysData[rich[0]].Minerals = stat.MineralList
				sysData[rich[1]].Minerals = stat.MineralList
			}
		}

		outData[id] = sysData
	}

	//csv out
	for id, stat := range warpP {
		for ind, planet := range outData[id] {
			str := []string{id, fs(stat.StarCount), fs(stat.HardPlanetsCount), fs(len(stat.MineralList)),
				fs(stat.MineralList), fs(stat.GasList), fs(stat.RichPlanets), fs(stat.MinesCount)}
			str = append(str, fs(ind), fs(planet.IsGas), fs(planet.Minerals), fs(planet.Temperature))
			w.Write(str)
		}
	}
	w.Flush()

	buf, err = json.Marshal(outData)
	if err != nil {
		panic(err)
	}
	identbuf := bytes.Buffer{}
	json.Indent(&identbuf, buf, "", "    ")
	err = ioutil.WriteFile("all_planets.json", identbuf.Bytes(), 0)
	if err != nil {
		panic(err)
	}
}

func genGasPlanet(stat WarpStat) Planet {
	res := Planet{
		IsGas:    true,
		Minerals: []int{},
	}
	sph := res.Spheres
	//3.1
	if stat.newElements {
		if !stat.g7 {
			sph[COREMADE] = r([]int{1, 2, 2}, NONE, EARTHANDNEW, NEW)
		}
	} else {
		sph[COREMADE] = r([]int{1, 2}, NONE, EARTH)
	}
	//3.2
	sph[COREVEL] = MOVING
	//3.3
	sph[MAGNET] = rr(NORM, STRONG)
	sph[RADIATIONBELT] = STRONG
	//4.1
	if !stat.g7 {
		sph[GASES] = EARTH
	} else if stat.eg {
		sph[GASES] = EARTHANDNEW
	} else {
		sph[GASES] = NEW
	}
	//4.2
	sph[ION] = rr(NORM, STRONG)
	//5.1
	sph[WATER] = rr(NONE, GASOUS, HARDANDGASOUS)
	//5.2
	if sph[WATER] != NONE {
		if !stat.newElements {
			sph[MIXTURES] = EARTH
		} else if stat.StarCount == 1 {
			sph[MIXTURES] = rr(EARTHANDNEW, NEW)
		} else {
			sph[MIXTURES] = rr(EARTHANDNEW, NEW, RADICAL)
		}
	}
	//6
	if has(stat.GasList, 1) || has(stat.GasList, 2) {
		sph[BIO] = rr(WAS, PRESENT)
	}

	res.Spheres = sph
	return res
}

func genHardPlanet(stat WarpStat, num int) Planet {
	res := Planet{
		IsGas:    false,
		Minerals: []int{},
	}
	sph := res.Spheres
	isClosest := num == 0
	var isGreen bool
	switch stat.HardPlanetsCount {
	case 2, 3:
		isGreen = num == 1
	case 4, 5:
		isGreen = num == 2 || num == 3
	}
	var x []int
	switch stat.HardPlanetsCount {
	case 1:
		x = []int{200, 150, 100}
	case 2:
		x = []int{100, 50, 15, 15}
	case 3:
		x = []int{100, 50, 15, 15, 0}
	case 4:
		x = []int{100, 50, 15, 15, 0, -15}
	case 5:
		x = []int{100, 50, 15, 15, 0, -15, -50}
	}
	res.Temperature = rr(x[num : num+3]...)

	//1
	if !stat.newElements {
		sph[COREMADE] = EARTH
	} else if stat.em {
		sph[COREMADE] = EARTHANDNEW
	} else {
		sph[COREMADE] = NEW
	}
	//2
	sph[COREVEL] = rr(MOVING, EXTINCT)
	//3
	if sph[COREVEL] == MOVING {
		sph[MAGNET] = rr(WEAK, NORM, STRONG)
	}
	//4
	sph[RADIATIONBELT] = sph[MAGNET]
	res.Radiation = STRONG - sph[RADIATIONBELT]
	//5
	sph[VULCAN] = rr(WAS, WEAK, NORM, STRONG)
	if has(stat.GasList, 4) || sph[VULCAN] == STRONG {
		res.Temperature = rr(200, 400, 500)
	}
	//6
	if !stat.newElements {
		sph[LITOMETAL] = EARTH
	} else if stat.em {
		sph[LITOMETAL] = EARTHANDNEW
	} else {
		sph[LITOMETAL] = NEW
	}
	//7
	if sph[VULCAN] == STRONG {
		sph[WATER] = rr(NONE, WAS)
	} else if isGreen {
		if sph[MAGNET] == NONE {
			sph[WATER] = rr(NONE, WAS, HARD)
		} else {
			sph[WATER] = rr(NONE, WAS, HARDANDGASOUS, HARDANDLIQUID, LIQUIDANDGASOUS, HARDANDLIQUIDANDGASOUS)
		}
	} else {
		if sph[MAGNET] == NONE {
			sph[WATER] = rr(NONE, WAS, HARD)
		} else {
			sph[WATER] = rr(NONE, WAS, HARD, LIQUID, GASOUS, HARDANDLIQUID, HARDANDGASOUS, LIQUIDANDGASOUS)
		}
	}
	//8
	if sph[WATER] != NONE && sph[WATER] != WAS {
		if !stat.newElements {
			sph[MIXTURES] = EARTH
		} else {
			sph[MIXTURES] = rr(EARTHANDNEW, NEW, RADICAL)
		}
	}
	//9
	if !stat.newElements {
		sph[PEDOMETALS] = EARTH
	} else if stat.em {
		sph[PEDOMETALS] = EARTHANDNEW
	} else {
		sph[PEDOMETALS] = NEW
	}
	//10
	if !isClosest {
		if !stat.newElements {
			sph[GASES] = EARTH
		} else if stat.eg {
			sph[GASES] = EARTHANDNEW
		} else {
			sph[GASES] = NEW
		}
		//11
		sph[ION] = rr(WEAK, NORM, STRONG)
		//12
		if !stat.newElements {
			sph[ATMOMETALS] = r([]int{2, 1}, NONE, EARTH)
		} else if stat.em {
			sph[ATMOMETALS] = r([]int{3, 1, 1}, NONE, EARTHANDNEW, NEW)
		} else {
			sph[ATMOMETALS] = r([]int{2, 1}, NONE, NEW)
		}
		//13
		if sph[MAGNET] != NONE && sph[WATER] != NONE && sph[WATER] != WAS {
			sph[OXYGEN] = r([]int{21, 1, 1, 1}, NONE, WEAK, NORM, STRONG)
		}
		//14
		if sph[OXYGEN] == NONE && sph[ION] == WEAK {
			sph[OZONE] = NONE
		} else if sph[OXYGEN] == STRONG {
			sph[OZONE] = STRONG
		} else {
			sph[OZONE] = rr(WEAK, NORM)
		}
	}
	//BIO
	if sph[OXYGEN]+sph[GASES]+sph[ATMOMETALS]+sph[OZONE]+sph[ION] == NONE ||
		(sph[MAGNET]+sph[RADIATIONBELT] == NONE && sph[ION] == NONE) ||
		sph[COREMADE] == NONE || sph[VULCAN] == STRONG || sph[WATER] == WATER ||
		isClosest || res.Temperature >= 200 || res.Temperature <= -100 {
		sph[BIO] = NONE
	} else if !res.IsGas && (sph[COREVEL] == EXTINCT || sph[VULCAN] == NORM || res.Temperature == -50) {
		sph[BIO] = WAS
	} else if sph[OXYGEN] > NONE ||
		(sph[GASES] > NONE && (has(stat.GasList, 1) || has(stat.GasList, 2))) ||
		sph[WATER] > NONE || (has(stat.MineralList, 1) || has(stat.MineralList, 2)) {
		sph[BIO] = PRESENT
	}

	res.Spheres = sph
	return res
}

func r(base []int, vals ...int) int {
	if len(base) != len(vals) {
		panic("base len != vals len")
	}
	s := sum(base)
	if s == 0 {
		panic("zero base")
	}
	n := rand.Intn(s)

	s = 0
	for i := 0; i < len(base); i++ {
		s += base[i]
		if s > n {
			return vals[i]
		}
	}
	return 0
}

func rr(vals ...int) int {
	n := rand.Intn(len(vals))
	return vals[n]
}

func has(a []int, n int) bool {
	for _, v := range a {
		if v == n {
			return true
		}
	}
	return false
}
func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}

func fs(v interface{}) string {
	str := fmt.Sprint(v)
	return strings.Replace(str, ".", ",", -1)
}
