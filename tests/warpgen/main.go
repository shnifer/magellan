package main

import (
	"encoding/json"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image"
	"io/ioutil"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const scale = 1

func init() {
	rand.Seed(time.Now().UnixNano())
	usedId = make(map[string]struct{}, 0)
}

func main() {
	buf, err := ioutil.ReadFile("starpos.json")
	if err != nil {
		panic(err)
	}
	var pts []image.Point
	err = json.Unmarshal(buf, &pts)
	if err != nil {
		panic(err)
	}
	var gal commons.Galaxy
	gal.Points = make(map[string]*commons.GalaxyPoint)
	var flag bool
	for _, pt := range pts {
		p := commons.GalaxyPoint{
			Pos:               pos(pt).Mul(scale),
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
		id := genID()
		if !flag {
			flag = true
			id = "solar"
		}
		gal.Points[id] = &p
	}
	res, err := json.Marshal(gal)
	ioutil.WriteFile("galaxy_warp.json", res, 0)
}

func pos(pt image.Point) v2.V2 {
	v := v2.V2{X: float64(pt.X), Y: float64(pt.Y)}.Add(v2.RandomInCircle(1))
	v.X = okr(v.X)
	v.Y = okr(v.Y)
	return v
}

var usedId map[string]struct{}

func genID() string {
	for {
		res := randLetter() + randLetter() + strconv.Itoa(rand.Intn(10))
		if _, exist := usedId[res]; !exist {
			usedId[res] = struct{}{}
			return res
		}
	}
}

func randLetter() string {
	n := byte(rand.Intn(26))
	s := []byte("A")[0]
	return string([]byte{s + n})
}

func okr(x float64) float64 {
	const sgn = 100
	return math.Floor(x*sgn) / sgn
}
