package main

import (
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
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

	scanT float64
	obj   *CosmoPoint

	totalT    float64
	maxRange2 float64

	countSprite *graph.Sprite
	countN      int

	scanRange  *graph.Sprite
	scanSector *graph.Sector

	onState func(state int)
}

func newScanner(cam *graph.Camera, onState func(int)) *scanner {
	defer LogFunc("newScanner")()

	const countN = 12

	sprite := graph.NewSprite(GetAtlasTex(commons.ScannerAN), cam.Deny())
	sprite.SetSize(15, 15)

	scanRange := graph.NewSprite(graph.CircleTex(), cam.Phys())
	scanRange.SetAlpha(0.5)
	scanRange.SetColor(colornames.Indigo)

	scanSector := graph.NewSector(cam.Phys())
	scanSector.SetColor(colornames.Goldenrod)

	return &scanner{
		countSprite: sprite,
		countN:      countN,
		scanRange:   scanRange,
		scanSector:  scanSector,
		onState:     onState,
	}
}

func (s *scanner) clicked(obj *CosmoPoint) {
	if Data.PilotData.Ship.Pos.Sub(obj.Pos).LenSqr() > s.maxRange2 {
		return
	}
	if s.obj == obj {
		return
	}
	s.stateSelect(obj)
}

func (s *scanner) update(ship v2.V2, dt float64) {
	s.scanRange.SetPos(ship)
	s.maxRange2 = Data.SP.Radar.Scan_range * Data.SP.Radar.Scan_range

	if s.state == scanZero {
		return
	}

	if Data.SP.Radar.Scan_speed > 0 {
		s.totalT = 1 / Data.SP.Radar.Scan_speed
	} else {
		s.totalT = 0
	}

	v := s.obj.Pos.Sub(ship)
	dist := v.Len() + 0.1
	ang := v.Dir()
	angW := math.Atan(s.obj.Size/dist) * v2.Rad2Deg * 0.8
	s.scanSector.SetCenterRadius(ship, dist)
	s.scanSector.SetAngles(ang-angW, ang+angW)

	if ship.Sub(s.obj.Pos).LenSqr() > s.maxRange2 {
		s.stateZero()
		return
	}

	if s.state == scanProgress {
		if s.totalT > 0 {
			s.scanT += dt
			if s.scanT > s.totalT {
				s.stateDone()
			}
		} else {
			s.stateZero()
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
	Range := Data.SP.Radar.Scan_range * 2
	s.scanRange.SetSize(Range, Range)
	Q.Add(s.scanRange, graph.Z_UNDER_OBJECT)

	if s.state == scanZero {
		return
	}

	Q.Add(s.scanSector, graph.Z_UNDER_OBJECT)

	if s.state != scanSelect {
		//Draw circle counter
		num := int(0.5 + s.scanT/s.totalT*float64(s.countN))
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

func (s *scanner) procScanned(obj *CosmoPoint) {
	commons.ClientLogGame(Client, "scan", "SCANNED ", obj.ID)
	gp, ok := Data.Galaxy.Points[obj.ID]
	if !ok {
		return
	}
	if gp.ScanData == "" {
		return
	}
	if gp.Type == commons.BUILDING_BEACON {
		commons.RequestRemoveBuilding(Client, obj.ID)
	}
	corp := commons.CorpNames[rand.Intn(5)]
	if rand.Intn(2) == 0 {
		//proc random mine
		mineFK, exist := gp.Mines[corp]
		if exist {
			commons.RequestRemoveBuilding(Client, mineFK)
		} else {
			commons.AddMine(Data, Client, obj.ID, corp)
		}
	} else {
		//proc random fishhouse
		fishhouseFK, exist := gp.FishHouses[corp]
		if exist {
			commons.RequestRemoveBuilding(Client, fishhouseFK)
		} else {
			commons.AddFishHouse(Data, Client, obj.ID, corp)
		}
	}

	return
}

func (s *scanner) procOnState() {
	if s.onState != nil {
		s.onState(s.state)
	}
}

func (s *scanner) stateZero() {
	s.state = scanZero
	s.work = ""
	s.scanT = 0
	s.obj = nil
	s.procOnState()
}

func (s *scanner) stateSelect(obj *CosmoPoint) {
	s.state = scanSelect
	s.work = ""
	s.scanT = 0
	s.obj = obj
	s.procOnState()
}

func (s *scanner) stateProgress(work string) {
	s.work = work
	s.scanT = 0
	s.state = scanProgress
	s.procOnState()
}

func (s *scanner) stateDone() {
	s.state = ScanDone
	s.procOnState()
}
