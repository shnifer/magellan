package graph

import (
	et "github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"image/color"
	"image"
	"github.com/hajimehoshi/ebiten"
)

func TopLeft() Point{
	return Point{0.0,0.0}
}
func Center() Point{
	return Point{0.5,0.5}
}

type Text struct {
	//common
	text string
	face font.Face
	color color.Color
	w,h int
	bounds image.Rectangle

	//for Draw method
	pos Point
	//in parts of w/h [0...1]
	pivot Point

	//for Image method
	filter ebiten.Filter
}

func NewText(text string, face font.Face, color color.Color) *Text{
	b,_ := font.BoundString(face, text)
	rect:=image.Rect(b.Min.X.Round(),b.Min.Y.Round(),b.Max.X.Round(),b.Max.Y.Round())

	res:=Text{
		text:text,
		face:face,
		color:color,
		bounds:rect,
		w:rect.Dx(),
		h:rect.Dy(),
		filter:ebiten.FilterDefault,
	}

	return &res
}

func (t *Text) SetPosPivot(pos,pivot Point){
	t.pos=pos
	t.pivot=pivot
}

func (t *Text) Draw (dst *ebiten.Image){
	et.Draw(dst, t.text, t.face,
		int(t.pos.X-t.pivot.Y*float64(t.w))-t.bounds.Min.X,
		int(t.pos.Y-t.pivot.Y*float64(t.h))-t.bounds.Min.Y, t.color)
}

func (t *Text) SetFilter(filter ebiten.Filter){
	t.filter = filter
}

func (t *Text) Image() *ebiten.Image{
	img,_:=ebiten.NewImage(t.w,t.h,t.filter)
	et.Draw(img,t.text,t.face, -t.bounds.Min.X, -t.bounds.Min.Y, t.color)
	return img
}