package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"time"
)

func update(img *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	img.DrawImage(Back, Op)
	text.Draw(img, "Text", Face, 100, 100, color.White)
	select {
	case <-Tick:
		fmt.Println(ebiten.CurrentFPS())
	default:

	}
	return nil
}

var Back *ebiten.Image
var Op *ebiten.DrawImageOptions
var Face font.Face
var Tick <-chan time.Time

func main() {
	Tick = time.Tick(time.Second)
	var err error
	Back, _, err = ebitenutil.NewImageFromFile("back.jpg", ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	Op = &ebiten.DrawImageOptions{}
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
	ebiten.Run(update, 320, 240, 1, "test")
}
