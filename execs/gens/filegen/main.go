package main

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	hls "github.com/gerow/go-color"
	"image/color"
	"io/ioutil"
	"math/rand"
	"strconv"
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

type Options struct {
	StarANCount        int
	SizeMassDevPercent float64
	OrbitDevPercent    float64

	SingleStar struct {
		R10  float64
		Size float64
		MaxG float64
	}

	DoubleStar struct {
		R10    float64
		Size   float64
		MaxG   float64
		Radius float64
		Period float64
	}

	TripleStar struct {
		R10    float64
		Size   float64
		MaxG   float64
		Radius float64
		Period float64
		Pair   struct {
			R10    float64
			Size   float64
			MaxG   float64
			Radius float64
			Period float64
		}
	}
}

var Opts Options

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

	for sysName, stat := range warpP {
		pref := sysName + "-"
		planets := allPlanet[sysName]
		_ = planets
		points := make(map[string]*GalaxyPoint)

		createStars(stat, points, pref)

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

func createStars(stat WarpStat, points map[string]*GalaxyPoint, pref string) {
	switch stat.StarCount {
	case 1:
		points[pref+"S"] = pOpts{
			t:    GPT_STAR,
			r10:  Opts.SingleStar.R10,
			size: Opts.SingleStar.Size,
			maxG: Opts.SingleStar.MaxG,
		}.gp()
	case 2:
		points[pref+"sv"] = &GalaxyPoint{IsVirtual: true}

		kOrbitPeriod := KDev(Opts.OrbitDevPercent)
		r := Opts.DoubleStar.Radius * kOrbitPeriod
		period := Opts.DoubleStar.Period * kOrbitPeriod

		kr := KDev(Opts.OrbitDevPercent)

		points[pref+"S1"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv",
			orbit:  r * kr,
			period: period,
			phase:  0,
			r10:    Opts.DoubleStar.R10 / kr,
			size:   Opts.DoubleStar.Size / kr,
			maxG:   Opts.DoubleStar.MaxG / kr,
		}.gp()
		points[pref+"S2"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv",
			orbit:  r / kr,
			period: period,
			phase:  180,
			r10:    Opts.DoubleStar.R10 * kr,
			size:   Opts.DoubleStar.Size * kr,
			maxG:   Opts.DoubleStar.MaxG * kr,
		}.gp()
	case 3:
		points[pref+"sv1"] = &GalaxyPoint{IsVirtual: true}

		kOrbitPeriod := KDev(Opts.OrbitDevPercent)
		r := Opts.TripleStar.Radius * kOrbitPeriod
		period := Opts.TripleStar.Period * kOrbitPeriod

		kr := KDev(Opts.OrbitDevPercent)

		points[pref+"S1"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv1",
			orbit:  r * kr,
			period: period,
			phase:  0,
			r10:    Opts.DoubleStar.R10 / kr,
			size:   Opts.DoubleStar.Size / kr,
			maxG:   Opts.DoubleStar.MaxG / kr,
		}.gp()
		points[pref+"sv2"] = &GalaxyPoint{
			IsVirtual: true,
			ParentID:  pref + "sv1",
			Orbit:     r / kr,
			Period:    period,
			AngPhase:  180,
		}

		kOrbitPeriod = KDev(Opts.OrbitDevPercent)
		r = Opts.TripleStar.Pair.Radius * kOrbitPeriod
		period = Opts.TripleStar.Pair.Period * kOrbitPeriod

		kr = KDev(Opts.OrbitDevPercent)
		points[pref+"S2"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv2",
			orbit:  r * kr,
			period: period,
			phase:  0,
			r10:    Opts.TripleStar.Pair.R10 / kr,
			size:   Opts.TripleStar.Pair.Size / kr,
			maxG:   Opts.TripleStar.Pair.MaxG / kr,
		}.gp()
		points[pref+"S3"] = pOpts{
			t:      GPT_STAR,
			parent: pref + "sv2",
			orbit:  r / kr,
			period: period,
			phase:  180,
			r10:    Opts.TripleStar.Pair.R10 * kr,
			size:   Opts.TripleStar.Pair.Size * kr,
			maxG:   Opts.TripleStar.Pair.MaxG * kr,
		}.gp()
	}
}

func sAN(t string, count int) string {
	if count == 0 {
		return ""
	} else {
		return t + "-" + strconv.Itoa(rand.Intn(count))
	}
}

type pOpts struct {
	parent string
	t      string
	orbit  float64
	period float64
	phase  float64
	size   float64
	r10    float64
	maxG   float64
	shps   [15]int
}

func (o pOpts) gp() *GalaxyPoint {
	okr := func(x float64) float64 {
		const sgn = 100
		return float64(int(x*sgn)) / sgn
	}

	count := 0
	switch o.t {
	case GPT_STAR:
		count = Opts.StarANCount
	}

	massSizeK := KDev(Opts.SizeMassDevPercent)

	zd := o.r10 / 3 * massSizeK
	maxG := o.maxG * massSizeK
	mass := maxG * zd * zd

	signatures := sphs2sigs(o.shps)

	return &GalaxyPoint{
		ParentID:   o.parent,
		Type:       o.t,
		SpriteAN:   sAN(o.t, count),
		Orbit:      o.orbit,
		Period:     o.period,
		AngPhase:   o.phase,
		Size:       okr(o.size * massSizeK),
		Mass:       okr(mass),
		GDepth:     okr(zd),
		Emissions:  nil,
		Signatures: signatures,
		Color:      randBright(),
	}
}

func randBright() color.RGBA {
	rgb := hls.HSL{
		S: 0.5 + 0.5*rand.Float64(),
		L: 0.8 + 0.2*rand.Float64(),
		H: rand.Float64(),
	}.ToRGB()
	return color.RGBA{
		R: uint8(rgb.R * 255),
		G: uint8(rgb.G * 255),
		B: uint8(rgb.B * 255),
		A: 255,
	}
}

func sphs2sigs(s [15]int) []Signature {
	res := make([]Signature, 0)

	add := func(a, b int) {
		res = append(res, Signature{
			TypeName: strconv.Itoa(a) + "-" + strconv.Itoa(b),
			Dev:      v2.RandomInCircle(1),
		})
	}

	for i, v := range s {
		if v == NONE {
			continue
		}
		switch i {
		case MAGNET, RADIATIONBELT, OXYGEN, OZONE, ION, COREVEL, VULCAN, BIO:
			add(i, v)
		case ATMOMETALS, GASES, PEDOMETALS, COREMADE, LITOMETAL, MIXTURES:
			if v == EARTHANDNEW {
				add(i, EARTH)
				add(i, NEW)
			} else {
				add(i, v)
			}
		case WATER:
			switch v {
			case HARDANDGASOUS:
				add(i, HARD)
				add(i, GASOUS)
			case HARDANDLIQUID:
				add(i, HARD)
				add(i, LIQUID)
			case LIQUIDANDGASOUS:
				add(i, LIQUID)
				add(i, GASOUS)
			case HARDANDLIQUIDANDGASOUS:
				add(i, HARD)
				add(i, LIQUID)
				add(i, GASOUS)
			default:
				add(i, v)
			}
		}
	}

	return res
}

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
