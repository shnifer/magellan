package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/skip2/go-qrcode"
	"log"
	"time"
)

type scanner struct {
	totalT    float64
	scanT     float64
	maxRange2 float64
	obj       *CosmoPoint

	countSprite *graph.Sprite
	countN      int

	scannedImg *graph.Sprite
}

func newScanner(cam *graph.Camera) *scanner {
	const countN = 8

	var totalT float64
	if Data.BSP.Scan_speed > 0 {
		totalT = 1 / Data.BSP.Scan_speed
	}
	sprite := graph.NewSprite(GetAtlasTex("trail"), cam, true, true)
	sprite.SetSize(15, 15)

	return &scanner{
		totalT:      totalT,
		maxRange2:   Data.BSP.Scan_range * Data.BSP.Scan_range,
		countSprite: sprite,
		countN:      countN,
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

func (s *scanner) update(dt float64) {
	if s.obj == nil {
		return
	}
	if Data.PilotData.Ship.Pos.Sub(s.obj.Pos).LenSqr() > s.maxRange2 {
		s.reset()
		return
	}
	s.scanT += dt
	if s.scanT > s.totalT {
		s.procScanned(s.obj)
		s.reset()
	}
}

func (s *scanner) Draw(img *ebiten.Image) {
	if s.scannedImg != nil {
		s.scannedImg.Draw(img)
	}
	if s.obj == nil {
		return
	}
	num := int(0.5 + s.scanT/s.totalT*float64(s.countN))
	obj := s.obj.Pos
	rng := s.obj.Size * 1.1
	if rng < 30 {
		rng = 30
	}
	for i := 0; i < num; i++ {
		pos := obj.AddMul(v2.InDir(-360/float64(s.countN)*float64(i)), rng)
		s.countSprite.SetPos(pos)
		s.countSprite.Draw(img)
	}
}

func (s *scanner) procScanned(obj *CosmoPoint) {
	log.Println("SCANNED ", obj.ID)
	if obj.ScanData == "" {
		return
	}
	qr, err := qrcode.New(obj.ScanData, qrcode.Medium)
	if err != nil {
		panic(err)
	}
	image, err := ebiten.NewImageFromImage(qr.Image(256), ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	tex := graph.TexFromImage(image, ebiten.FilterDefault, 0, 0, 0)
	s.scannedImg = graph.NewSprite(tex, nil, false, false)
	s.scannedImg.SetPivot(graph.TopLeft())
	s.scannedImg.SetPos(graph.ScrP(0, 0))
	time.AfterFunc(time.Second*3, func() { s.scannedImg = nil })
}
