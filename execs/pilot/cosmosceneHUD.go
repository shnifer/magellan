package main

import (
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
)

type cosmoSceneHUD struct {
	background *graph.Sprite
	compass    *graph.Sprite

	arrowSize float64
	rulerSize float64
	//hud
	thrustLevel   *graph.Sprite
	thrustControl *graph.Sprite
	turnLevel     *graph.Sprite
	turnControl   *graph.Sprite
	rulerV        *graph.Sprite
	rulerH        *graph.Sprite
	//
	f9 *graph.Frame9HUD
}

const (
	rulerX     = 0.08
	rulerY     = 0.16
	arrSize    = 0.05
	wide       = 0.6
	rulerWideK = 1.1

	compassSize = 0.55
)

func newCosmoSceneHUD(cam *graph.Camera) cosmoSceneHUD {
	var background *graph.Sprite

	if !DEFVAL.LowQ {
		background = NewAtlasSpriteHUD(commons.DefaultBackgroundAN)
		background.SetSize(float64(WinW), float64(WinH))
		background.SetPivot(graph.TopLeft())
		background.SetColor(colornames.Dimgrey)
	}

	compass := NewAtlasSprite(commons.CompassAN, cam.FixS())
	compassSize := float64(WinH) * compassSize
	compass.SetSize(compassSize, compassSize)
	compass.SetAlpha(1)

	f9 := NewAtlasFrame9HUD(commons.Frame9AN, WinW, WinH, graph.Z_HUD-1)

	arrowSize := float64(WinH) * arrSize
	arrowTex := GetAtlasTex(commons.ThrustArrowAN)
	thrustLevel := graph.NewSpriteHUD(arrowTex)
	thrustLevel.SetSizeProportion(arrowSize)
	thrustLevel.SetAng(-90)
	thrustControl := graph.NewSpriteHUD(arrowTex)
	thrustControl.SetSizeProportion(arrowSize)
	thrustControl.SetAng(90)
	turnLevel := graph.NewSpriteHUD(arrowTex)
	turnLevel.SetSizeProportion(arrowSize)
	turnLevel.SetAng(0)
	turnControl := graph.NewSpriteHUD(arrowTex)
	turnControl.SetSizeProportion(arrowSize)
	turnControl.SetAng(180)

	rulerSize := float64(WinH) * wide
	rulerH := NewAtlasSpriteHUD(commons.RulerHAN)
	rulerH.SetSizeProportion(rulerSize * rulerWideK)
	rulerH.SetPos(graph.ScrP(0.5, rulerY))
	rulerV := NewAtlasSpriteHUD(commons.RulerVAN)
	rulerV.SetSizeProportion(rulerSize * rulerWideK)
	rulerV.SetPos(graph.ScrP(rulerX, 0.5))

	return cosmoSceneHUD{
		background:    background,
		compass:       compass,
		f9:            f9,
		thrustControl: thrustControl,
		thrustLevel:   thrustLevel,
		turnControl:   turnControl,
		turnLevel:     turnLevel,
		rulerV:        rulerV,
		rulerH:        rulerH,
		arrowSize:     arrowSize,
		rulerSize:     rulerSize,
	}
}

func (s *cosmoScene) UpdateHUD() {
	var p v2.V2
	arrS := s.hud.arrowSize * 0.6

	ruler := func(x float64) float64 {
		return 0.5 - s.hud.rulerSize/2*x
	}

	vPos := graph.ScrP(rulerX, 0.5)
	p = vPos.Add(v2.V2{X: arrS, Y: ruler(s.thrustLevel)})
	s.hud.thrustLevel.SetPos(p)

	p = vPos.Add(v2.V2{X: -arrS, Y: ruler(input.GetF("forward"))})
	s.hud.thrustControl.SetPos(p)

	hPos := graph.ScrP(0.5, rulerY)
	p = hPos.Add(v2.V2{X: ruler(s.maneurLevel), Y: arrS})
	s.hud.turnLevel.SetPos(p)

	p = hPos.Add(v2.V2{X: ruler(input.GetF("turn")), Y: -arrS})
	s.hud.turnControl.SetPos(p)

	s.hud.compass.SetPos(Data.PilotData.Ship.Pos)
}

func (h cosmoSceneHUD) Req(Q *graph.DrawQueue) {

	if !DEFVAL.LowQ {
		Q.Add(h.background, graph.Z_STAT_BACKGROUND)
	}
	Q.Add(h.compass, graph.Z_HUD)

	Q.Add(h.rulerV, graph.Z_HUD)
	Q.Add(h.rulerH, graph.Z_HUD)
	Q.Add(h.thrustLevel, graph.Z_HUD+1)
	Q.Add(h.thrustControl, graph.Z_HUD+1)
	Q.Add(h.turnLevel, graph.Z_HUD+1)
	Q.Add(h.turnControl, graph.Z_HUD+1)

	Q.Append(h.f9)
}
