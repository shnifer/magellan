package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"strings"
)

const (
	minecorptagprefix = "mine~"
)

func (s *cosmoScene) updateControl(dt float64) {
	if s.inputFocus == inputText {
		s.textInput.Update(dt)
	}
	if s.inputFocus == inputMain {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mousex, mousey := ebiten.CursorPosition()
			s.procMouseClick(mousex, mousey)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			Data.NaviData.ActiveMarker = false
		}
	}

	s.scanner.update(s.shipRB.RB().Pos, dt)
	s.cosmoPanels.update(dt)
}

func (s *cosmoScene) procMouseClick(x, y int) {
	//PANELS
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
	//COSMOOBJECTS
	worldPos := s.cam.UnApply(v2.V2{X: float64(x), Y: float64(y)})
	for id, obj := range Data.Galaxy.Points {
		d := worldPos.Sub(obj.Pos).Len()
		if d < obj.Size ||
			d < draw.Mark_size/s.cam.Scale*graph.GS() {
			really := s.scanner.clicked(s.objects[id])
			if really {
				return
			}
		}
	}
	//NAVI MARKER
	Data.NaviData.ActiveMarker = true
	Data.NaviData.MarkerPos = worldPos
}

func (s *cosmoScene) procButtonClick(tag string) {
	switch tag {
	case "button_mine":
		s.scanner.Start(tag)
	case "button_landing":
		s.scanner.Start(tag)
	case "button_scan":
		s.scanner.Start(tag)
	case "button_beacon":
		s.startBeaconTextInput()
	case "button_orbit":
		if !Data.NaviData.IsOrbiting {
			Data.NaviData.IsOrbiting = true
			Data.NaviData.OrbitObjectID = s.scanner.obj.ID
			ClientLogGame(Client, "landing", s.scanner.obj.ID)
			found := false
			for i, v := range Data.NaviData.Landing {
				if v == Data.NaviData.OrbitObjectID {
					Data.NaviData.Landing =
						append(Data.NaviData.Landing[:i], Data.NaviData.Landing[i+1:]...)
					found = true
					break
				}
			}
			if !found {
				Log(LVL_ERROR, "orbiting done without needed landing module")
			}
		}

	case "button_leaveorbit":
		Data.NaviData.IsOrbiting = false
		s.scanner.stateZero()
	default:
		if strings.HasPrefix(tag, minecorptagprefix) {
			msg, ok := s.checkMine()
			if ok {
				corp := tag[len(minecorptagprefix):]
				s.doneMine(corp)
				s.scanner.stateSelect(s.scanner.obj)
			} else {
				s.announce.AddMsg(msg, colornames.Red, 2)
			}
		} else {
			Log(LVL_ERROR, "Unknown button tag ", tag)
		}
	}
}

func (s *cosmoScene) scanState(scanState int) {
	switch scanState {
	case scanZero:
		s.cosmoPanels.left.Highlight("")
		s.cosmoPanels.activeLeft(false)
		s.cosmoPanels.activeRight(false)
		s.cosmoPanels.left.Enable()
		Data.NaviData.IsScanning = false
	case scanSelect:
		s.cosmoPanels.left.Highlight("")
		s.cosmoPanels.activeLeft(true)
		s.cosmoPanels.activeRight(false)
		s.cosmoPanels.left.Enable()
		Data.NaviData.IsScanning = true
		Data.NaviData.ScanObjectID = s.scanner.obj.ID
	case scanProgress:
		s.cosmoPanels.left.Highlight(s.scanner.work)
	case ScanDone:
		switch s.scanner.work {
		case "button_mine":
			msg, ok := s.checkMine()
			if ok {
				s.cosmoPanels.rightMines()
				s.cosmoPanels.activeRight(true)
				s.announce.AddMsg(msg, colornames.Green, 2)
			} else {
				s.announce.AddMsg(msg, colornames.Red, 2)
				s.scanner.stateZero()
			}
		case "button_landing":
			if s.checkLanding() {
				s.cosmoPanels.rightLanding()
				s.cosmoPanels.activeRight(true)
				s.cosmoPanels.left.Disable()
			} else {
				s.announce.AddMsg("высадка невозможна", colornames.Red, 2)
				s.scanner.stateZero()
			}
		case "button_scan":
			s.doneScan()
			s.scanner.stateZero()
		}
	default:
		Log(LVL_ERROR, "Unknown scan state ", scanState)
	}
}

//todo: show signature and name
func (s *cosmoScene) doneScan() {
	id := s.scanner.obj.ID
	gp, ok := Data.Galaxy.Points[id]
	if !ok {
		return
	}
	key := "scan"
	msg := gp.ScanData
	if msg == "" {
		msg = "id: " + id
	}
	if gp.Type == BUILDING_BLACKBOX {
		RequestRemoveBuilding(Client, id)
		key = "blackbox"
	}
	ClientLogGame(Client, key, "SCANNED ", msg)
}

func (s *cosmoScene) doneMine(corp string) {
	AddMine(Data, Client, s.scanner.obj.ID, corp)
	msg := "planet " + s.scanner.obj.ID + " corp " + CompanyNameByOwner(corp)
	ClientLogGame(Client, "mine", msg)
	for i, c := range Data.NaviData.Mines {
		if c == corp {
			Data.NaviData.Mines = append(Data.NaviData.Mines[:i], Data.NaviData.Mines[i+1:]...)
			return
		}
	}
	Log(LVL_ERROR, "we placed mine that we had not on board")
}

func (s *cosmoScene) checkLanding() bool {
	found := false
	for _, id := range Data.NaviData.Landing {
		if s.scanner.obj.ID == id {
			found = true
		}
	}
	return found
}

func (s *cosmoScene) checkMine() (msg string, ok bool) {
	gp := Data.Galaxy.Points[s.scanner.obj.ID]
	if len(gp.Mines) > 0 {
		return "у нас есть a mine already", false
	}

	hasknown := make([]string, 0)
	var hasunknown bool

	var know bool
	for _, has := range gp.Minerals {
		know = false
		for _, known := range Data.BSP.KnownMinerals {
			if has == known.ID {
				know = true
				hasknown = append(hasknown, known.UserName)
				break
			}
		}
		if !know {
			hasunknown = true
		}
	}

	if len(hasknown) == 0 {
		if hasunknown {
			return "неизвестные minerals", false
		} else {
			return "nothing здесь нет", false
		}
	}
	msg = "Добыча: " + strings.Join(hasknown, ",")
	if hasunknown {
		msg += "\nтакже неизвестные minerals"
	}
	return msg, true
}

func (s *cosmoScene) onBeaconTextInput(text string, done bool) {
	s.inputFocus = inputMain
	if !done {
		return
	}
	//place beacon
	AddBeacon(Data, Client, text)
	ClientLogGame(Client, "beacon", text)
	Data.NaviData.BeaconCount--
}

func (s *cosmoScene) startBeaconTextInput() {
	s.inputFocus = inputText
	s.textInput.SetText("")
}