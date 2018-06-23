package main

import (
	"os"
	"log"
	"io/ioutil"
	"image"
	_"image/jpeg"
	_"image/gif"
	_"golang.org/x/image/bmp"
	"bytes"
	"image/png"
)

//change to non-transparent to white
func main(){
	if len(os.Args)==1 {
		log.Println("use convertsprite filename.png")
		return
	}
	fn:=os.Args[1]
	b,err:=ioutil.ReadFile(fn)
	if err!=nil{
		panic(err)
	}
	buf:=bytes.NewBuffer(b)
	img, _, err:=image.Decode(buf)
	if err!=nil{
		panic(err)
	}

	bounds:=img.Bounds()
	bounds.Max=bounds.Max.Div(2)
	res:=image.NewRGBA(bounds)
	for x:=0; x<img.Bounds().Max.X/2; x++{
		for y:=0; y<img.Bounds().Max.Y/2; y++{
			src:=img.At(x*2,y*2)
			_,_,_,a:=src.RGBA()
			if a>0 {
	//			src=color.White
			}
			res.Set(x,y,src)
		}
	}

	outf, err := os.Create("_"+fn)
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	if err := png.Encode(outf, res); err != nil {
		panic(err)
	}
}
