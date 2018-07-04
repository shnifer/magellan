package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func (s *cosmoScene) updateControl(dt float64) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mousex, mousey := ebiten.CursorPosition()
		s.procMouseClick(mousex, mousey)
	}

	s.scanner.update(s.shipRB.RB().Pos, dt)
	s.cosmoPanels.update(dt)
}

func (s *cosmoScene) procMouseClick(x, y int) {
	if s.cosmoPanels != nil {
		if tag, ok := s.cosmoPanels.left.ProcMouse(x, y); ok {
			s.procButtonClick(tag)
			return
		}
		if tag, ok := s.cosmoPanels.right.ProcMouse(x, y); ok {
			s.procButtonClick(tag)
			return
		}
		if tag, ok := s.cosmoPanels.top.ProcMouse(x, y); ok {
			s.procButtonClick(tag)
			return
		}
	}

	worldPos := s.cam.UnApply(v2.V2{X: float64(x), Y: float64(y)})
	for id, obj := range Data.Galaxy.Points {
		if worldPos.Sub(obj.Pos).LenSqr() < (obj.Size * obj.Size) {
			s.scanner.clicked(s.objects[id])
			return
		}
	}

	Data.NaviData.ActiveMarker = true
	Data.NaviData.MarkerPos = worldPos
	s.naviMarkerT = DEFVAL.NaviMarketDuration
}

func (s *cosmoScene) procButtonClick(tag string) {
	switch tag {
	case "button_mine":
	case "button_landing":
	case "button_beacon":
		//todo:enter text
		AddBeacon(Data, Client, "just a test beacon")
		ClientLogGame(Client, "ADD BEACON KEY", "just a test beacon")
		Data.NaviData.BeaconCount--
	default:
		Log(LVL_ERROR, "Unknown button tag ", tag)
	}
}

func (s *cosmoScene) scanState(scanState int) {
	switch scanState {
	case scanZero:
		s.cosmoPanels.activeLeft(false)
		s.cosmoPanels.activeRight(false)
	case scanSelect:
		s.cosmoPanels.activeLeft(true)
		s.cosmoPanels.activeRight(false)
	case scanProgress:
	case ScanDone:
	default:
		Log(LVL_ERROR, "Unknown scan state ", scanState)
	}
}
