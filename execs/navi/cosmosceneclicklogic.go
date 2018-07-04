package main

import (
	"github.com/Shnifer/magellan/v2"
	"log"
)

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
	log.Println(tag)
}
