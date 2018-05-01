package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	_"image/jpeg"
	"fmt"
	"golang.org/x/image/font"
	"io/ioutil"
	"github.com/golang/freetype/truetype"
	"image/color"
)

func update(img *ebiten.Image) error{
	if ebiten.IsRunningSlowly() {
		return nil
	}
	img.DrawImage(Back,Op)
	msg:=fmt.Sprintf("FPS: %v",ebiten.CurrentFPS())
	text.Draw(img, "Text", Face,100,100,color.White)
	ebitenutil.DebugPrint(img, msg)
	return nil
}

var Back *ebiten.Image
var Op *ebiten.DrawImageOptions
var Face font.Face

func main(){
	var err error
	Back,_,err=ebitenutil.NewImageFromFile("back.jpg",ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	Op=&ebiten.DrawImageOptions{}
	b, err := ioutil.ReadFile("interdim.ttf")
	if err != nil {
		panic(err)
	}
	f, err := truetype.Parse(b)
	if err != nil {
		panic(err)
	}
	Face = truetype.NewFace(f, &truetype.Options{Size: 20})
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(update,320,240,1,"test")
}
