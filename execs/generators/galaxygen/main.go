package main

import (
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/image/colornames"
	"io/ioutil"
	"os"
	"time"
	"github.com/Shnifer/magellan/v2"
	"strconv"
	"image/color"
)

type Options struct {
	N int

	MinRes    int
	MaxCap    int
	CloseDist int

	DotOutSize int
	DotInSize  int
}

var Opts Options

type StarPos map[string]v2.V2

func main() {
	stop := timer("ALL")
	defer stop()

	dat, err := ioutil.ReadFile("ini.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(dat, &Opts)

	backf, err := os.Open("back.png")
	if err != nil {
		panic(err)
	}
	defer backf.Close()
	back, _, err := image.Decode(backf)
	if err != nil {
		panic(err)
	}
	out := image.NewRGBA(back.Bounds())
	factDensity := image.NewRGBA(back.Bounds())
	draw.Draw(out, out.Bounds(), back, image.ZP, draw.Over)

	f, err := os.Open("density.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	density, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	DensityF := func(x int, y int) byte {
		R, _, _, _ := density.At(x, y).RGBA()
		r := int(R >> 8)

		res := 255 - r

		if Opts.MinRes > 0 {
			if res < Opts.MinRes {
				res = 0
			}
		} else {
			if res > 0 {
				res -= Opts.MinRes
			}
		}

		if res > Opts.MaxCap && Opts.MaxCap > 0 {
			res = Opts.MaxCap
		}

		if res > 255 {
			res = 255
		}
		return byte(res)
	}

	for x:=0;x<density.Bounds().Max.X;x++{
		for y:=0;y<density.Bounds().Max.Y;y++{
			d:=DensityF(x,y)
			factDensity.Set(x,y,color.RGBA{R:d, G:d, B:d, A:255})
		}
	}


	//rand.Seed(time.Now().Unix())
	RPG := CreateRandomPointGenerator(density.Bounds(), DensityF)

	stars := make([]image.Point, Opts.N)

	for i := range stars {
		stars[i] = RPG()
	}

	stars = deleteClose(stars)

	starPos:=make(StarPos, len(stars))
	for i, s := range stars {
		name:=genID()
		if i==0{
			name = "solar"
		}
		p:=s.Sub(density.Bounds().Max.Div(2))
		v:=v2.V2{X: float64(p.X), Y: float64(p.Y)}.Add(v2.RandomInCircle(1))
		starPos[name] = v
	}

	res, err := json.Marshal(starPos)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("starpos.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(res)

	kx := float64(back.Bounds().Max.X / density.Bounds().Max.X)
	ky := float64(back.Bounds().Max.Y / density.Bounds().Max.Y)
	log.Println("kx,ky ", kx, ky)

	OutSize := Opts.DotOutSize
	InSize := Opts.DotInSize

	var r2 int
	for _, star := range stars {
		X := int(kx * float64(star.X))
		Y := int(ky * float64(star.Y))
		for x := X - OutSize; x <= X+OutSize; x++ {
			for y := Y - OutSize; y <= Y+OutSize; y++ {
				r2 = (x-X)*(x-X) + (y-Y)*(y-Y)
				if r2 <= OutSize*OutSize && r2 > InSize*InSize {
					out.Set(x, y, colornames.Orange)
				}
			}
		}
	}

	outf, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(outf, out); err != nil {
		outf.Close()
		log.Fatal(err)
	}

	if err := outf.Close(); err != nil {
		log.Fatal(err)
	}

	factDensityf, err := os.Create("fact_density.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(factDensityf, factDensity); err != nil {
		factDensityf.Close()
		log.Fatal(err)
	}

	if err := factDensityf.Close(); err != nil {
		log.Fatal(err)
	}
}

func timer(caption string) func() {
	start := time.Now()
	return func() {
		fmt.Println(caption, time.Now().Sub(start))
	}
}

func deleteClose(stars []image.Point) (res []image.Point) {
	res = make([]image.Point, 0, len(stars))
	var f bool
	var v image.Point
	var r int
	for _, star := range stars {
		f = false
		for _, checkS := range res {
			v = star.Sub(checkS)
			if v.X*v.X+v.Y*v.Y <= Opts.CloseDist*Opts.CloseDist {
				f = true
				r++
				break
			}
		}
		if !f {
			res = append(res, star)
		}
	}

	log.Printf("Removed %v close stars. %v left", r, len(res))
	return res
}


var usedId map[string]struct{}
func init(){
	usedId = make(map[string]struct{})
}
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
