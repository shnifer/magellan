package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

type cosmoPanels struct {
	leftB, leftM, leftL int
	left, top, right    *ButtonsPanel
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

	panelOpts = ButtonsPanelOpts{
		PivotP:       graph.ScrP(0, 0),
		PivotV:       graph.TopLeft(),
		BorderSpace:  20 * graph.GS(),
		ButtonSpace:  10 * graph.GS(),
		ButtonLayer:  graph.Z_STAT_HUD + 100,
		ButtonSize:   v2.V2{X: 150, Y: 50}.Mul(graph.GS()),
		CaptionLayer: graph.Z_STAT_HUD + 101,
		SlideT:       0.5,
		SlideV:       v2.V2{X: 0, Y: 0.7},
	}

	topPanel := NewButtonsPanel(panelOpts)

	return &cosmoPanels{
		left:  leftPanel,
		right: rightPanel,
		top:   topPanel,
	}
}

func (p *cosmoPanels) recalcLeft() {
	if p.leftL == len(Data.NaviData.Landing) &&
		p.leftM == len(Data.NaviData.Mines) {
		return
	}
	p.leftM = len(Data.NaviData.Mines)
	p.leftL = len(Data.NaviData.Landing)

	p.left.ClearButtons()
	tex := GetAtlasTex(commons.ButtonAN)
	bo := ButtonOpts{
		Tex:    tex,
		Face:   Fonts[Face_mono],
		CapClr: color.White,
		Clr:    color.White,
	}
	bo.Caption = "SCAN"
	bo.Tags = "button_scan"
	p.left.AddButton(bo)

	if len(Data.NaviData.Mines) > 0 {
		bo.Caption = fmt.Sprintf("MINE [%v]", p.leftM)
		bo.Tags = "button_mine"
		p.left.AddButton(bo)
	}
	if len(Data.NaviData.Landing) > 0 {
		bo.Caption = "LANDING"
		bo.Tags = "button_landing"
		p.left.AddButton(bo)
	}
}

func (p *cosmoPanels) rightMines() {
	p.right.ClearButtons()
	tex := GetAtlasTex(commons.ButtonAN)

	mines := make(map[string]int)
	for _, corp := range Data.NaviData.Mines {
		mines[corp]++
	}
	for _, corp := range commons.CorpNames {
		if mines[corp] == 0 {
			continue
		}
		bo := ButtonOpts{
			Tex:     tex,
			Face:    Fonts[Face_mono],
			CapClr:  color.White,
			Clr:     commons.ColorByOwner(corp),
			Caption: commons.CompanyNameByOwner(corp),
			Tags:    minecorptagprefix + corp,
		}
		p.right.AddButton(bo)
	}
}

func (p *cosmoPanels) rightLanding() {
	p.right.ClearButtons()
	tex := GetAtlasTex(commons.ButtonAN)

	bo := ButtonOpts{
		Tex:     tex,
		Face:    Fonts[Face_mono],
		CapClr:  color.White,
		Clr:     color.White,
		Caption: "DO ORBIT",
		Tags:    "button_orbit",
	}
	p.right.AddButton(bo)

	bo = ButtonOpts{
		Tex:     tex,
		Face:    Fonts[Face_mono],
		CapClr:  color.White,
		Clr:     color.White,
		Caption: "LEAVE ORBIT",
		Tags:    "button_leaveorbit",
	}
	p.right.AddButton(bo)

}

func (p *cosmoPanels) recalcTop() {
	if p.leftB == Data.NaviData.BeaconCount {
		return
	}
	p.leftB = Data.NaviData.BeaconCount

	tex := GetAtlasTex(commons.ButtonAN)
	bo := ButtonOpts{
		Tex:          tex,
		Face:         Fonts[Face_mono],
		CapClr:       color.White,
		Clr:          color.White,
		Caption:      fmt.Sprintf("BEACON [%v]", p.leftB),
		Tags:         "button_beacon",
		HighlightClr: colornames.Green,
	}
	p.top.ClearButtons()
	p.top.AddButton(bo)
	p.top.SetActive(p.leftB > 0)
}

func (p *cosmoPanels) update(dt float64) {
	p.recalcLeft()
	p.recalcTop()

	p.left.Update(dt)
	p.right.Update(dt)
	p.top.Update(dt)
}

func (p *cosmoPanels) activeLeft(active bool) {
	p.left.SetActive(active)
}
func (p *cosmoPanels) activeRight(active bool) {
	p.right.SetActive(active)
}

func (p *cosmoPanels) Req(Q *graph.DrawQueue) {
	Q.Append(p.left)
	Q.Append(p.right)
	Q.Append(p.top)
}
