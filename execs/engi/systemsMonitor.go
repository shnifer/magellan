package main

import (
	"github.com/Shnifer/magellan/draw"
	. "github.com/Shnifer/magellan/graph"
	"math"
)

type systemsMonitor struct {
	params CamParams

	all *Sprite
}

func newSystemsMonitor() *systemsMonitor {

	cam := NewCamera()
	cam.DenyGlobalScale = true
	cam.Center = ScrP(0.5, 0.5)

	scale1 := CalcGlobalScale(int(float64(WinH) / float64(DEFVAL.SpriteSizeH) * 1000))
	scale2 := CalcGlobalScale(int(float64(WinW) / float64(DEFVAL.SpriteSizeW) * 1000))
	scale := math.Min(scale1, scale2)

	cam.Scale = scale
	cam.Recalc()
	param := cam.Phys()

	all := draw.NewAtlasSprite("engi_all", param)

	res := systemsMonitor{
		params: param,
		all:    all,
	}

	return &res
}

func (s *systemsMonitor) Req(Q *DrawQueue) {
	Q.Add(s.all, Z_GAME_OBJECT)
}
