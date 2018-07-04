package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

type cosmoPanels struct {
	left, right *ButtonsPanel
}

func newCosmoPanels() *cosmoPanels {
	panelOpts := ButtonsPanelOpts{
		PivotP:       graph.ScrP(0, 1),
		PivotV:       graph.BotLeft(),
		BorderSpace:  20 * graph.GS(),
		ButtonSpace:  10 * graph.GS(),
		ButtonLayer:  graph.Z_STAT_HUD + 100,
		ButtonSize:   v2.V2{X: 150, Y: 50}.Mul(graph.GS()),
		CaptionLayer: graph.Z_STAT_HUD + 101,
		SlideT:       0.5,
		SlideV:       v2.V2{X: 0.8, Y: 0},
	}
	leftPanel := NewButtonsPanel(panelOpts)
	leftPanel.SetActive(true)

	panelOpts = ButtonsPanelOpts{
		PivotP:       graph.ScrP(1, 1),
		PivotV:       graph.BotRight(),
		BorderSpace:  20 * graph.GS(),
		ButtonSpace:  10 * graph.GS(),
		ButtonLayer:  graph.Z_STAT_HUD + 100,
		ButtonSize:   v2.V2{X: 150, Y: 50}.Mul(graph.GS()),
		CaptionLayer: graph.Z_STAT_HUD + 101,
		SlideT:       0.5,
		SlideV:       v2.V2{X: -1, Y: 0},
	}

	rightPanel := NewButtonsPanel(panelOpts)

	return &cosmoPanels{
		left:  leftPanel,
		right: rightPanel,
	}
}

func (p *cosmoPanels) recalcLeft() {
	tex := GetAtlasTex(commons.ButtonAN)
	bo := ButtonOpts{
		Tex:    tex,
		Face:   Fonts[Face_mono],
		CapClr: color.White,
		Clr:    color.White,
	}
	p.left.ClearButtons()
	if Data.NaviData.BeaconCount > 0 {
		bo.Caption = fmt.Sprintf("МАЯК [%v]", Data.NaviData.BeaconCount)
		bo.Tags = "button_beacon"
		p.left.AddButton(bo)
	}
	if len(Data.NaviData.Mines) > 0 {
		bo.Caption = fmt.Sprintf("ШАХТА [%v]", Data.NaviData.BeaconCount)
		bo.Tags = "button_mine"
		p.left.AddButton(bo)
	}
	if len(Data.NaviData.Landing) > 0 {
		bo.Caption = "ВЫСАДКА"
		bo.Tags = "button_landing"
		p.left.AddButton(bo)
	}
}

func (p *cosmoPanels) update(dt float64) {
	p.left.Update(dt)
	p.right.Update(dt)
}

func (p *cosmoPanels) activeLeft(active bool) {
	p.left.SetActive(active)
}
func (p *cosmoPanels) activeRight(active bool) {
	p.right.SetActive(active)
}

func (p *cosmoPanels) Req() *graph.DrawQueue {
	Q := graph.NewDrawQueue()
	Q.Append(p.left)
	Q.Append(p.right)
	return Q
}
