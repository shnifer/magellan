package graph

import (
	"bytes"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"math"
)

const Deg2Rad = math.Pi / 180

//const Rad2Deg = 180 / math.Pi

type Sprite struct {
	tex       Tex
	op        ebiten.DrawImageOptions
	camParams CamParams
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
	//number of sprite from sheet
	spriteN int
}

func NewSprite(tex Tex, params CamParams) *Sprite {
	op := ebiten.DrawImageOptions{}
	op.Filter = tex.filter
	srcRect := image.Rect(0, 0, tex.sw, tex.sh)
	op.SourceRect = &srcRect

	w, h := float64(tex.sw), float64(tex.sh)

	res := &Sprite{
		tex:       tex,
		op:        op,
		sx:        1,
		sy:        1,
		pivot:     v2.V2{X: w / 2, Y: h / 2},
		color:     color.White,
		alpha:     1,
		camParams: params,
	}

	res.calcGeom()

	return res
}

//without cam
func NewSpriteHUD(tex Tex) *Sprite {
	return NewSprite(tex, NoCam)
}

func fileLoader(filename string) (io.Reader, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(b), err
}

func NewSpriteFromFile(filename string, smoothFilter bool, sw, sh int, count int, params CamParams) (*Sprite, error) {
	tex, err := GetTex(filename, smoothFilter, sw, sh, count, fileLoader)
	if err != nil {
		return nil, err
	}
	return NewSprite(tex, params), nil
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

func (s *Sprite) SetSizeProportion(size float64) {
	sx := size / float64(s.tex.sw)
	sy := size / float64(s.tex.sh)
	var scale float64
	if sx > sy {
		scale = sy
	} else {
		scale = sx
	}
	s.sx = scale
	s.sy = scale
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
	*op = s.op
	G := s.g1
	//Flip vert before cam coords
	if s.camParams.Cam != nil {
		G.Scale(1, -1)
		if s.camParams.DenyScale {
			G.Scale(1/s.camParams.Cam.Scale, 1/s.camParams.Cam.Scale)
		}
		if s.camParams.DenyAngle {
			G.Rotate(s.camParams.Cam.AngleDeg * Deg2Rad)
		}
	}
	G.Concat(s.g2)
	if s.camParams.Cam != nil {
		G.Concat(s.camParams.Cam.Geom())
	}
	op.GeoM = G
	return s.tex.image, op
}

func (s *Sprite) skipDrawCheck() bool {
	if s == nil {
		log.Println("Draw called for nil Sprite")
		return true
	}

	//Check for camera clipping
	if s.camParams.Cam != nil {
		w := float64(s.tex.sw) * s.sx
		h := float64(s.tex.sh) * s.sy
		inRect := s.camParams.Cam.RectInSpace(s.pos, w, h)
		if !inRect {
			return true
		}
	}
	return false
}

func (s *Sprite) Draw(dest *ebiten.Image) {
	if s.skipDrawCheck() {
		return
	}

	img, op := s.ImageOp()
	dest.DrawImage(img, op)
}

//MUST support multiple draw with different parameters
func (s *Sprite) DrawF() (drawF, string) {
	if s.skipDrawCheck() {
		return drawFZero, ""
	}
	//so we calc draw ops on s.DrawF() call not drawF resolve
	img, op := s.ImageOp()
	f := func(dest *ebiten.Image) {
		dest.DrawImage(img, op)
	}
	return f, s.tex.name
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

func (s *Sprite) TexImageDispose() {
	s.tex.image.Dispose()
}

func (s *Sprite) Cols() int {
	return s.tex.cols
}
func (s *Sprite) Rows() int {
	return s.tex.rows
}
