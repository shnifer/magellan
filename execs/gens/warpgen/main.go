package main

import (
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"io/ioutil"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Options struct {
	PosScale float64

	Size float64
	R10 float64
	Mass float64

	DevPercent float64

	RedR float64
	DR float64
	SpawnOutK float64
}

var Opts Options

func main() {
	buf, err := ioutil.ReadFile("starpos.json")
	if err != nil {
		panic(err)
	}
	var pts map[string]v2.V2
	err = json.Unmarshal(buf, &pts)
	if err != nil {
		panic(err)
	}

	dat, err := ioutil.ReadFile("warpgen_ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &Opts)

	var gal commons.Galaxy
	gal.Points = make(map[string]*commons.GalaxyPoint)
	for id, pt := range pts {
		kMass:=commons.KDev(Opts.DevPercent)
		kR:=commons.KDev(Opts.DevPercent)
		r := Opts.RedR*kR
		kR=commons.KDev(Opts.DevPercent)
		dr := Opts.DR*kR
		p := commons.GalaxyPoint{
			Pos:               okrV2(pt.Mul(Opts.PosScale)),
			Type:              commons.GPT_WARP,
			Size:              okr(Opts.Size*kMass),
			Mass:              okr(Opts.Mass*kMass),
			GDepth:            okr(Opts.R10/3*kMass),
			WarpRedOutDist:    okr(r),
			WarpGreenInDist:   okr(r + dr),
			WarpGreenOutDist:  okr(r + 2*dr),
			WarpYellowOutDist: okr(r + 3*dr),
			WarpSpawnDistance: okr((r + 3*dr)*Opts.SpawnOutK),
			GreenColor:        colornames.White,
			InnerColor:        randomColor(),
			OuterColor:        randomColor(),
			Color:             randomColor(),
		}
		gal.Points[id] = &p
	}
	res, err := json.Marshal(gal)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile("galaxy_warp.json", res, 0)
	if err != nil {
		panic(err)
	}
}

func okr(x float64) float64 {
	const sgn = 100
	return math.Floor(x*sgn) / sgn
}

func randomColor() color.RGBA {
	r := func() byte {
		return byte(rand.Intn(256))
	}
	return color.RGBA{
		R: r(),
		G: r(),
		B: r(),
		A: 255,
	}
}

func okrV2(v v2.V2) v2.V2{
	return v2.V2{
		X: okr(v.X),
		Y: okr(v.Y),
	}
}