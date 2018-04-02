package graph

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
	"math"
)

const Deg2Rad = math.Pi / 180

type Point struct {
	X, Y float64
}

type Sprite struct {
	tex      Tex
	col, row int
	op       *ebiten.DrawImageOptions
	//before and past cam parts of geom
	dirty  bool
	g1, g2 ebiten.GeoM
	//pos and rot point, in sprite before scale
	//in pxls
	pivot Point
	//basic scale here
	sx, sy float64
	//place of center
	pos Point
	//in rad
	angle float64
	//alpha normed
	color color.Color
	//additional alpha k [0...1]
	alpha float64
	//if we plan to apply camera -- flip vertically, cz world coords Y is up and screen - dwn
	cam *Camera
	//FixedSize
	denyCamScale bool
}

func NewSprite(tex Tex, cam *Camera, denyCamScale bool) *Sprite {
	op := &ebiten.DrawImageOptions{}
	op.Filter = tex.filter
	srcRect := image.Rect(0, 0, tex.sw, tex.sh)
	op.SourceRect = &srcRect

	w, h := float64(tex.sw), float64(tex.sh)

	res := &Sprite{
		tex:          tex,
		op:           op,
		sx:           1,
		sy:           1,
		pivot:        Point{w / 2, h / 2},
		color:        color.White,
		alpha:        1,
		cam:          cam,
		denyCamScale: denyCamScale,
	}

	res.calcGeom()

	return res
}

func NewSpriteFromFile(filename string, filter ebiten.Filter, sw, sh int, cam *Camera, denyCamScale bool) (*Sprite, error) {
	tex, err := GetTex(filename, filter, sw, sh)
	if err != nil {
		return nil, err
	}
	return NewSprite(tex, cam, denyCamScale), nil
}

func (s *Sprite) recalcColorM() {
	const MaxColor = 0xffff
	s.op.ColorM.Reset()
	r, g, b, a := s.color.RGBA()
	s.op.ColorM.Scale(s.alpha*float64(r)/MaxColor, s.alpha*float64(g)/MaxColor, s.alpha*float64(b)/MaxColor, s.alpha*float64(a)/MaxColor)
}

func (s *Sprite) SetColor(color color.Color) {
	s.color = color
	s.dirty = true
}

func (s *Sprite) SetAlpha(a float64) {
	s.alpha = a
	s.dirty = true
}

func (s *Sprite) SetScale(x, y float64) {
	s.sx = x
	s.sy = y
	s.dirty = true
}

func (s *Sprite) SetSize(x, y float64) {
	s.sx = x / float64(s.tex.sh)
	s.sy = y / float64(s.tex.sw)
	s.dirty = true
}

func (s *Sprite) SetPivot(pivot Point) {
	s.pivot = pivot
	s.dirty = true
}

func (s *Sprite) SetPos(pos Point) {
	s.pos = pos
	s.dirty = true
}

func (s *Sprite) SetAng(angleDeg float64) {
	s.angle = angleDeg * Deg2Rad
	s.dirty = true
}

func (s *Sprite) SetPosAng(pos Point, angle float64) {
	s.pos = pos
	s.angle = angle
	s.dirty = true
}

func (s *Sprite) calcGeom() {
	G := ebiten.GeoM{}
	G.Translate(-s.pivot.X, -s.pivot.Y)
	G.Scale(s.sx, s.sy)
	s.g1 = G
	G.Reset()
	G.Rotate(s.angle)
	G.Translate(s.pos.X, s.pos.Y)
	s.g2 = G
}

//Copy options, so cam apply do not change
func (s *Sprite) ImageOp() (*ebiten.Image, *ebiten.DrawImageOptions) {
	if s.dirty {
		s.calcGeom()
		s.dirty = false
	}
	op := new(ebiten.DrawImageOptions)
	*op = *s.op
	G := s.g1
	//Flip vert before cam coords
	if s.cam != nil {
		G.Scale(1, -1)
		if s.denyCamScale {
			G.Scale(1/s.cam.Scale, 1/s.cam.Scale)
		}
	}
	G.Concat(s.g2)
	if s.cam != nil {
		G.Concat(s.cam.Geom())
	}
	op.GeoM = G
	return s.tex.image, op
}
