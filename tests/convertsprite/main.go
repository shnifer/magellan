package main

import (
	"bytes"
	_ "golang.org/x/image/bmp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

//change to non-transparent to white
func main() {
	if len(os.Args) == 1 {
		log.Println("use convertsprite filename.png")
		return
	}
	fn := os.Args[1]
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(b)
	img, _, err := image.Decode(buf)
	if err != nil {
		panic(err)
	}

	const scale = 2
	bounds := img.Bounds()
	bounds.Max = bounds.Max.Div(scale)
	res := image.NewRGBA(bounds)
	for x := 0; x < img.Bounds().Max.X/scale; x++ {
		for y := 0; y < img.Bounds().Max.Y/scale; y++ {
			src := img.At(x*scale, y*scale)
			_, _, _, a := src.RGBA()
			if a > 0 {
				//			src=color.White
			}
			res.Set(x, y, src)
		}
	}

	outf, err := os.Create("_" + fn)
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	if err := png.Encode(outf, res); err != nil {
		panic(err)
	}
}
