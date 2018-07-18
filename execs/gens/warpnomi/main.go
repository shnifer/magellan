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
}

type WarpStat struct {
	HardPlanetsCount int
	MineralsCount    int
	RangeFromSolar   float64
	MineralList      []int
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
		hp := extract(hpBase)
		ms := 0
		if hp > 0 {
			ms = extract(msBase)
		}
		minerals := make([]int, ms)
		for i := 0; i < ms; i++ {
			var min int
			if solarRange < Opts.NoobRange {
				min = extractDiffer(noobBase, minerals[:i])
				if min > 0 {
					mnrlsBase[min]--
				}
			}
			if min == 0 {
				min = extractDiffer(mnrlsBase, minerals[:i])
			}
			minerals[i] = min
		}
		sort.Ints(minerals)

		minstrs := make([]string, len(Opts.Minerals))
		for _, v := range minerals {
			minstrs[v-1] = "1"
		}

		w.Write(append([]string{id, fs(hp), fs(ms), fs(solarRange)}, minstrs...))

		stats[id] = WarpStat{
			RangeFromSolar:   solarRange,
			HardPlanetsCount: hp,
			MineralsCount:    ms,
			MineralList:      minerals,
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
	return fmt.Sprint(v)
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
