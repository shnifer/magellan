package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Shnifer/magellan/v2"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Options struct {
	SysWithHardPlanets []int
	SysWithMinerals    []int
	Minerals           []int

	NoobRange    float64
	NoobMinerals []int

	DoubleStarPercent1Min  int
	TripleStartPercent2Min int

	Gas7Percent int

	EarthGasPercent   int
	EarthMetalPercent int
}

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

var Opts Options

func main() {
	f, err := os.Create("warpnomi.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.UseCRLF = true
	w.Write([]string{"id", "звёзд", "твёрдых планет", "нов. металлов в системе", "нов. газов в системе", "до солнца",
		"М1", "М2", "М3", "М4", "М5", "М6", "М7", "G1", "G2", "G3", "G4", "G5", "G6", "G7", "Железо", "Никель", "Магний", "Титан", "Аллюминий",
		"Азот", "Углекислота", "Криптон", "Гелий", "Сера"})

	buf, err := ioutil.ReadFile("starpos.json")
	if err != nil {
		panic(err)
	}
	var pts map[string]v2.V2
	err = json.Unmarshal(buf, &pts)
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadFile("warpnomi_ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &Opts)

	solarP := pts["solar"]
	stats := make(map[string]WarpStat)
	total := len(pts) - 1
	hpBase := newBase(total, Opts.SysWithHardPlanets)
	s := 0
	for _, v := range Opts.SysWithHardPlanets {
		s += v
	}
	msBase := newBase(s, Opts.SysWithMinerals)
	mnrlsBase := append([]int{0}, Opts.Minerals...)
	noobBase := append([]int{0}, Opts.NoobMinerals...)

	log.Println("mnrlsBase ", mnrlsBase)
	log.Println("noobBase ", noobBase)

	for id, p := range pts {
		if id == "solar" {
			continue
		}
		solarRange := p.Sub(solarP).Len()
		var useNoob bool
		if solarRange < Opts.NoobRange && sum(noobBase) > 0 && msBase[1] > 0 {
			useNoob = true
		}
		var hp int
		var ok bool
		for !ok {
			hp = get(hpBase)
			if !useNoob || hp > 0 {
				ok = true
			}
		}
		hpBase[hp]--
		ms := 0
		if hp > 0 {
			if useNoob {
				var ok bool
				for !ok {
					ms = get(msBase)
					if ms > 0 {
						ok = true
					}
				}
				msBase[ms]--
			} else {
				ms = extract(msBase)
			}
		}
		minerals := make([]int, ms)
		for i := 0; i < ms; i++ {
			if i == 0 && solarRange < Opts.NoobRange {
				min := extract(noobBase)
				if min > 0 {
					mnrlsBase[min]--
					minerals[i] = min
					continue
				}
			}
			if i == 0 {
				n := extract(mnrlsBase)
				if n > 0 {
					minerals[i] = n
				} else {
					log.Println("auchtung")
				}
			} else {
				var n, z int
				var ok bool
				for !ok {
					z++
					n = get(mnrlsBase)
					if n == minerals[0] {
						continue
					}
					if minerals[0] <= 3 && n <= 3 || minerals[0] > 3 && n > 3 {
						ok = true
					}
					if z > 100 {
						log.Println("Can't gen new mineral for [0]=", minerals[0], " have ", mnrlsBase)
						n = 0
						ok = true
					}
				}
				if n > 0 {
					mnrlsBase[n]--
					minerals[i] = n
				}
			}
		}
		sort.Ints(minerals)

		minstrs := make([]string, len(Opts.Minerals))
		for _, v := range minerals {
			minstrs[v-1] = "1"
		}

		var starCount int
		var g7 bool
		switch ms {
		case 0:
			starCount = 1
			if hp == 0 {
				if rand.Intn(100) < Opts.Gas7Percent {
					g7 = true
					if rand.Intn(100) < Opts.DoubleStarPercent1Min {
						starCount = 2
					} else if rand.Intn(100) < Opts.TripleStartPercent2Min {
						starCount = 3
					}
				}
			}
		case 1:
			if rand.Intn(100) < Opts.DoubleStarPercent1Min {
				starCount = 2
			} else {
				starCount = 1
			}
		case 2:
			if rand.Intn(100) < Opts.TripleStartPercent2Min {
				starCount = 3
			} else {
				starCount = 2
			}
		}

		var newElements = g7 || ms > 0
		gas := make([]int, 0)
		if newElements {
			gas = genGas(starCount, g7, minerals)
		}
		gasstrs := make([]string, 7)
		for _, v := range gas {
			gasstrs[v-1] = "1"
		}

		earthGas := make([]int, 0)
		earthMetal := make([]int, 0)
		if starCount < 3 {
			if !has(minerals, 1) && !has(minerals, 4) &&
				!has(gas, 2) && !has(gas, 3) && !has(gas, 4) {
				earthGas = genEGas(minerals, gas)
			}
			if !g7 && !has(minerals, 1) && !has(minerals, 2) && !has(minerals, 3) {
				earthMetal = genEMetal(minerals, gas)
			}
		}
		Emetstrs := make([]string, 5)
		for _, v := range earthMetal {
			Emetstrs[v-1] = "1"
		}
		Egasstrs := make([]string, 5)
		for _, v := range earthGas {
			Egasstrs[v-1] = "1"
		}

		strs := []string{id, fs(starCount), fs(hp), fs(ms), fs(len(gas)), fs(solarRange)}
		strs = append(strs, minstrs...)
		strs = append(strs, gasstrs...)
		strs = append(strs, Emetstrs...)
		strs = append(strs, Egasstrs...)

		w.Write(strs)

		stats[id] = WarpStat{
			StarCount:        starCount,
			RangeFromSolar:   solarRange,
			HardPlanetsCount: hp,
			MineralsCount:    ms,
			MineralList:      minerals,
			GasList:          gas,
			EMetals:          earthMetal,
			EGas:             earthGas,
		}
	}
	w.Flush()

	dat, err = json.Marshal(stats)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("warpnomi.json", dat, 0)
}

func fs(v interface{}) string {
	str := fmt.Sprint(v)
	return strings.Replace(str, ".", ",", -1)
}

func extract(base []int) int {
	n := get(base)
	if n < len(base) {
		if base[n] > 0 {
			base[n]--
		}
	}
	return n
}

func extractDiffer(base []int, has []int) int {
	if len(has) >= len(base) {
		return 0
	}
	var n int
	var ok bool
	for !ok {
		n = get(base)
		ok = true
		for _, v := range has {
			if n == v {
				ok = false
				break
			}
		}
	}
	if n < len(base) {
		if base[n] > 0 {
			base[n]--
		}
	}
	return n
}

func get(base []int) int {
	if len(base) == 0 {
		return 0
	}
	s := 0
	for _, v := range base {
		s += v
	}
	if s == 0 {
		return 0
	}
	n := rand.Intn(s)
	s = 0
	for i, v := range base {
		if s+v >= n+1 {
			return i
		}
		s += v
	}
	return len(base) - 1
}

func newBase(total int, opts []int) []int {
	s := 0
	for _, v := range opts {
		s += v
	}
	base := append([]int{total - s}, opts...)
	return base
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

func genGas(stars int, g7 bool, minerals []int) []int {
	hasm := func(n int) bool {
		return has(minerals, n)
	}
	totalG := 1
	if stars == 3 {
		totalG = 2
	} else if stars == 2 {
		if rand.Intn(2) == 0 {
			totalG = 2
		}
	}
	if g7 {
		totalG--
	}
	if totalG == 0 {
		return []int{7}
	}
	res := make([]int, totalG)
	var group int
	if g7 {
		group = 1
	} else if hasm(4) {
		group = 0
	} else if rand.Intn(2) == 0 {
		group = 1
	}
	for i := 0; i < totalG; i++ {
		var n int
		var ok bool
		var z int
		for !ok {
			z++
			if z > 100 {
				log.Println(i, group, minerals, res[0])
			}
			if group == 0 {
				switch rand.Intn(5) {
				case 0:
					n = 2
				case 1, 2:
					n = 3
				case 3, 4:
					n = 4
				}
			} else {
				switch rand.Intn(5) {
				case 0:
					n = 1
				case 1, 2:
					n = 5
				case 3, 4:
					n = 6
				}
			}
			if i == 0 {
				ok = true
			} else if res[0] != n {
				ok = true
				if n == 1 && (hasm(1) || hasm(4)) {
					ok = false
				}
				if n == 2 && (hasm(5) || hasm(6)) {
					ok = false
				}
				if n == 5 && hasm(4) {
					ok = false
				}
				if n == 6 && (hasm(2) || hasm(4)) {
					ok = false
				}
			}
			if ok {
				res[i] = n
			}
		}
	}
	if g7 {
		res = append(res, 7)
	}
	sort.Ints(res)
	return res
}

func genEGas(min, gas []int) []int {
	res := make([]int, 0)
	var ok bool
	var z int
	for !ok {
		z++
		if z > 100 {
			log.Println("egas", min, gas)
		}
		res = res[:0]
		for i := 1; i <= 5; i++ {
			if i == 5 && has(gas, 6) {
				continue
			}
			if rand.Intn(100) < Opts.EarthGasPercent {
				res = append(res, i)
			}
		}
		if len(res) >= 2 {
			ok = true
		}
	}
	return res
}

func genEMetal(min, gas []int) []int {
	var ok bool
	res := make([]int, 0)
	var z int
	for !ok {
		z++
		if z > 100 {
			log.Println("emet", min, gas)
		}
		res = res[:0]
		for i := 1; i <= 5; i++ {
			if (i == 1 || i == 2) && has(min, 7) {
				continue
			}
			if (i == 3 || i == 4) && has(min, 5) {
				continue
			}
			if i == 5 && has(min, 6) {
				continue
			}
			if rand.Intn(100) < Opts.EarthMetalPercent {
				res = append(res, i)
			}
		}
		if len(res) >= 2 {
			ok = true
		}
		if has(min, 5) && has(min, 7) {
			if len(res) >= 1 {
				ok = true
			}
		}
	}
	return res
}
