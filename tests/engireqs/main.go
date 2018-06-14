package main

import (
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/ranma"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"image/color"
)

const addr = "http://tagunil.ru:8000/system/"

var slotSprite *graph.Sprite
var dotSprite *graph.Sprite
var smokeSprite *graph.Sprite

var Ranma *ranma.Ranma

func run(image *ebiten.Image) error {

	procClick()

	if ebiten.IsRunningSlowly() {
		return nil
	}
	image.Fill(color.White)
	drawSlots(image)
	return nil
}

func main() {
	draw.InitTexAtlas()

	slotSprite = draw.NewAtlasSprite("BUILDING_BEACON", graph.NoCam)
	slotSprite.SetSize(50, 50)

	dotSprite = draw.NewAtlasSprite("trail", graph.NoCam)
	dotSprite.SetSize(60, 60)

	smokeSprite = draw.NewAtlasSprite("smoke", graph.NoCam)
	smokeSprite.SetSize(60, 60)

	Ranma = ranma.NewRanma(addr, true, 1000, 16)

	ebiten.Run(run, 1325, 725, 1, "Engi")
}

func drawSlots(image *ebiten.Image) {
	for n := 0; n < 8; n++ {
		for m := 0; m < 16; m++ {
			p := pos(n, m)
			if Ranma.GetInByte(n, m) {
				smokeSprite.SetPos(p)
				smokeSprite.Draw(image)
			}
			slotSprite.SetPos(p)
			slotSprite.Draw(image)

			if Ranma.GetOutByte(n, m) {
				dotSprite.SetPos(p)
				dotSprite.Draw(image)
			}
		}
	}
}

func procClick() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		v := v2.V2{X: float64(x), Y: float64(y)}
		for sn := 0; sn < 8; sn++ {
			for bn := 0; bn < 16; bn++ {
				d := pos(sn, bn).Sub(v).Len()
				if d < 30 {
					Ranma.XorInByte(sn, bn)
				}
			}
		}
	}
}

func pos(s, b int) v2.V2 {
	return v2.V2{100, 100}.AddMul(v2.V2{75, 0}, float64(b)).AddMul(v2.V2{0, 75}, float64(s))
}
