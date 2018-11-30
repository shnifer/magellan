package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/draw"
	"github.com/shnifer/magellan/graph"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/storage"
	"github.com/shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
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

	objectMarks map[string]map[string]string

	//map[id]fullKey
	objectNamesFK map[string]string
	objectNames   map[string]string

	q *graph.DrawQueue

	//drag
	isDragging  bool
	dragLastPos v2.V2

	focus int

	selectedID string
	nameInput  *TextInput

	back *graph.Sprite

	wormHolesT float64
	usedNs     map[int]struct{}
	randomNs   map[int]int

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

	back := NewAtlasSprite("mission_control_back", graph.NoCam)
	back.SetSize(float64(WinW), float64(WinH))
	back.SetPivot(graph.TopLeft())

	res := &scene{
		cam:           cam,
		objects:       make(map[string]*CosmoPoint),
		objectsID:     make([]string, 0),
		objectNamesFK: make(map[string]string),
		objectNames:   make(map[string]string),
		q:             graph.NewDrawQueue(),
		sonar:         sonar,
		sonarBack:     sonarBack,
		sonarPos:      sonarPos,
		back:          back,
		usedNs:        make(map[int]struct{}),
		randomNs:      make(map[int]int),
		objectMarks:   make(map[string]map[string]string),
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
	s.objectNames = make(map[string]string)
	s.objectMarks = make(map[string]map[string]string)
	s.focus = focus_main
	s.showSignature = false
	s.sigs = []commons.Signature{}
	s.usedNs = commons.GetWormHolesNs()
	s.randomNs = make(map[int]int)
	s.wormHolesT = 10

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
		s.objects[gp.ID] = newCP(gp, s.cam.Phys(), s.objectMarks[gp.ID], s.objectNames[gp.ID])
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

	if GalaxyName == commons.WARP_Galaxy_ID && DEFVAL.ShowWormHoles {
		s.wormHolesT += dt
		if s.wormHolesT > 10 {
			s.wormHolesT = 0
			s.setTopCaptions()
		}
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
	if s.back != nil {
		s.q.Add(s.back, graph.Z_BACKGROUND)
	}
	s.q.Run(window)
}

func (s *scene) updatePosition(dt float64) {
	_, wheel := ebiten.Wheel()
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
	if fk, ok := s.objectNamesFK[objectID]; ok {
		objKey, err := storage.ReadKey(fk)
		if err != nil {
			Log(LVL_ERROR, "scene.objectNamesFK strange fullKey", fk, "error", err)
		} else {
			nameDisk.Remove(objKey)
			delete(s.objectNamesFK, objectID)
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
	_, ok := s.objects[rec.planetID]
	if !ok {
		Log(LVL_ERROR, "received name for non-exist planet ", rec.planetID)
		return
	}
	s.objectNames[rec.planetID] = rec.name
	s.remakePoint(rec.planetID)
}

func (s *scene) EventDelName(fk string) {
	for id, key := range s.objectNamesFK {
		if fk == key {
			s.objectNamesFK[id] = ""
			delete(s.objectNames, id)
			s.remakePoint(id)
		}
	}
}

func (s *scene) EventAddBuilding(build commons.Building, fk string) {
	switch build.Type {
	case commons.BUILDING_MINE:
		//for warp map we don't care of planet, Mines shown for systems
		if GalaxyName == commons.WARP_Galaxy_ID {
			build.PlanetID = build.GalaxyID
			build.GalaxyID = commons.WARP_Galaxy_ID
		}
		if _, ok := CurGalaxy.Points[build.PlanetID]; !ok {
			Log(LVL_ERROR, "mine on unknown target name ", build.PlanetID)
		}
		CurGalaxy.AddBuilding(build)
	case commons.BUILDING_BEACON, commons.BUILDING_BLACKBOX:
		if build.GalaxyID != GalaxyName {
			return
		}
		if GalaxyName == commons.WARP_Galaxy_ID {
			var msg string
			if build.Type == commons.BUILDING_BEACON {
				msg = "МАЯК: " + build.Message
			} else {
				msg = "ЧЯ: " + build.Message
			}
			if _, ok := s.objectMarks[build.PlanetID]; ok {
				s.objectMarks[build.PlanetID][fk] = msg
			} else {
				s.objectMarks[build.PlanetID] = make(map[string]string)
				s.objectMarks[build.PlanetID][fk] = msg
			}

		} else {
			CurGalaxy.AddBuilding(build)
		}
	}
	s.remakePoint(build.PlanetID)
}

func (s *scene) EventDelBuilding(build commons.Building, fk string) {
	if GalaxyName == commons.WARP_Galaxy_ID {
		if build.GalaxyID != GalaxyName {
			return
		}
		_, ok := s.objectMarks[build.PlanetID]
		if !ok {
			return
		}
		delete(s.objectMarks[build.PlanetID], fk)
	} else {
		CurGalaxy.DelBuilding(build)
	}
	s.remakePoint(build.PlanetID)
}

func newCP(gp *commons.GalaxyPoint, param graph.CamParams, marks map[string]string, name string) *CosmoPoint {
	res := NewCosmoPoint(gp, param)
	var msg string
	var clr color.Color = color.White

	if name != "" {
		msg += name + "\n"
		clr = captionColor
	}

	if marks != nil {
		arr := make([]string, 0, len(marks))
		for _, s := range marks {
			arr = append(arr, s)
		}
		sort.Sort(sort.StringSlice(arr))
		for _, s := range arr {
			msg += s + "\n"
		}
	}

	if DEFVAL.DebugControl {
		if GalaxyName == commons.WARP_Galaxy_ID {
			msg += "(" + gp.ID + ")"
			if len(gp.Minerals) > 0 {
				msg += " " + fmt.Sprint(gp.Minerals)
			}
		} else if len(gp.Minerals) > 0 {
			msg += fmt.Sprint(gp.Minerals)
			clr = colornames.Red
		}
	}

	msg = strings.TrimSuffix(msg, "\n")
	res.SetCaption(msg, clr)

	return res
}

func (s *scene) remakePoint(id string) {
	changeCP(s.objects, id,
		newCP(CurGalaxy.Points[id], s.cam.Phys(), s.objectMarks[id], s.objectNames[id]))
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

//depricated
func drawWormHoleArrows(Q *graph.DrawQueue, cam *graph.Camera) {
	dirs := commons.GetCurrentWormHoleDirectionSys()
	for _, d := range dirs {
		from, ok := CurGalaxy.Points[d.Src]
		if !ok {
			Log(LVL_WARN, "can't found system ", d.Src)
		}
		to, ok := CurGalaxy.Points[d.Dest]
		if !ok {
			Log(LVL_WARN, "can't found system ", d.Src)
		}
		graph.Arrow(Q, cam, from.Pos, to.Pos, colornames.Red, 30, 30, graph.Z_HUD)
	}
}

func (s *scene) setTopCaptions() {
	dirs := commons.GetCurrentWormHoleDirectionN()
	for sys, d := range dirs {
		dst := d.Dest
		if d.Dest == 0 {
			dst = s.getRandomDest(d.Src)
		} else {
			s.randomNs[d.Src] = 0
		}
		msg := fmt.Sprintf("ИЗЛУЧЕНИЕ:\n%05d%05d", d.Src, dst)
		s.objects[sys].SetCaptionTop(msg, color.White)
	}
}
func (s *scene) getRandomDest(whN int) int {
	v := s.randomNs[whN]
	if v > 0 {
		return v
	}
	var ok bool
	for !ok {
		v = rand.Intn(100000-100) + 100
		_, exist := s.usedNs[v]
		ok = !exist
	}
	s.randomNs[whN] = v
	return v
}
