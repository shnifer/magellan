package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
)

type warpSceneHUD struct {
	caption *graph.Text

	//trail
	trailT float64
	trail  *graph.FadingArray

	thrustLevelHUD   *graph.Sprite
	thrustControlHUD *graph.Sprite
	turnLevelHUD     *graph.Sprite
	turnControlHUD   *graph.Sprite

	distCircle *graph.CircleLine
}

func newWarpSceneHUD(cam *graph.Camera) warpSceneHUD {
	caption := graph.NewText("Warp scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	clo := graph.CircleLineOpts{
		Layer:  graph.Z_STAT_HUD + 10,
		Clr:    colornames.Oldlace,
		PCount: 32,
	}
	distCircle := graph.NewCircleLine(cam.Center, float64(WinH)*0.3, clo)

	arrowTex := GetAtlasTex(ThrustArrowAN)
	thrustLevelHUD := graph.NewSpriteHUD(arrowTex)
	thrustLevelHUD.SetSize(50, 50)
	thrustLevelHUD.SetAng(-90)
	thrustLevelHUD.SetAlpha(0.7)
	thrustControlHUD := graph.NewSpriteHUD(arrowTex)
	thrustControlHUD.SetSize(50, 50)
	thrustControlHUD.SetAng(90)
	thrustControlHUD.SetAlpha(0.5)
	turnLevelHUD := graph.NewSpriteHUD(arrowTex)
	turnLevelHUD.SetSize(50, 50)
	turnLevelHUD.SetAng(0)
	turnLevelHUD.SetAlpha(0.7)
	turnControlHUD := graph.NewSpriteHUD(arrowTex)
	turnControlHUD.SetSize(50, 50)
	turnControlHUD.SetAng(180)
	turnControlHUD.SetAlpha(0.5)

	trail := graph.NewFadingArray(GetAtlasTex(TrailAN), trailLifeTime/trailPeriod,
		cam.Deny())

	return warpSceneHUD{
		trail:            trail,
		caption:          caption,
		distCircle:       distCircle,
		thrustLevelHUD:   thrustLevelHUD,
		turnLevelHUD:     turnLevelHUD,
		turnControlHUD:   turnControlHUD,
		thrustControlHUD: thrustControlHUD,
	}
}

func (s *warpScene) updateHUD() {
	s.hud.trailT += dt
	if s.hud.trailT > trailPeriod {
		s.hud.trailT -= trailPeriod

		s.hud.trail.Add(graph.ArrayElem{
			Size:     5,
			Pos:      Data.PilotData.Ship.Pos,
			LifeTime: trailLifeTime,
		})
	}
	s.hud.trail.Update(dt)

	s.hud.thrustLevelHUD.SetPos(graph.ScrP(0.15, 0.9-0.8*s.thrustLevel))
	s.hud.thrustControlHUD.SetPos(graph.ScrP(0.1, 0.9-0.8*input.WarpLevel("warpspeed")))

	s.hud.turnLevelHUD.SetPos(graph.ScrP(0.5-0.4*s.maneurLevel, 0.15))
	s.hud.turnControlHUD.SetPos(graph.ScrP(0.5-0.4*input.GetF("turn"), 0.1))
}

func (h warpSceneHUD) Req(Q *graph.DrawQueue) {
	Q.Add(h.trail, graph.Z_UNDER_OBJECT)
	Q.Add(h.thrustLevelHUD, graph.Z_HUD)
	Q.Add(h.thrustControlHUD, graph.Z_HUD)
	Q.Add(h.turnLevelHUD, graph.Z_HUD)
	Q.Add(h.turnControlHUD, graph.Z_HUD)

	Q.Add(h.caption, graph.Z_STAT_HUD)
	Q.Append(h.distCircle)
}

func (s *warpScene) drawScale(Q *graph.DrawQueue) {
	//Scale factor hud
	camScale := s.cam.Scale * graph.GS()
	maxLen := float64(WinW) * 0.8
	order := math.Floor(math.Log10(maxLen / camScale))
	val := math.Pow10(int(order))
	l := camScale * val

	from := graph.ScrP(0.1, 0.9)
	to := from.AddMul(v2.V2{X: 1, Y: 0}, l)
	mid := from.AddMul(v2.V2{X: 1, Y: 0}, l/2)
	mid.Y += 10

	tick := v2.V2{X: 0, Y: 5}

	graph.LineScr(Q, from, to, colornames.White, graph.Z_STAT_HUD+10)
	graph.LineScr(Q, from.Sub(tick), from.Add(tick), colornames.White, graph.Z_STAT_HUD+10)
	graph.LineScr(Q, to.Sub(tick), to.Add(tick), colornames.White, graph.Z_STAT_HUD+10)

	msg := fmt.Sprintf("%v", val)
	scaleText := graph.NewText(msg, Fonts[Face_mono], colornames.White)
	scaleText.SetPosPivot(mid, graph.TopMiddle())
	Q.Add(scaleText, graph.Z_STAT_HUD+10)

	circleRadPx := float64(WinH) * 0.3
	physRad := circleRadPx / s.cam.Scale / graph.GS()

	msg = fmt.Sprintf("circle radius: %f", physRad)
	physRadText := graph.NewText(msg, Fonts[Face_mono], colornames.Oldlace)
	physRadText.SetPosPivot(graph.ScrP(0.5, 0.4), graph.TopMiddle())
	Q.Add(physRadText, graph.Z_STAT_HUD+10)
}

func (s *warpScene) drawGravity(Q *graph.DrawQueue) {
	scale := float64(WinH) * 0.3 / (s.cam.Scale * graph.GS())
	ship := Data.PilotData.Ship.Pos
	thrust := Data.PilotData.ThrustVector
	drawv := func(v v2.V2, clr color.Color) {
		graph.Line(Q, s.cam, ship, ship.AddMul(v, scale), clr, graph.Z_STAT_HUD+10)
	}
	for _, v := range s.gravityReport {
		drawv(v, colornames.Deepskyblue)
	}
	drawv(s.gravityAcc, colornames.Lightblue)
	drawv(Data.PilotData.ThrustVector, colornames.Darkolivegreen)
	drawv(thrust.Add(s.gravityAcc), colornames.White)
}
