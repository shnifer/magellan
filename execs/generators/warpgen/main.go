package main

import (
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
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

	dat, err := ioutil.ReadFile("ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &Opts)

	var gal commons.Galaxy
	gal.Points = make(map[string]*commons.GalaxyPoint)
	for id, pt := range pts {
		p := commons.GalaxyPoint{
			Pos:               pt.Mul(Opts.Scale),
			Type:              commons.GPT_WARP,
			Size:              1,
			Mass:              okr(1 + rand.Float64()),
			GDepth: 0.1,
			WarpSpawnDistance: 5,
			WarpRedOutDist:    1,
			WarpGreenInDist:   2,
			WarpGreenOutDist:  3,
			WarpYellowOutDist: 4,
			GreenColor: colornames.Green,
			InnerColor: colornames.Firebrick,
			OuterColor: colornames.Lightyellow,
			Color:             colornames.White,
		}
		gal.Points[id] = &p
	}
	res, err := json.Marshal(gal)
	if err!=nil{
		panic(err)
	}
	err=ioutil.WriteFile("galaxy_warp.json", res, 0)
	if err!=nil{
		panic(err)
	}
}

func okr(x float64) float64 {
	const sgn = 100
	return math.Floor(x*sgn) / sgn
}
