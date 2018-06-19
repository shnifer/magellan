package main

import (
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/graph/qr"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"math"
	"time"
)

type scanner struct {
	scanT float64
	obj   *CosmoPoint

	totalT    float64
	maxRange2 float64

	countSprite *graph.Sprite
	countN      int

	scanRange  *graph.Sprite
	scanSector *graph.Sector

	scannedImg *graph.Sprite
}

func newScanner(cam *graph.Camera) *scanner {
	defer LogFunc("newScanner")()

	const countN = 12

	sprite := graph.NewSprite(GetAtlasTex("trail"), cam.Deny())
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
	}
}

func (s *scanner) reset() {
	s.scanT = 0
	s.obj = nil
}

func (s *scanner) clicked(obj *CosmoPoint) {
	if Data.PilotData.Ship.Pos.Sub(obj.Pos).LenSqr() > s.maxRange2 {
		return
	}
	if s.obj == obj {
		return
	}
	s.reset()
	s.obj = obj
}

func (s *scanner) update(ship v2.V2, dt float64) {
	s.scanRange.SetPos(ship)
	s.maxRange2 = Data.SP.Radar.Scan_range * Data.SP.Radar.Scan_range

	if s.obj == nil {
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

	if Data.PilotData.Ship.Pos.Sub(s.obj.Pos).LenSqr() > s.maxRange2 {
		s.reset()
		return
	}
	s.scanT += dt
	if s.totalT > 0 {
		if s.scanT > s.totalT {
			s.procScanned(s.obj)
			s.reset()
		}
	} else {
		s.reset()
	}
}

func (s *scanner) Req() *graph.DrawQueue {
	Q := graph.NewDrawQueue()

	Range := Data.SP.Radar.Scan_range * 2
	s.scanRange.SetSize(Range, Range)
	Q.Add(s.scanRange, graph.Z_UNDER_OBJECT)

	if s.scannedImg != nil {
		Q.Add(s.scannedImg, graph.Z_STAT_HUD)
	}

	if s.obj == nil {
		return Q
	}

	Q.Add(s.scanSector, graph.Z_UNDER_OBJECT)

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
	return Q
}

func (s *scanner) procScanned(obj *CosmoPoint) {
	LogGame("scan", "SCANNED ", obj.ID)
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

	s.scannedImg = qr.NewQRSpriteHUD(gp.ScanData, 256)
	s.scannedImg.SetPivot(graph.TopLeft())
	s.scannedImg.SetPos(graph.ScrP(0, 0))
	time.AfterFunc(time.Second*3, func() { s.scannedImg.TexImageDispose(); s.scannedImg = nil })
}
