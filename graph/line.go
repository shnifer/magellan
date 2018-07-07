package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
)

const defaultLineLen = 10
const defaultLineWid = 1

var lineSprite *Sprite

func init() {
	defLine := defaultLineTex()
	lineSprite = NewSpriteHUD(defLine)

}

func Line(Q *DrawQueue, cam *Camera, from, to v2.V2, clr color.Color, layer int) {
	LineScr(Q, cam.Apply(from), cam.Apply(to), clr, layer)
}

func LineHUD(Q *DrawQueue, cam *Camera, from v2.V2, v V2, clr color.Color, layer int) {
	pos := cam.Apply(from)
	LineScr(Q, pos, pos.Add(v), clr, layer)
}

func LineScr(Q *DrawQueue, from, to v2.V2, clr color.Color, layer int) {

	if clipped(from, to) {
		return
	}
	v := to.Sub(from)
	l := v.Len()
	a := v.Dir() + 180

	lineSprite.SetColor(clr)
	lineSprite.SetPivot(BotLeft())
	lineSprite.SetPos(from)
	lineSprite.SetScale(1, l/defaultLineLen)
	lineSprite.SetAng(a)

	Q.Add(lineSprite, layer)
}

func clipped(from, to v2.V2) bool {
	if (int(from.X) == int(to.X)) &&
		(int(from.Y) == int(to.Y)) {
		//do not draw 1 pixel lines
		//MB bad for some cases but we are fine
		return true
	}
	if from.X >= 0 && from.X <= winW && from.Y >= 0 && from.Y <= winH {
		return false
	}
	if to.X >= 0 && to.X <= winW && to.Y >= 0 && to.Y <= winH {
		return false
	}
	if (from.X < 0 && to.X < 0) || (from.X > winW && to.X > winW) ||
		(from.Y < 0 && to.Y < 0) || (from.Y > winH && to.Y > winH) {
		return true
	}
	return false
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
