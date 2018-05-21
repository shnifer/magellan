package main

import (
	"image"
	"image/draw"
	"image/png"
	_"image/jpeg"
	"log"
	//"math/rand"
	"fmt"
	"os"
	"time"
	"golang.org/x/image/colornames"
)

func main() {
	stop := timer("ALL")
	defer stop()

	backf,err:=os.Open("back.png")
	if err!=nil{
		panic(err)
	}
	defer backf.Close()
	back, _, err:=image.Decode(backf)
	if err!=nil{
		panic(err)
	}
	out := image.NewRGBA(back.Bounds())
	draw.Draw(out,out.Bounds(),back,image.ZP,draw.Over)

	f,err:=os.Open("density.png")
	if err!=nil{
		panic(err)
	}
	defer f.Close()

	density, _, err:=image.Decode(f)
	if err!=nil{
		panic(err)
	}

	const minDensity = 0

	DensityF:= func(x int, y int) byte{
		R,_,_,_ := density.At(x,y).RGBA()
		r := R>>8

		res:=255-byte(r)

		if res<minDensity{
			res=0
		}
		return res
	}


	//rand.Seed(time.Now().Unix())
	RPG := CreateRandomPointGenerator(density.Bounds(),DensityF)

	stars := make([]image.Point, 1000)

	for i := range stars {
		stars[i] = RPG()
	}

	kx:=float64(back.Bounds().Max.X/density.Bounds().Max.X)
	ky:=float64(back.Bounds().Max.Y/density.Bounds().Max.Y)
	log.Println("kx,ky ",kx,ky)

	const dotOutSize = 5
	const dotInnerSize = 2

	var r2 int
	for _,star:=range stars{
		X:=int(kx*float64(star.X))
		Y:=int(ky*float64(star.Y))
		for x:=X-dotOutSize;x<=X+dotOutSize;x++{
			for y:=Y-dotOutSize;y<=Y+dotOutSize;y++{
				r2=(x-X)*(x-X)+(y-Y)*(y-Y)
				if r2<=dotOutSize*dotOutSize && r2>dotInnerSize*dotInnerSize {
					out.Set(x,y,colornames.Orange)
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
}

func timer(caption string) func() {
	start := time.Now()
	return func() {
		fmt.Println(caption, time.Now().Sub(start))
	}
}
