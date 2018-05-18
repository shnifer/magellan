package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	//"math/rand"
	"fmt"
	"os"
	"time"
)

const GalaxyRadius = 2000  //in light years
const HalfImageSize = 2000 //in pixels, better even

func main() {
	stop := timer("ALL")
	defer stop()

	//rand.Seed(time.Now().Unix())
	RPG := CreateRandomPointGenerator(Dens)

	stars := make([]image.Point, 1000)

	for i, _ := range stars {
		stars[i] = RPG()
	}

	r := image.Rect(-HalfImageSize, -HalfImageSize, HalfImageSize, HalfImageSize)
	img := image.NewRGBA(r)
	c := color.Black
	draw.Draw(img, r, &image.Uniform{c}, image.ZP, draw.Src)

	for x := -HalfImageSize; x < HalfImageSize; x++ {
		for y := -HalfImageSize; y < HalfImageSize; y++ {
			d := Dens(x, y)
			d /= 4

			img.Set(x, y, color.NRGBA{
				R: 0, G: 0, B: d, A: 255,
			})
		}
	}

	for _, star := range stars {
		x := star.X
		y := star.Y
		img.Set(x, y, color.NRGBA{
			R: 255, G: 255, B: 0, A: 255,
		})
	}

	f, err := os.Create("image.png")
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func timer(caption string) func() {
	start := time.Now()
	return func() {
		fmt.Println(caption, time.Now().Sub(start))
	}
}
