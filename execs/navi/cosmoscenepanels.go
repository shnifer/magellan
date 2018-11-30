package main

import (
	"fmt"
	"github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/draw"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
)

type cosmoPanels struct {
	leftB, leftM, leftL    int
	left, top, right, back *ButtonsPanel

	showSignature bool
	sigs          []commons.Signature
	sonar         *SonarHUD
	sonarBack     *graph.Sprite
	sonarPos      v2.V2
	sonarSize     float64
	sonarName     string
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

	panelOpts = ButtonsPanelOpts{
		PivotP:       graph.ScrP(1, 0),
		PivotV:       graph.TopRight(),
		BorderSpace:  20 * graph.GS(),
		ButtonSpace:  10 * graph.GS(),
		ButtonLayer:  graph.Z_STAT_HUD + 100,
		ButtonSize:   v2.V2{X: 150, Y: 50}.Mul(graph.GS()),
		CaptionLayer: graph.Z_STAT_HUD + 101,
		SlideT:       0.5,
		SlideV:       v2.V2{X: 0, Y: 1},
	}

	tex := GetAtlasTex(commons.ButtonAN)
	backPanel := NewButtonsPanel(panelOpts)
	bo := ButtonOpts{
		Tex:     tex,
		Face:    Fonts[Face_mono],
		CapClr:  color.White,
		Clr:     color.White,
		Caption: "Домой!",
		Tags:    button_home,
	}
	backPanel.AddButton(bo)

	sonarSize := float64(WinH / 3)
	offset := float64(WinH / 10)
	sonarPos := graph.ScrP(1, 0).AddMul(v2.V2{X: -1, Y: 1}, sonarSize/2+offset)
	sonar := NewSonarHUD(sonarPos, sonarSize, graph.NoCam, graph.Z_HUD)

	sonarBack := graph.NewSprite(graph.CircleTex(), graph.NoCam)
	sonarBack.SetPos(sonarPos)
	sonarBack.SetSize(sonarSize, sonarSize)
	sonarBack.SetColor(color.Black)
	sonarBack.SetAlpha(0.8)

	return &cosmoPanels{
		left:  leftPanel,
		right: rightPanel,
		top:   topPanel,
		back:  backPanel,

		sonar:     sonar,
		sonarPos:  sonarPos,
		sonarSize: sonarSize,
		sonarBack: sonarBack,
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
	bo.Tags = button_scan
	p.left.AddButton(bo)

	if len(Data.NaviData.Mines) > 0 {
		bo.Caption = fmt.Sprintf("MINE [%v]", p.leftM)
		bo.Tags = button_mine
		p.left.AddButton(bo)
	}
	if len(Data.NaviData.Landing) > 0 {
		bo.Caption = "LANDING"
		bo.Tags = button_landing
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
		Tags:    button_orbit,
	}
	p.right.AddButton(bo)

	bo = ButtonOpts{
		Tex:     tex,
		Face:    Fonts[Face_mono],
		CapClr:  color.White,
		Clr:     color.White,
		Caption: "LEAVE ORBIT",
		Tags:    button_leaveorbit,
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
		Tags:         button_beacon,
		HighlightClr: colornames.Green,
	}
	p.top.ClearButtons()
	p.top.AddButton(bo)
	p.top.SetActive(p.leftB > 0)
}

func (p *cosmoPanels) update(dt float64) {
	p.recalcLeft()
	p.recalcTop()
	canHome := Data.NaviData.CanLandhome && Data.GalaxyID == commons.START_Galaxy_ID
	p.back.SetActive(canHome)

	p.left.Update(dt)
	p.right.Update(dt)
	p.top.Update(dt)
	p.back.Update(dt)

	p.sonar.ActiveSignatures(p.sigs)
	p.sonar.Update(dt)
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
	Q.Append(p.back)

	if p.showSignature {
		Q.Add(p.sonarBack, graph.Z_HUD-1)
		Q.Append(p.sonar)
		T := graph.NewText(p.sonarName, Fonts[Face_cap], colornames.Orangered)
		T.SetPosPivot(p.sonarPos.Sub(graph.ScrP(0, 0.2)), graph.Center())
		Q.Add(T, graph.Z_STAT_HUD)
	}
}

func (p *cosmoPanels) ShowSonar(sigs []commons.Signature, name string) {
	p.showSignature = true
	p.sigs = sigs
	p.sonarName = name
}

func (p *cosmoPanels) CloseSonar() {
	p.showSignature = false
	p.sigs = []commons.Signature{}
	p.sonarName = ""
}
