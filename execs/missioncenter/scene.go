package main

import (
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"log"
	"sort"
)

const (
	focus_main = iota
	focus_enterName
)

type scene struct {
	cam       *graph.Camera
	objects   map[string]*CosmoPoint
	objectsID []string

	q *graph.DrawQueue

	//drag
	isDragging  bool
	dragLastPos v2.V2

	focus int

	selectedID string
	nameInput  *TextInput
}

func newScene() *scene {
	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Scale = 1
	cam.Recalc()

	res := &scene{
		cam:       cam,
		objects:   make(map[string]*CosmoPoint),
		objectsID: make([]string, 0),
		q:         graph.NewDrawQueue(),
	}

	textPanel := NewAtlasSprite(commons.TextPanelAN, graph.NoCam)
	textPanel.SetPos(graph.ScrP(0.5, 0))
	textPanel.SetPivot(graph.TopMiddle())
	size := graph.ScrP(0.6, 0.1)
	textPanel.SetSize(size.X, size.Y)

	res.nameInput = NewTextInput(textPanel, Fonts[Face_cap], colornames.White, graph.Z_HUD+1, res.onNameTextInput)

	return res
}

func (s *scene) init() {
	s.objects = make(map[string]*CosmoPoint, len(CurGalaxy.Ordered))
	s.objectsID = make([]string, 0, len(CurGalaxy.Ordered))
	for _, gp := range CurGalaxy.Ordered {
		if gp.IsVirtual {
			continue
		}
		s.objects[gp.ID] = NewCosmoPoint(gp, s.cam.Phys())
		s.objectsID = append(s.objectsID, gp.ID)
	}
	sort.Strings(s.objectsID)
}

func (s *scene) update(dt float64) {
	CurGalaxy.Update(sessionTime)
	for _, gp := range CurGalaxy.Ordered {
		s.objects[gp.ID].Pos = gp.Pos
	}

	s.updateposition(dt)
	if s.focus == focus_enterName {
		s.nameInput.Update(dt)
	}
}

func (s *scene) draw(window *ebiten.Image) {
	s.q.Clear()
	for _, id := range s.objectsID {
		s.q.Append(s.objects[id])
	}
	if s.focus == focus_enterName {
		s.q.Append(s.nameInput)
	}
	s.q.Run(window)
}

func (s *scene) updateposition(dt float64) {
	_, wheel := ebiten.MouseWheel()
	if wheel != 0 {
		if wheel > 0 {
			s.cam.Scale *= 1.19
		} else {
			s.cam.Scale /= 1.19
		}
		s.cam.Scale = commons.Clamp(s.cam.Scale, 0.1, 100)
		s.cam.Recalc()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if pressed, id := s.objectClicked(); pressed {
			s.procObjectClick(id)
		} else {
			//startDrag
			s.dragLastPos = mousePosV()
			s.isDragging = true
		}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.isDragging = false
	} else if s.isDragging {
		//dragging
		pos := mousePosV()
		delta := pos.Sub(s.dragLastPos)
		delta.X = -delta.X
		if delta != v2.ZV {
			s.dragLastPos = pos
			s.cam.Pos.DoAddMul(delta, 1/s.cam.Scale)
			s.cam.Recalc()
		}
	}
}

func mousePosV() v2.V2 {
	x, y := ebiten.CursorPosition()
	return v2.V2{X: float64(x), Y: float64(y)}
}

func (s *scene) objectClicked() (pressed bool, id string) {
	p := s.cam.UnApply(mousePosV())
	for id, obj := range s.objects {
		if obj.Pos.Sub(p).LenSqr() < obj.Size*obj.Size {
			return true, id
		}
	}
	return false, ""
}

func (s *scene) procObjectClick(id string) {
	log.Println("clicked ", id)
	s.focus = focus_enterName
	s.selectedID = id
}

func (s *scene) onNameTextInput(text string, done bool) {
	s.focus = focus_main
	if done {
		s.nameObject(GalaxyName, s.selectedID, text)
	}
	s.selectedID = ""
}

func (s *scene) nameObject(galaxyName, objectID, newName string) {
	s.objects[s.selectedID].SetCaption(newName, colornames.White)
	nameDisk.Add(galaxyName, objectID, newName)
}
