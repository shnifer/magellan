package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"sort"
	"github.com/Shnifer/magellan/commons"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/Shnifer/magellan/v2"
	"log"
)

type scene struct{
	cam *graph.Camera
	objects map[string]*draw.CosmoPoint
	objectsID []string

	q *graph.DrawQueue

	//drag
	dragLastPos v2.V2
}

func newScene() *scene{
	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Scale = 1
	cam.Recalc()

	return &scene{
		cam: cam,
		objects: make(map[string]*draw.CosmoPoint),
		objectsID: make([]string,0),
		q: graph.NewDrawQueue(),
	}
}

func (s *scene) init (){
	s.objects = make(map[string]*draw.CosmoPoint, len(CurGalaxy.Ordered))
	s.objectsID = make([]string, 0,len(CurGalaxy.Ordered))
	for _,gp:=range CurGalaxy.Ordered{
		if gp.IsVirtual{
			continue
		}
		s.objects[gp.ID] = draw.NewCosmoPoint(gp, s.cam.Phys())
		s.objectsID = append(s.objectsID, gp.ID)
	}
	sort.Strings(s.objectsID)
}

func (s *scene) update(dt float64){
	CurGalaxy.Update(sessionTime)
	for _,gp:= range CurGalaxy.Ordered{
		s.objects[gp.ID].Pos = gp.Pos
	}

	s.updateposition(dt)
}

func (s *scene) draw(window *ebiten.Image) {
	s.q.Clear()
	for _,id:=range s.objectsID{
		s.q.Append(s.objects[id])
	}
	s.q.Run(window)
}

func (s *scene) updateposition(dt float64){
	_,wheel:=ebiten.MouseWheel()
	if wheel!=0 {
		if wheel > 0 {
			s.cam.Scale *= 1.41
		} else {
			s.cam.Scale /= 1.41
		}
		s.cam.Scale = commons.Clamp(s.cam.Scale, 0.1, 100)
		s.cam.Recalc()
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		//startDrag
		s.dragLastPos = mousePosV()
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		//dragging
		pos:=mousePosV()
		delta := pos.Sub(s.dragLastPos)
		delta.X = -delta.X
		log.Println(s.dragLastPos, pos )
		if delta!=v2.ZV {
			s.dragLastPos = pos
			log.Println(delta, s.cam.Scale)
			s.cam.Pos.DoAddMul(delta, 1/s.cam.Scale)
			log.Println(s.cam.Pos)
			s.cam.Recalc()
		}
	}
}

func mousePosV() v2.V2{
	x,y:=ebiten.CursorPosition()
	return v2.V2{X:float64(x), Y:float64(y)}
}