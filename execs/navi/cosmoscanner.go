package main

import (
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"math"
)

const (
	scanZero = iota
	scanSelect
	scanProgress
	ScanDone
)

type scanner struct {
	state int
	work  string

	scanPart float64
	obj      *CosmoPoint

	maxRange  float64
	scanSpeed float64

	countSprite *graph.Sprite
	countN      int

	scanRange  *graph.CircleLine
	dropRange  *graph.CircleLine
	scanSector *graph.Sector

	onState func(state int)
}

func newScanner(cam *graph.Camera, onState func(int)) *scanner {
	defer LogFunc("newScanner")()

	const countN = 12

	sprite := graph.NewSprite(GetAtlasTex(commons.ScannerAN), cam.Deny())
	sprite.SetSize(15, 15)

	opts := graph.CircleLineOpts{
		Params: cam.Phys(),
		Layer:  graph.Z_UNDER_OBJECT,
		PCount: 32,
		Clr:    colornames.Indigo,
	}
	scanRange := graph.NewCircleLine(v2.ZV, 0, opts)
	dropRange := graph.NewCircleLine(v2.ZV, 0, opts)

	scanSector := graph.NewSector(cam.Phys())
	scanSector.SetColor(colornames.Goldenrod)

	return &scanner{
		countSprite: sprite,
		countN:      countN,
		scanRange:   scanRange,
		dropRange:   dropRange,
		scanSector:  scanSector,
		onState:     onState,
	}
}

func (s *scanner) clicked(obj *CosmoPoint) bool {
	/*	if Data.PilotData.Ship.Pos.Sub(obj.Pos).LenSqr() > s.maxRange2 {
			return false
		}
	*/if s.obj == obj {
		return false
	}
	s.stateSelect(obj)
	return true
}

func (s *scanner) update(ship v2.V2, dt float64) {
	s.scanRange.SetPos(ship)
	s.dropRange.SetPos(ship)

	if s.state == scanZero {
		return
	}

	//graphics
	v := s.obj.Pos.Sub(ship)
	dist := v.Len() + 0.1
	ang := v.Dir()
	angW := math.Atan(s.obj.Size/dist) * v2.Rad2Deg * 0.8
	s.scanSector.SetCenterRadius(ship, dist)
	s.scanSector.SetAngles(ang-angW, ang+angW)

	if s.state == scanSelect {
		return
	}

	if s.work == button_scan {
		s.maxRange = Data.SP.Scanner.ScanRange
		s.scanSpeed = Data.SP.Scanner.ScanSpeed
	} else {
		s.maxRange = Data.SP.Scanner.DropRange
		s.scanSpeed = Data.SP.Scanner.DropSpeed
	}

	if ship.Sub(s.obj.Pos).LenSqr() > s.maxRange*s.maxRange {
		s.stateZero()
		return
	}

	if s.state == scanProgress {
		s.scanPart += dt * s.scanSpeed
		if s.scanPart > 1 {
			s.stateDone()
		}
	}
}

func (s *scanner) Start(work string) {
	if s.state == scanZero || s.obj == nil {
		return
	}
	if s.state == scanProgress && s.work == work {
		return
	}

	s.stateProgress(work)
}

func (s *scanner) Req(Q *graph.DrawQueue) {
	switch s.work {
	case button_scan:
		s.scanRange.SetColor(colornames.Mediumpurple)
		s.dropRange.SetColor(colornames.Dimgray)
	case button_mine, button_landing:
		s.scanRange.SetColor(colornames.Dimgray)
		s.dropRange.SetColor(colornames.Orangered)
	default:
		s.scanRange.SetColor(colornames.Indigo)
		s.dropRange.SetColor(colornames.Indianred)
	}
	s.scanRange.SetRadius(Data.SP.Scanner.ScanRange)
	s.dropRange.SetRadius(Data.SP.Scanner.DropRange)
	Q.Append(s.scanRange)
	Q.Append(s.dropRange)

	if s.state == scanZero {
		return
	}

	Q.Add(s.scanSector, graph.Z_UNDER_OBJECT)

	if s.state != scanSelect {
		//Draw circle counter
		num := int(0.5 + s.scanPart*float64(s.countN))
		obj := s.obj.Pos
		rng := s.obj.Size * 1.1
		if rng < 30 {
			rng = 30
		}
		for i := 0; i < num; i++ {
			pos := obj.AddMul(v2.InDir(-360/float64(s.countN)*float64(i)), rng)
			s.countSprite.SetPos(pos)
			Q.Add(s.countSprite, graph.Z_UNDER_OBJECT)
		}
	}
}

func (s *scanner) procOnState() {
	if s.onState != nil {
		s.onState(s.state)
	}
}

func (s *scanner) stateZero() {
	s.state = scanZero
	s.work = ""
	s.scanPart = 0
	s.obj = nil
	s.procOnState()
}

func (s *scanner) stateSelect(obj *CosmoPoint) {
	s.state = scanSelect
	s.work = ""
	s.scanPart = 0
	s.obj = obj
	s.procOnState()
}

func (s *scanner) stateProgress(work string) {
	s.work = work
	s.scanPart = 0
	s.state = scanProgress
	s.procOnState()
}

func (s *scanner) stateDone() {
	s.state = ScanDone
	s.procOnState()
}
