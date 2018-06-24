package draw

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
)

var mixerTemp *ebiten.Image

func init() {
	mixerTemp, _ = ebiten.NewImage(1024, 1024, ebiten.FilterDefault)
}

type SlidingSphere struct {
	*graph.Sprite

	w, h          int
	slidingSprite *graph.SlidingSprite
	maskSprite    *graph.Sprite
	periodS       float64
}

func NewAtlasSlidingSphere(atlasName string, params graph.CamParams, periodS float64) *SlidingSphere {
	tex := GetSlidingAtlasTex(atlasName)
	underSprite := graph.NewSprite(tex, graph.NoCam)
	underSprite.SetPivot(graph.TopLeft())
	slidingSprite := graph.NewSlidingSprite(underSprite)
	w, h := tex.Size()
	//mixerTemp, _ := ebiten.NewImage(w, h, ebiten.FilterDefault)
	sprite := graph.NewSprite(tex, params)

	mask := graph.NewSprite(graph.CircleTex(), graph.NoCam)
	mask.SetPivot(graph.TopLeft())
	mask.SetSize(float64(w), float64(h))

	return &SlidingSphere{
		Sprite: sprite,
		//	mixerTemp:          mixerTemp,
		w:             w,
		h:             h,
		slidingSprite: slidingSprite,
		maskSprite:    mask,
		periodS:       periodS,
	}
}

func (ss *SlidingSphere) DrawF() (graph.DrawF, string) {
	if ss.Sprite.SkipDrawCheck() {
		return graph.DrawFZero, ""
	}

	f := func(dest *ebiten.Image) {
		mixerTemp.Clear()
		ss.slidingSprite.Draw(mixerTemp)
		img, op := ss.maskSprite.ImageOp()
		op.CompositeMode = ebiten.CompositeModeDestinationIn
		mixerTemp.DrawImage(img, op)
		t := graph.TexFromImage(mixerTemp, ebiten.FilterDefault, ss.w, ss.h, 0, "~slidingthing")
		ss.Sprite.SetTex(t)
		ss.Draw(dest)
	}
	return f, ""
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
