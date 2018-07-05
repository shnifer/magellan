package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

const defaultLineLen = 10
const defaultLineWid = 1

var defLine Tex

func init() {
	defLine = defaultLineTex()
}

//TODO: optimize! Do not create each time or trackpredictor will die
func Line(cam *Camera, from, to v2.V2, clr color.Color) *Sprite {
	return LineScr(cam.Apply(from), cam.Apply(to), clr)
}

func LineHUD(cam *Camera, from v2.V2, v V2, clr color.Color) *Sprite {
	pos := cam.Apply(from)
	return LineScr(pos, pos.Add(v), clr)
}

func LineScr(from, to v2.V2, clr color.Color) *Sprite {
	v := to.Sub(from)
	l := v.Len()
	a := v.Dir() + 180

	res := NewSpriteHUD(defLine)
	res.SetColor(clr)
	res.SetPivot(BotLeft())
	res.SetPos(from)
	res.SetScale(1, l/defaultLineLen)
	res.SetAng(a)
	return res
}

func defaultLineTex() Tex {
	img, _ := ebiten.NewImage(defaultLineWid, defaultLineLen, ebiten.FilterDefault)
	p := make([]byte, defaultLineLen*defaultLineWid*4)
	for i := range p {
		p[i] = 255
	}

	img.ReplacePixels(p)
	return TexFromImage(img, ebiten.FilterDefault, 0, 0, 0, "~line")
}
