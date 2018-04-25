package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
	"math"
)

const Deg2Rad = math.Pi / 180

type Sprite struct {
	tex      Tex
	col, row int
	op       *ebiten.DrawImageOptions
	//before and past cam parts of geom
	dirty      bool
	colorDirty bool
	g1, g2     ebiten.GeoM
	//pos and rot point, in sprite before scale
	//in pxls
	pivot v2.V2
	//basic scale here
	sx, sy float64
	//place of center
	pos v2.V2
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
	//FixedAngle
	denyCamAngle bool
	//number of sprite from sheet
	spriteN int
}

func NewSprite(tex Tex, cam *Camera, denyCamScale, denyCamAngle bool) *Sprite {
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
		pivot:        v2.V2{X: w / 2, Y: h / 2},
		color:        color.White,
		alpha:        1,
		cam:          cam,
		denyCamScale: denyCamScale,
		denyCamAngle: denyCamAngle,
	}

	res.calcGeom()

	return res
}

//without cam
func NewSpriteHUD(tex Tex) *Sprite {
	return NewSprite(tex, nil, false, false)
}

func NewSpriteFromFile(filename string, filter ebiten.Filter, sw, sh int, count int, cam *Camera, denyCamScale, denyCamAngle bool) (*Sprite, error) {
	tex, err := GetTex(filename, filter, sw, sh, count)
	if err != nil {
		return nil, err
	}
	return NewSprite(tex, cam, denyCamScale, denyCamAngle), nil
}

func (s *Sprite) recalcColorM() {
	const MaxColor = 0xffff
	s.op.ColorM.Reset()
	r, g, b, a := s.color.RGBA()
	s.op.ColorM.Scale(s.alpha*float64(r)/MaxColor, s.alpha*float64(g)/MaxColor, s.alpha*float64(b)/MaxColor, s.alpha*float64(a)/MaxColor)
}

func (s *Sprite) SetColor(color color.Color) {
	s.color = color
	s.colorDirty = true
}

func (s *Sprite) SetAlpha(a float64) {
	s.alpha = a
	s.colorDirty = true
}

func (s *Sprite) SetScale(x, y float64) {
	s.sx = x
	s.sy = y
	s.dirty = true
}

func (s *Sprite) SetSize(x, y float64) {
	s.sx = x / float64(s.tex.sw)
	s.sy = y / float64(s.tex.sh)
	s.dirty = true
}

//pivotPartial is [0..1,0..1] vector of pivot point in parts of image size
func (s *Sprite) SetPivot(pivotPartial v2.V2) {
	s.pivot = v2.V2{
		X: pivotPartial.X * float64(s.tex.sw),
		Y: pivotPartial.Y * float64(s.tex.sh),
	}
	s.dirty = true
}

func (s *Sprite) SetPos(pos v2.V2) {
	s.pos = pos
	s.dirty = true
}

func (s *Sprite) SetAng(angleDeg float64) {
	s.angle = angleDeg * Deg2Rad
	s.dirty = true
}

func (s *Sprite) SetPosAng(pos v2.V2, angle float64) {
	s.pos = pos
	s.angle = angle * Deg2Rad
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
	if s.colorDirty {
		s.recalcColorM()
		s.colorDirty = false
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
		if s.denyCamAngle {
			G.Rotate(s.cam.AngleDeg * Deg2Rad)
		}
	}
	G.Concat(s.g2)
	if s.cam != nil {
		G.Concat(s.cam.Geom())
	}
	op.GeoM = G
	return s.tex.image, op
}

func (s *Sprite) Draw(dest *ebiten.Image) {
	img, op := s.ImageOp()
	dest.DrawImage(img, op)
}

func (s *Sprite) SpriteN() int {
	return s.spriteN
}

func (s *Sprite) SpritesCount() int {
	return s.tex.count
}

func (s *Sprite) SetSpriteN(n int) {
	n = n % s.tex.count
	s.spriteN = n
	nx := n % s.tex.cols
	ny := n / s.tex.cols
	rect := image.Rect(nx*s.tex.sw, ny*s.tex.sh, (nx+1)*s.tex.sw, (ny+1)*s.tex.sh)
	s.op.SourceRect = &rect
}

func (s *Sprite) NextSprite() {
	s.SetSpriteN(s.spriteN + 1)
}
