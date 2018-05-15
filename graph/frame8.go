package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
)

type V2 = v2.V2

type Frame9HUD struct {
	*Sprite

	sprite      [9]*Sprite
	scale       float64
	actualScale float64
	w, h        float64

	//counted once from sprites
	//base minimals
	//top, left, right, bot float64
}

func NewFrame9(sprites [9]*Sprite, w, h float64) (res *Frame9HUD) {
	res = &Frame9HUD{
		sprite: sprites,
		scale:  1,
		w:      w,
		h:      h,
	}
	res.setPivots()
	res.recalc()
	return res
}

func (f9 *Frame9HUD) MinSize() (w, h float64) {
	max := func(a, b float64) float64 {
		if a > b {
			return a
		} else {
			return b
		}
	}
	return max(f9.wN(0)+f9.wN(2), f9.wN(6)+f9.wN(8)),
		max(f9.hN(0)+f9.hN(6), f9.hN(2)+f9.hN(8))
}

func (f9 *Frame9HUD) SetSize(w, h float64) {
	f9.w, f9.h = w, h
	f9.recalc()
}
func (f9 *Frame9HUD) SetScale(scale float64) {
	f9.scale = scale
	f9.recalc()
}

func (f9 *Frame9HUD) recalc() {
	minW, minH := f9.MinSize()
	actualS := f9.scale
	if minW > 0 {
		maxS := f9.w / minW
		if actualS > maxS {
			actualS = maxS
		}
	}
	if minH > 0 {
		maxS := f9.h / minH
		if actualS > maxS {
			actualS = maxS
		}
	}
	f9.actualScale = actualS

	for i, sprite := range f9.sprite {
		if sprite == nil {
			continue
		}
		switch i {
		case 0:
			sprite.SetScale(actualS, actualS)
			sprite.SetPos(V2{X: 0, Y: 0})
		case 1:
			sprite.SetSize(f9.w-actualS*(f9.wN(0)+f9.wN(2)), float64(sprite.tex.sh)*actualS)
			sprite.SetPos(V2{X: f9.wN(0) * actualS, Y: 0})
		case 2:
			sprite.SetScale(actualS, actualS)
			sprite.SetPos(V2{X: f9.w, Y: 0})
		case 3:
			sprite.SetSize(float64(sprite.tex.sw)*actualS, f9.h-actualS*(f9.hN(0)+f9.hN(6)))
			sprite.SetPos(V2{X: 0, Y: f9.hN(0) * actualS})
		case 4:
			//skip middle
		case 5:
			sprite.SetSize(float64(sprite.tex.sw)*actualS, f9.h-actualS*(f9.hN(3)+f9.hN(8)))
			sprite.SetPos(V2{X: f9.w, Y: f9.hN(2) * actualS})
		case 6:
			sprite.SetScale(actualS, actualS)
			sprite.SetPos(V2{X: 0, Y: f9.h})
		case 7:
			sprite.SetSize(f9.w-actualS*(f9.wN(6)+f9.wN(8)), float64(sprite.tex.sh)*actualS)
			sprite.SetPos(V2{X: f9.wN(6) * actualS, Y: f9.h})
		case 8:
			sprite.SetScale(actualS, actualS)
			sprite.SetPos(V2{X: f9.w, Y: f9.h})
		}
	}
	image, _ := ebiten.NewImage(int(f9.w), int(f9.h), ebiten.FilterDefault)
	for i := 0; i < 9; i++ {
		if f9.sprite[i] != nil {
			f9.sprite[i].Draw(image)
		}
	}

	tex := TexFromImage(image, ebiten.FilterDefault, 0, 0, 1, "")
	f9.Sprite = NewSpriteHUD(tex)
	f9.Sprite.SetPivot(TopLeft())
}

/*
func (f9 *Frame9HUD) baseBorders() {
	var top, left, right, bot float64
	set := func(x *float64, v int) {
		if *x == 0 {
			*x = float64(v)
		}
	}

	if f9.sprite[0] != nil {
		set(&top, f9.sprite[0].tex.sh)
		set(&left, f9.sprite[0].tex.sw)
	}
	if f9.sprite[8] != nil {
		set(&bot, f9.sprite[8].tex.sh)
		set(&right, f9.sprite[8].tex.sw)
	}
	if f9.sprite[2] != nil {
		set(&top, f9.sprite[2].tex.sh)
		set(&right, f9.sprite[2].tex.sw)
	}
	if f9.sprite[6] != nil {
		set(&bot, f9.sprite[6].tex.sh)
		set(&left, f9.sprite[6].tex.sw)
	}
	if f9.sprite[1] != nil {
		set(&top, f9.sprite[1].tex.sh)
	}
	if f9.sprite[3] != nil {
		set(&left, f9.sprite[3].tex.sw)
	}
	if f9.sprite[5] != nil {
		set(&right, f9.sprite[5].tex.sw)
	}
	if f9.sprite[7] != nil {
		set(&bot, f9.sprite[7].tex.sh)
	}

	f9.top, f9.left, f9.right, f9.bot = top, left, right, bot
}
*/

/*
func (f9 *Frame9HUD) GetCenter() (center V2, w,h float64){
	w = f9.w-f9.left-f9.right
	h = f9.h-f9.top-f9.bot
	center = V2{X: f9.left+w/2, Y:f9.top+h/2}
	return center,w,h
}
*/

func (f9 *Frame9HUD) setPivots() {
	set := func(i int, pivot V2) {
		if f9.sprite[i] != nil {
			f9.sprite[i].SetPivot(pivot)
		}
	}
	set(0, TopLeft())
	set(1, TopLeft())
	set(2, TopRight())
	set(3, TopLeft())
	set(4, TopLeft())
	set(5, TopRight())
	set(6, BotLeft())
	set(7, BotLeft())
	set(8, BotRight())
}

func (f9 *Frame9HUD) wN(n int) float64 {
	if f9.sprite[n] == nil {
		return 0
	}
	return float64(f9.sprite[n].tex.sw)
}

func (f9 *Frame9HUD) hN(n int) float64 {
	if f9.sprite[n] == nil {
		return 0
	}
	return float64(f9.sprite[n].tex.sh)
}
