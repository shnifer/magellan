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
	back    *graph.Sprite

	//trail
	trailT float64
	trail  *graph.FadingArray

	thrustLevel   *graph.Sprite
	thrustControl *graph.Sprite
	turnLevel     *graph.Sprite
	turnControl   *graph.Sprite
	rulerV        *graph.Sprite
	rulerH        *graph.Sprite
	arrowSize     float64
	rulerSize     float64

	compass    *graph.Sprite
	distCircle *graph.CircleLine
	physRadText *graph.Text
}

func newWarpSceneHUD(cam *graph.Camera) warpSceneHUD {
	caption := graph.NewText("Warp scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	clo := graph.CircleLineOpts{
		Layer:  graph.Z_STAT_HUD + 10,
		Clr:    colornames.Oldlace,
		PCount: 64,
	}
	distCircle := graph.NewCircleLine(cam.Center, float64(WinH)*0.3, clo)

	var background *graph.Sprite
	if !DEFVAL.LowQ {
		background = NewAtlasSpriteHUD(WarpBackgroundAN)
		background.SetSize(float64(WinW), float64(WinH))
		background.SetPivot(graph.TopLeft())
		background.SetColor(colornames.Dimgrey)
	}

	arrowSize := float64(WinH) * arrSize
	arrowTex := GetAtlasTex(ThrustArrowAN)
	thrustLevel := graph.NewSpriteHUD(arrowTex)
	thrustLevel.SetSize(50, 50)
	thrustLevel.SetAng(-90)
	thrustLevel.SetAlpha(0.7)
	thrustControl := graph.NewSpriteHUD(arrowTex)
	thrustControl.SetSize(50, 50)
	thrustControl.SetAng(90)
	thrustControl.SetAlpha(0.5)
	turnLevel := graph.NewSpriteHUD(arrowTex)
	turnLevel.SetSize(50, 50)
	turnLevel.SetAng(0)
	turnLevel.SetAlpha(0.7)
	turnControl := graph.NewSpriteHUD(arrowTex)
	turnControl.SetSize(50, 50)
	turnControl.SetAng(180)
	turnControl.SetAlpha(0.5)

	compass := NewAtlasSprite(CompassAN, cam.FixS())
	compassSize := float64(WinH) * compassSize
	compass.SetSize(compassSize, compassSize)
	compass.SetAlpha(1)

	trail := graph.NewFadingArray(GetAtlasTex(TrailAN), trailLifeTime/trailPeriod,
		cam.Deny())

	rulerSize := float64(WinH) * wide
	rulerH := NewAtlasSpriteHUD(RulerWarpHAN)
	rulerH.SetSizeProportion(rulerSize * rulerWideK)
	rulerH.SetPos(graph.ScrP(0.5, rulerY))
	rulerV := NewAtlasSpriteHUD(RulerWarpVAN)
	rulerV.SetSizeProportion(rulerSize * rulerWideK)
	rulerV.SetPos(graph.ScrP(rulerX, 0.5))

	return warpSceneHUD{
		back:          background,
		trail:         trail,
		compass:       compass,
		caption:       caption,
		distCircle:    distCircle,
		thrustLevel:   thrustLevel,
		turnLevel:     turnLevel,
		turnControl:   turnControl,
		thrustControl: thrustControl,
		arrowSize:     arrowSize,
		rulerSize:     rulerSize,
		rulerV:        rulerV,
		rulerH:        rulerH,
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

	var p v2.V2
	arrS := s.hud.arrowSize * 0.6

	ruler := func(x float64) float64 {
		return -s.hud.rulerSize / 2 * x
	}
	vRuler := func(x float64) float64 {
		return s.hud.rulerSize/2 - s.hud.rulerSize*x
	}
	vPos := graph.ScrP(rulerX, 0.5)
	p = vPos.Add(v2.V2{X: arrS, Y: vRuler(s.thrustLevel)})
	s.hud.thrustLevel.SetPos(p)

	p = vPos.Add(v2.V2{X: -arrS, Y: vRuler(input.WarpLevel("warpspeed"))})
	s.hud.thrustControl.SetPos(p)

	hPos := graph.ScrP(0.5, rulerY)
	p = hPos.Add(v2.V2{X: ruler(s.maneurLevel), Y: arrS})
	s.hud.turnLevel.SetPos(p)

	turnInput := input.GetF("turn")
	p = hPos.Add(v2.V2{X: ruler(turnInput), Y: -arrS})
	s.hud.turnControl.SetPos(p)

	s.hud.compass.SetPos(Data.PilotData.Ship.Pos)

	circleRadPx := float64(WinH) * 0.3
	physRad := circleRadPx / s.cam.Scale / graph.GS()

	msg := fmt.Sprintf("circle radius: %f", physRad)
	physRadText := graph.NewText(msg, Fonts[Face_mono], colornames.Oldlace)
	physRadText.SetPosPivot(graph.ScrP(0.5, 0.4), graph.TopMiddle())
	s.hud.physRadText = physRadText
}

func (h warpSceneHUD) Req(Q *graph.DrawQueue) {
	Q.Add(h.back, graph.Z_STAT_BACKGROUND)
	Q.Add(h.trail, graph.Z_UNDER_OBJECT)
	Q.Add(h.thrustLevel, graph.Z_HUD)
	Q.Add(h.thrustControl, graph.Z_HUD)
	Q.Add(h.turnLevel, graph.Z_HUD)
	Q.Add(h.turnControl, graph.Z_HUD)
	Q.Add(h.compass, graph.Z_HUD)
	Q.Add(h.rulerV, graph.Z_HUD)
	Q.Add(h.rulerH, graph.Z_HUD)

	//Q.Add(h.caption, graph.Z_STAT_HUD)
	Q.Append(h.distCircle)
	Q.Add(h.physRadText, graph.Z_STAT_HUD+10)
}

func (s *warpScene) drawScale(Q *graph.DrawQueue) {
	//PosScale factor hud
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
