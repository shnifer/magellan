package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
)

//Array: many sprites of the same texture with timeout
type part struct {
	sprite *Sprite

	isActive bool

	leftTime float64
}

func (sa *FadingArray) newPart(elem ArrayElem) (p part) {
	p.isActive = true
	p.leftTime = elem.LifeTime
	sprite := NewSprite(sa.tex, sa.cam, sa.denyScale, sa.denyAngle)
	sprite.SetPosAng(elem.Pos, elem.AngDeg)
	sprite.SetSize(elem.Size, elem.Size)
	sprite.SetSpriteN(elem.SpriteN)
	p.sprite = sprite
	return p
}

type ArrayElem struct {
	Pos      v2.V2
	Size     float64
	AngDeg   float64
	SpriteN  int
	LifeTime float64
}

type FadingArray struct {
	tex                  Tex
	cam                  *Camera
	denyScale, denyAngle bool
	array                []part

	cap  int
	cur  int
	used int
}

func NewFadingArray(tex Tex, cap int, cam *Camera, denyCamScale, denyCamAngle bool) (res *FadingArray) {
	res = new(FadingArray)
	res.tex = tex
	res.cap = cap
	res.cam = cam
	res.denyScale = denyCamScale
	res.denyAngle = denyCamAngle
	res.array = make([]part, cap, cap)
	return res
}

func (sa *FadingArray) findNextInd() int {

	if sa.used == sa.cap {
		sa.used--
		sa.cur = (sa.cur + 1) % sa.cap
		return sa.cur
	}
	for sa.array[sa.cur].isActive {
		sa.cur = (sa.cur + 1) % sa.cap
	}
	return sa.cur
}

func (sa *FadingArray) Add(elem ArrayElem) {
	ind := sa.findNextInd()
	sa.array[ind] = sa.newPart(elem)
	sa.used++
}

func (sa *FadingArray) Update(dt float64) {
	for i := range sa.array {
		if !sa.array[i].isActive {
			continue
		}
		sa.array[i].leftTime -= dt
		if sa.array[i].leftTime < 0 {
			sa.array[i].leftTime = 0
			sa.array[i].isActive = false
			sa.used--
			continue
		}
	}
}

func (sa *FadingArray) Draw(img *ebiten.Image) {
	for _, part := range sa.array {
		if !part.isActive {
			continue
		}
		part.sprite.Draw(img)
	}
}

func (sa *FadingArray) DrawF() (drawF, string) {
	return sa.Draw, sa.tex.name
}

func (sa *FadingArray) Clear() {
	sa.used = 0
	sa.cur = 0
	for i := range sa.array {
		sa.array[i].isActive = false
		sa.array[i].leftTime = 0
	}
}
