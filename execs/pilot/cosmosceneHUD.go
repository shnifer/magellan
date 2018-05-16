package main

import (
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
	rulerX     = 0.1
	rulerY     = 0.1
	arrSize    = 0.05
	wide       = 0.6
	rulerWideK = 1.1

	compassSize = 0.7
)

func newCosmoSceneHUD(cam *graph.Camera) cosmoSceneHUD {
	background := NewAtlasSpriteHUD("background")
	background.SetSize(float64(WinW), float64(WinH))
	background.SetPivot(graph.TopLeft())
	background.SetColor(colornames.Dimgrey)

	compass := NewAtlasSprite("compass", cam, true, false)
	compassSize := float64(WinH) * compassSize
	compass.SetSize(compassSize, compassSize)
	compass.SetAlpha(1)

	f9 := NewAtlasFrame9HUD("front9", WinW, WinH)
	f9.SetScale(0.5)

	arrowSize := float64(WinH) * arrSize
	arrowTex := GetAtlasTex("arrow")
	thrustLevel := graph.NewSpriteHUD(arrowTex)
	thrustLevel.SetSizeProportion(arrowSize)
	thrustLevel.SetAng(-90)
	thrustLevel.SetAlpha(0.8)
	thrustControl := graph.NewSpriteHUD(arrowTex)
	thrustControl.SetSizeProportion(arrowSize)
	thrustControl.SetAng(90)
	thrustControl.SetAlpha(0.8)
	turnLevel := graph.NewSpriteHUD(arrowTex)
	turnLevel.SetSizeProportion(arrowSize)
	turnLevel.SetAng(0)
	turnLevel.SetAlpha(0.8)
	turnControl := graph.NewSpriteHUD(arrowTex)
	turnControl.SetSizeProportion(arrowSize)
	turnControl.SetAng(180)
	turnControl.SetAlpha(0.8)

	rulerSize := graph.ScrP(wide*rulerWideK, wide*rulerWideK)
	rulerH := NewAtlasSpriteHUD("rulerH")
	rulerH.SetSizeProportion(rulerSize.X)
	rulerH.SetPos(graph.ScrP(0.5, rulerY))
	rulerV := NewAtlasSpriteHUD("rulerV")
	rulerV.SetSizeProportion(rulerSize.Y)
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
	}
}

func (s *cosmoScene) UpdateHUD() {
	var p v2.V2
	arrS := s.hud.arrowSize
	p = graph.ScrP(rulerX, 0.5-wide/2*s.thrustLevel).AddMul(v2.InDir(-90), arrS)
	s.hud.thrustLevel.SetPos(p)

	p = graph.ScrP(rulerX, 0.5-wide/2*input.GetF("forward")).AddMul(v2.InDir(90), arrS)
	s.hud.thrustControl.SetPos(p)

	p = graph.ScrP(0.5-0.3*s.maneurLevel, rulerY).AddMul(v2.InDir(0), arrS)
	s.hud.turnLevel.SetPos(p)

	p = graph.ScrP(0.5-0.3*input.GetF("turn"), rulerY).AddMul(v2.InDir(180), arrS)
	s.hud.turnControl.SetPos(p)

	s.hud.compass.SetPos(Data.PilotData.Ship.Pos)
}

func (h cosmoSceneHUD) Req() *graph.DrawQueue {
	Q := graph.NewDrawQueue()
	Q.Add(h.background, graph.Z_STAT_BACKGROUND)
	Q.Add(h.compass, graph.Z_BACKGROUND)

	Q.Add(h.rulerV, graph.Z_HUD)
	Q.Add(h.rulerH, graph.Z_HUD)
	Q.Add(h.thrustLevel, graph.Z_HUD+1)
	Q.Add(h.thrustControl, graph.Z_HUD+1)
	Q.Add(h.turnLevel, graph.Z_HUD+1)
	Q.Add(h.turnControl, graph.Z_HUD+1)

	Q.Add(h.f9, graph.Z_HUD-1)

	return Q
}
