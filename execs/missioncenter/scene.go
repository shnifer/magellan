package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/storage"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"sort"
	"strconv"
)

const (
	focus_main = iota
	focus_enterName
)

var (
	captionColor = colornames.Yellow
)

type scene struct {
	cam       *graph.Camera
	objects   map[string]*CosmoPoint
	objectsID []string

	//map[id]fullKey
	objectNamesFK map[string]string

	q *graph.DrawQueue

	//drag
	isDragging  bool
	dragLastPos v2.V2

	focus int

	selectedID string
	nameInput  *TextInput

	showSignature bool
	sigs          []commons.Signature
	sonar         *SonarHUD
	sonarBack     *graph.Sprite
	sonarPos      v2.V2
	sonarName     *graph.Text
}

func newScene() *scene {
	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Scale = 1
	cam.Recalc()

	sonarSize := float64(WinH / 3)
	offset := float64(WinH / 10)
	sonarPos := graph.ScrP(1, 0).AddMul(v2.V2{X: -1, Y: 1}, sonarSize/2+offset)
	sonar := NewSonarHUD(sonarPos, sonarSize, graph.NoCam, graph.Z_HUD)

	sonarBack := graph.NewSprite(graph.CircleTex(), graph.NoCam)
	sonarBack.SetPos(sonarPos)
	sonarBack.SetSize(sonarSize, sonarSize)
	sonarBack.SetColor(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	sonarBack.SetAlpha(0.8)

	res := &scene{
		cam:           cam,
		objects:       make(map[string]*CosmoPoint),
		objectsID:     make([]string, 0),
		objectNamesFK: make(map[string]string),
		q:             graph.NewDrawQueue(),
		sonar:         sonar,
		sonarBack:     sonarBack,
		sonarPos:      sonarPos,
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
	s.objectNamesFK = make(map[string]string)
	s.focus = focus_main
	s.showSignature = false
	s.sigs = []commons.Signature{}
	if GalaxyName == commons.WARP_Galaxy_ID {
		s.cam.Scale = 0.001
		s.cam.Pos = CurGalaxy.Points["solar"].Pos
	} else {
		s.cam.Pos = v2.ZV
		s.cam.Scale = 0.01
	}
	s.cam.Recalc()

	for _, gp := range CurGalaxy.Ordered {
		if gp.IsVirtual {
			continue
		}
		s.objects[gp.ID] = newCP(gp, s.cam.Phys())
		s.objectsID = append(s.objectsID, gp.ID)
	}

	for objKey, data := range buildingData {
		build, err := commons.Building{}.Decode([]byte(data))
		if err != nil {
			Log(LVL_ERROR, "can't decode building data ", data)
			continue
		}
		fk := objKey.FullKey()
		s.EventAddBuilding(build, fk)
	}

	for objKey, data := range namesData {
		rec, err := nameRec{}.decode(data)
		if err != nil {
			Log(LVL_ERROR, "can't decode objKey ", objKey, " data ", data, " ", err)
			continue
		}
		s.EventAddName(rec, objKey.FullKey())
	}
	sort.Strings(s.objectsID)
}

func (s *scene) update(dt float64) {
	CurGalaxy.Update(sessionTime)
	for _, gp := range CurGalaxy.Ordered {
		obj := s.objects[gp.ID]
		if obj == nil {
			continue
		}
		obj.Pos = gp.Pos
		obj.Update(dt)
	}

	s.updatePosition(dt)
	if s.focus == focus_enterName {
		s.nameInput.Update(dt)
	}

	s.sonar.ActiveSignatures(s.sigs)
	s.sonar.Update(dt)
}

func (s *scene) draw(window *ebiten.Image) {
	s.q.Clear()
	t := graph.NewText(GalaxyName, Fonts[Face_cap], color.White)
	t.SetPosPivot(graph.ScrP(0.1, 0.2), v2.ZV)
	s.q.Add(t, graph.Z_STAT_HUD)
	for _, id := range s.objectsID {
		s.q.Append(s.objects[id])
	}
	if s.focus == focus_enterName {
		s.q.Append(s.nameInput)
	}
	if DEFVAL.DebugControl {
		if s.showSignature {
			s.q.Add(s.sonarBack, graph.Z_HUD-1)
			s.q.Append(s.sonar)
			s.q.Add(s.sonarName, graph.Z_STAT_HUD)
		}
	}
	s.q.Run(window)
}

func (s *scene) updatePosition(dt float64) {
	_, wheel := ebiten.MouseWheel()
	if wheel != 0 {
		newScale := s.cam.Scale
		if wheel > 0 {
			newScale *= 1.19
		} else {
			newScale /= 1.19
		}
		newScale = commons.Clamp(newScale, 0.00001, 1000)
		x, y := ebiten.CursorPosition()
		mp := s.cam.UnApply(v2.V2{X: float64(x), Y: float64(y)})
		v := s.cam.Pos.Sub(mp)
		s.cam.Pos = mp.AddMul(v, s.cam.Scale/newScale)
		s.cam.Scale = newScale
		s.cam.Recalc()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if pressed, id := s.objectClicked(); pressed {
			if GalaxyName == commons.WARP_Galaxy_ID {
				s.procWarpObjectClick(id)
			} else {
				s.procStarObjectClick(id)
			}
		} else {
			//startDrag
			s.dragLastPos = mousePosV()
			s.isDragging = true
			s.showSignature = false
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if pressed, id := s.objectClicked(); GalaxyName == commons.WARP_Galaxy_ID && pressed {
			s.procWarpRightObjectClick(id)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && GalaxyName != commons.WARP_Galaxy_ID {
		changeGalaxy <- commons.WARP_Galaxy_ID
	}
}

func mousePosV() v2.V2 {
	x, y := ebiten.CursorPosition()
	return v2.V2{X: float64(x), Y: float64(y)}
}

func (s *scene) objectClicked() (pressed bool, id string) {
	p := s.cam.UnApply(mousePosV())
	for id, obj := range s.objects {
		dist := math.Max(obj.Size, 10/s.cam.Scale)
		if obj.Pos.Sub(p).Len() < dist {
			return true, id
		}
	}
	return false, ""
}

func (s *scene) procWarpObjectClick(id string) {
	s.focus = focus_enterName
	s.selectedID = id
	caption := s.objects[id].GetCaption()
	s.nameInput.SetText(caption)
	s.showSignature = true
	s.sigs = CurGalaxy.Points[id].Signatures
	t := graph.NewText(id, Fonts[Face_cap], color.White)
	t.SetPosPivot(s.sonarPos, graph.Center())
	s.sonarName = t
}

func (s *scene) procStarObjectClick(id string) {
	s.showSignature = true
	s.sigs = CurGalaxy.Points[id].Signatures
	t := graph.NewText(id, Fonts[Face_cap], color.White)
	t.SetPosPivot(s.sonarPos, graph.Center())
	s.sonarName = t
}

func (s *scene) procWarpRightObjectClick(id string) {
	if DEFVAL.DebugControl {
		changeGalaxy <- id
	}
}

func (s *scene) onNameTextInput(text string, done bool) {
	s.focus = focus_main
	if done {
		s.requestNewName(GalaxyName, s.selectedID, text)
	}
	s.selectedID = ""
	s.showSignature = false
}

func (s *scene) requestNewName(galaxyName, objectID, newName string) {
	s.objects[objectID].SetCaption(newName, captionColor)
	if fk, ok := s.objectNamesFK[objectID]; ok {
		objKey, err := storage.ReadKey(fk)
		if err != nil {
			Log(LVL_ERROR, "scene.objectNamesFK strange fullKey", fk, "error", err)
		} else {
			nameDisk.Remove(objKey)
		}
	}

	key := objectID + " " + strconv.Itoa(nameDisk.NextID())
	rec := nameRec{
		planetID: objectID,
		name:     newName,
	}
	nameDisk.Add(galaxyName, key, rec.encode())
}

func (s *scene) EventAddName(rec nameRec, fk string) {
	s.objectNamesFK[rec.planetID] = fk
	obj, ok := s.objects[rec.planetID]
	if !ok {
		Log(LVL_ERROR, "received name for non-exist planet ", rec.planetID)
		return
	}
	obj.SetCaption(rec.name, captionColor)
}

func (s *scene) EventDelName(fk string) {
	for id, key := range s.objectNamesFK {
		if fk == key {
			s.objects[id].SetCaption("", captionColor)
			s.objectNamesFK[id] = ""
		}
	}
}

func (s *scene) EventAddBuilding(build commons.Building, fk string) {
	if build.Type != commons.BUILDING_MINE {
		return
	}
	//for warp map we don't care of planet, Mines shown for systems
	if GalaxyName == commons.WARP_Galaxy_ID {
		build.PlanetID = build.GalaxyID
		build.GalaxyID = commons.WARP_Galaxy_ID
	}
	if _, ok := CurGalaxy.Points[build.PlanetID]; !ok {
		Log(LVL_ERROR, "mine on unknown target name ", build.PlanetID)
	}
	CurGalaxy.AddBuilding(build)
	changeCP(s.objects, build.PlanetID,
		newCP(CurGalaxy.Points[build.PlanetID], s.cam.Phys()))
}

func (s *scene) EventDelBuilding(build commons.Building, fk string) {
	CurGalaxy.DelBuilding(build)
	changeCP(s.objects, build.PlanetID,
		newCP(CurGalaxy.Points[build.PlanetID], s.cam.Phys()))
}

func newCP(gp *commons.GalaxyPoint, param graph.CamParams) *CosmoPoint {
	res := NewCosmoPoint(gp, param)
	if DEFVAL.DebugControl {
		if GalaxyName == commons.WARP_Galaxy_ID {
			msg := gp.ID
			if len(gp.Minerals) > 0 {
				msg = msg + " " + fmt.Sprint(gp.Minerals)
			}
			res.SetCaption(msg, color.White)
		} else if len(gp.Minerals) > 0 {
			res.SetCaption(fmt.Sprint(gp.Minerals), colornames.Red)
		}
	}
	return res
}

func changeCP(objs map[string]*CosmoPoint, id string, point *CosmoPoint) {
	if point == nil {
		Log(LVL_ERROR, "scene change CosmoPoint with nil value")
		return
	}
	if _, ok := objs[id]; ok {
		*objs[id] = *point
	} else {
		objs[id] = point
	}
}
