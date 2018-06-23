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
	if s.sprite.skipDrawCheck() {
		return
	}

	img, op := s.ImageOp()
	dest.DrawImage(img, op)
}

//MUST support multiple draw with different parameters
func (s *SlidingSprite) DrawF() (drawF, string) {
	if s.sprite.skipDrawCheck() {
		return drawFZero, ""
	}
	//so we calc draw ops on s.DrawF() call not drawF resolve
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

func SlidingTex(source Tex) (result Tex) {
	result = source
	w, h := source.image.Size()
	newImage, _ := ebiten.NewImage(w+h, h, source.filter)
	op := &ebiten.DrawImageOptions{}
	newImage.DrawImage(source.image, op)
	rect := image.Rect(0, 0, source.sw, h)
	op.SourceRect = &rect
	op.GeoM.Translate(float64(w), 0)
	newImage.DrawImage(source.image, op)
	result.image = newImage
	return result
}
