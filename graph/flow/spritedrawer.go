package flow

import (
	"github.com/Shnifer/magellan/graph"
	"math/rand"
)

type spriteDrawer struct {
	cs    *graph.CycledSprite
	layer int
}

type SpriteDrawerParams struct {
	Sprite       *graph.Sprite
	CycleType    int
	DoRandomLine bool
	FPS          float64
	Layer        int
}

func (params SpriteDrawerParams) New() updDrawPointer {
	var cs *graph.CycledSprite

	if params.DoRandomLine {
		rr := rand.Intn(params.Sprite.Rows())
		c := params.Sprite.Cols()
		min := rr * c
		max := min + c - 1
		limit := params.Sprite.SpritesCount() - 1
		if max > limit {
			max = limit
		}
		cs = graph.NewCycledSpriteRange(params.Sprite, params.CycleType, params.FPS, min, max)
	} else {
		cs = graph.NewCycledSprite(params.Sprite, params.CycleType, params.FPS)
	}
	return &spriteDrawer{
		cs:    cs,
		layer: params.Layer,
	}
}

func (sd *spriteDrawer) update(dt float64) {
	sd.cs.Update(dt)
}

func (sd *spriteDrawer) drawPoint(p Point) *graph.DrawQueue {
	sd.cs.SetPos(p.pos)

	if ang, ok := p.attr["Ang"]; ok {
		sd.cs.SetAng(ang)
	}
	if size, ok := p.attr["Size"]; ok {
		sd.cs.SetSizeProportion(size)
	}
	if alpha, ok := p.attr["Alpha"]; ok {
		sd.cs.SetAlpha(alpha)
	}

	res := graph.NewDrawQueue()
	res.Add(sd.cs, sd.layer)
	return res
}
