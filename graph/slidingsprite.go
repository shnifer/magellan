package graph

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"math"
)

type SlidingSprite struct {
	sprite *Sprite
	//slide value, 0..1 is full horizontal round
	slide float64
	l     float64
}

func NewSlidingSprite(sprite *Sprite) *SlidingSprite {
	l := float64(sprite.tex.image.Bounds().Max.X - sprite.tex.sw)
	return &SlidingSprite{
		sprite: sprite,
		l:      l,
	}
}

func (s *SlidingSprite) normSlide() {
	s.slide = s.slide - math.Floor(s.slide)
}

//cut integer part from slide, so you can pass anything
func (s *SlidingSprite) SetSlide(slide float64) {
	s.slide = slide
	s.normSlide()
}

func (s *SlidingSprite) AddSlide(deltaSlide float64) {
	s.slide += deltaSlide
	s.normSlide()
}

func (s *SlidingSprite) Draw(dest *ebiten.Image) {
	if s.sprite.SkipDrawCheck() {
		return
	}

	img, op := s.ImageOp()
	dest.DrawImage(img, op)
}

//MUST support multiple draw with different parameters
func (s *SlidingSprite) DrawF() (DrawF, string) {
	if s.sprite.SkipDrawCheck() {
		return DrawFZero, ""
	}
	//so we calc draw ops on s.DrawF() call not DrawF resolve
	img, op := s.ImageOp()
	f := func(dest *ebiten.Image) {
		dest.DrawImage(img, op)
	}
	return f, s.sprite.tex.name
}

//Copy options, so cam apply do not change
func (s *SlidingSprite) ImageOp() (*ebiten.Image, *ebiten.DrawImageOptions) {
	img, op := s.sprite.ImageOp()
	sr := image.Rect(0, 0, s.sprite.tex.sw, s.sprite.tex.sh)
	move := image.Pt(int(s.l*s.slide), 0)
	sr = sr.Add(move)
	op.SourceRect = &sr
	return img, op
}
