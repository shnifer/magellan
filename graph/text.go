package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	et "github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"strings"
)

const interlinesK = 1.2

func TopLeft() v2.V2 {
	return v2.V2{X: 0.0, Y: 0.0}
}
func Center() v2.V2 {
	return v2.V2{X: 0.5, Y: 0.5}
}

func MiddleBottom() v2.V2 {
	return v2.V2{X: 0.5, Y: 1}
}

func BottomRight() v2.V2 {
	return v2.V2{X: 1, Y: 1}
}

type Text struct {
	//common
	text   []string
	face   font.Face
	strH   int
	color  color.Color
	w, h   int
	bounds image.Rectangle

	//for Draw method
	pos v2.V2
	//in parts of w/h [0...1]
	pivot v2.V2

	//for Image method
	filter ebiten.Filter
}

func NewText(text string, face font.Face, color color.Color) *Text {

	strs := strings.Split(text, "\n")
	var rect image.Rectangle
	strH := int(float64(face.Metrics().Height.Round()) * interlinesK)

	for i, str := range strs {
		b, _ := font.BoundString(face, str)
		r := image.Rect(b.Min.X.Round(), b.Min.Y.Round()+i*strH, b.Max.X.Round(), b.Max.Y.Round()+i*strH)
		rect = rect.Union(r)
	}

	res := Text{
		text:   strs,
		strH:   strH,
		face:   face,
		color:  color,
		bounds: rect,
		w:      rect.Dx(),
		h:      rect.Dy(),
		filter: ebiten.FilterDefault,
	}

	return &res
}

func (t *Text) SetPosPivot(pos, pivot v2.V2) {
	t.pos = pos
	t.pivot = pivot
}

func (t *Text) Draw(dst *ebiten.Image) {
	for i, str := range t.text {
		et.Draw(dst, str, t.face,
			int(t.pos.X-t.pivot.Y*float64(t.w))-t.bounds.Min.X,
			int(t.pos.Y-t.pivot.Y*float64(t.h))-t.bounds.Min.Y+i*t.strH, t.color)
	}
}

func (t *Text) SetFilter(filter ebiten.Filter) {
	t.filter = filter
}

func (t *Text) Image() *ebiten.Image {
	img, _ := ebiten.NewImage(t.w, t.h, t.filter)
	for i, str := range t.text {
		et.Draw(img, str, t.face, -t.bounds.Min.X, -t.bounds.Min.Y+i*t.strH, t.color)
	}
	return img
}
