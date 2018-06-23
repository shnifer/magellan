package draw

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
)

type SlidingSphere struct {
	*graph.Sprite

	temp          *ebiten.Image
	w, h          int
	slidingSprite *graph.SlidingSprite
	maskSprite    *graph.Sprite
	layer         int
	periodS       float64
}

func NewAtlasSlidingSphere(atlasName string, params graph.CamParams, layer int, periodS float64) *SlidingSphere {
	tex := GetAtlasTex(atlasName)
	tex = graph.SlidingTex(tex)
	underSprite := graph.NewSprite(tex, graph.NoCam)
	underSprite.SetPivot(graph.TopLeft())
	slidingSprite := graph.NewSlidingSprite(underSprite)
	w, h := tex.Size()
	temp, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	sprite := graph.NewSprite(tex, params)

	mask := graph.NewSprite(graph.CircleTex(), graph.NoCam)
	mask.SetPivot(graph.TopLeft())
	mask.SetSize(float64(w), float64(h))

	return &SlidingSphere{
		Sprite:        sprite,
		temp:          temp,
		w:             w,
		h:             h,
		slidingSprite: slidingSprite,
		maskSprite:    mask,
		layer:         layer,
		periodS:       periodS,
	}
}

func (ss *SlidingSphere) Req() (res *graph.DrawQueue) {
	res = graph.NewDrawQueue()

	ss.temp.Clear()
	ss.slidingSprite.Draw(ss.temp)
	img, op := ss.maskSprite.ImageOp()
	op.CompositeMode = ebiten.CompositeModeDestinationIn
	ss.temp.DrawImage(img, op)

	t := graph.TexFromImage(ss.temp, ebiten.FilterDefault, ss.w, ss.h, 0, "~slidingthing")
	ss.Sprite.SetTex(t)

	res.Add(ss.Sprite, ss.layer)

	return res
}

func (ss *SlidingSphere) Update(dt float64) {
	ss.AddSlide(dt / ss.periodS)
}

func (ss *SlidingSphere) SetSlide(slide float64) {
	ss.slidingSprite.SetSlide(slide)
}

func (ss *SlidingSphere) AddSlide(deltaSlide float64) {
	ss.slidingSprite.AddSlide(deltaSlide)
}
