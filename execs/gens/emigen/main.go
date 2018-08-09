package main

import (
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	files := make([]string, 0)
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		name := f.Name()
		if strings.HasPrefix(name, "galaxy_") &&
			strings.HasSuffix(name, ".json") &&
			name != "galaxy_solar.json" && !f.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	wg := new(sync.WaitGroup)
	for _, name := range files {
		log.Println("run ", name)
		go procFile(name, wg)
		wg.Add(1)
	}
	wg.Wait()
}

func procFile(fn string, wg *sync.WaitGroup) {
	defer wg.Done()

	var G *commons.Galaxy
	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dat, &G)
	if err != nil {
		panic(err)
	}

	for i, p := range G.Points {
		p.Emissions = genEmissions(p)
		p.Signatures = procSignatures(p)
		G.Points[i] = p
	}

	dat, err = json.Marshal(G)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fn, dat, 0)
	if err != nil {
		panic(err)
	}
	log.Println("done " + fn)
}

func genEmissions(p *commons.GalaxyPoint) []commons.Emission {
	res := make([]commons.Emission, 0)

	t := p.Type
	var chance EmiChance
	var distrib EmiDistib
	switch t {
	case commons.GPT_STAR:
		chance = Opts.Chance.Star
		distrib = Opts.Star
	case commons.GPT_HARDPLANET:
		chance = Opts.Chance.Hard
		distrib = Opts.Hard
	case commons.GPT_GASPLANET:
		chance = Opts.Chance.Gas
		distrib = Opts.Gas
	case commons.GPT_ASTEROID:
		chance = Opts.Chance.Asteroid
		distrib = Opts.Asteroid
	case commons.GPT_WARP:
		chance = Opts.Chance.Warp
		distrib = Opts.Warp
	default:
		log.Println("unknown type ", t)
		return res
	}

	initEmi(&res, chance, t)

	doFirst := rand.Intn(100) < chance.First
	if !doFirst {
		return res
	}
	newEmi(&res, chance, distrib)

	for rand.Intn(100) < chance.More {
		newEmi(&res, chance, distrib)
	}

	return res
}

func initEmi(E *[]commons.Emission, c EmiChance, t string) {
	switch t {
	case commons.GPT_STAR:
		*E = append(*E, emi(c, commons.EMI_DMG_HEAT))
	case commons.GPT_ASTEROID:
		*E = append(*E, emi(c, commons.EMI_DMG_MECH))
	}
}

func newEmi(E *[]commons.Emission, c EmiChance, d EmiDistib) {
	var ok bool
	var t int
	for !ok {
		t = d.gen()
		if !has(*E, t) {
			ok = true
		}
	}
	*E = append(*E, emi(c, strconv.Itoa(t)))
}

func emi(c EmiChance, t string) commons.Emission {
	farRange := math.Round(commons.KDev(Opts.Dev) * c.Far)
	closeRange := math.Round(commons.KDev(Opts.Dev) * c.Close)
	if closeRange > farRange*0.9 {
		closeRange = farRange * 0.9
	}
	okr := func(x float64) float64 {
		const sgn = 100
		return float64(int(x*sgn)) / sgn
	}
	return commons.Emission{
		Type:      t,
		FarRange:  farRange,
		FarValue:  0,
		MainRange: closeRange,
		MainValue: okr(c.Force * commons.KDev(Opts.Dev)),
	}
}

func procSignatures(p *commons.GalaxyPoint) []commons.Signature {
	res := make([]commons.Signature, 0)
	//drop old damage sigs
	for _, s := range p.Signatures {
		if s.TypeName[0] != 'X' {
			res = append(res, s)
		}
	}

	for _, e := range p.Emissions {
		res = append(res, commons.Signature{
			TypeName:  "X" + e.Type,
			Dev:       okrV2(v2.RandomInCircle(1)),
			SigString: "",
		})
	}

	return res
}

func has(E []commons.Emission, t int) bool {
	for _, e := range E {
		if e.Type == strconv.Itoa(t) {
			return true
		}
	}
	return false
}

func okrV2(v v2.V2) v2.V2 {
	f := func(f float64) float64 {
		return math.Round(f*1000) / 1000
	}
	return v2.V2{X: f(v.X), Y: f(v.Y)}
}
