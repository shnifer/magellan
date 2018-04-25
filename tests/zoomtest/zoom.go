package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"image"
	_ "image/jpeg"
	_ "time"
)

var tex *ebiten.Image

func update(window *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			op := &ebiten.DrawImageOptions{}
			r := image.Rect(0, 0, 400, 300)
			_ = r
			//op.SourceRect = &r
			op.GeoM.Scale(0.5, 0.5)
			op.GeoM.Translate(400*float64(x), 300*float64(y))
			window.DrawImage(tex, op)
		}
	}
	/*
		op:=&ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.5,1.5)
		window.DrawImage(tex,op)
	*/
	msg := fmt.Sprintf("FPS: %v", ebiten.CurrentFPS())
	ebitenutil.DebugPrint(window, msg)
	return nil
}

func main() {

	img, _, err := ebitenutil.NewImageFromFile("res/textures/background.jpg", ebiten.FilterDefault)
	if err != nil {
		img, _, err = ebitenutil.NewImageFromFile("background.jpg", ebiten.FilterDefault)
		if err != nil {
			panic(err)
		}
	}
	tex = img
	ebiten.Run(update, 1200, 900, 1, "test")
}
