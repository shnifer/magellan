package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	et "github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image"
	"image/color"
)

func TopLeft() v2.V2 {
	return v2.V2{X: 0.0, Y: 0.0}
}
func Center() v2.V2 {
	return v2.V2{X: 0.5, Y: 0.5}
}

func MiddleBottom() v2.V2 {
	return v2.V2{X: 0.5, Y: 1}
}

type Text struct {
	//common
	text   string
	face   font.Face
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
	b, _ := font.BoundString(face, text)
	rect := image.Rect(b.Min.X.Round(), b.Min.Y.Round(), b.Max.X.Round(), b.Max.Y.Round())

	res := Text{
		text:   text,
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
	et.Draw(dst, t.text, t.face,
		int(t.pos.X-t.pivot.Y*float64(t.w))-t.bounds.Min.X,
		int(t.pos.Y-t.pivot.Y*float64(t.h))-t.bounds.Min.Y, t.color)
}

func (t *Text) SetFilter(filter ebiten.Filter) {
	t.filter = filter
}

func (t *Text) Image() *ebiten.Image {
	img, _ := ebiten.NewImage(t.w, t.h, t.filter)
	et.Draw(img, t.text, t.face, -t.bounds.Min.X, -t.bounds.Min.Y, t.color)
	return img
}
