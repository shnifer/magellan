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
	Scale float64
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
		r := 2 + 2*rand.Float64()
		dr := r * (0.3 + 0.4*rand.Float64())
		p := commons.GalaxyPoint{
			Pos:               pt.Mul(Opts.Scale),
			Type:              commons.GPT_WARP,
			Size:              3,
			Mass:              okr(1 + rand.Float64()),
			GDepth:            0.1,
			WarpSpawnDistance: 5,
			WarpRedOutDist:    r,
			WarpGreenInDist:   r + dr,
			WarpGreenOutDist:  r + 2*dr,
			WarpYellowOutDist: r + 3*dr,
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