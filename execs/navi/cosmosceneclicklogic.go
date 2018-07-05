package main

import (
	."github.com/Shnifer/magellan/log"
	."github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"strings"
)

const (
	minecorptagprefix = "mine~"
)
func (s *cosmoScene) updateControl(dt float64) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mousex, mousey := ebiten.CursorPosition()
		s.procMouseClick(mousex, mousey)
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
		d:=worldPos.Sub(obj.Pos).Len()
		if d < obj.Size ||
			d < draw.Mark_size/s.cam.Scale*graph.GS(){
			s.scanner.clicked(s.objects[id])
			return
		}
	}
	//NAVI MARKER
	Data.NaviData.ActiveMarker = true
	Data.NaviData.MarkerPos = worldPos
	s.naviMarkerT = DEFVAL.NaviMarketDuration
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
		//todo:enter text
		AddBeacon(Data, Client, "just a test beacon")
		ClientLogGame(Client, "ADD BEACON KEY", "just a test beacon")
		Data.NaviData.BeaconCount--
	case "button_orbit":

	case "button_leaveorbit":
		s.scanner.stateZero()
	default:
		if strings.HasPrefix(tag, minecorptagprefix){
			//todo:checks
			corp:=tag[len(minecorptagprefix):]
			s.doneMine(corp)
			s.scanner.stateSelect(s.scanner.obj)
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
	case scanSelect:
		s.cosmoPanels.left.Highlight("")
		s.cosmoPanels.activeLeft(true)
		s.cosmoPanels.activeRight(false)
		s.cosmoPanels.left.Enable()
	case scanProgress:
		s.cosmoPanels.left.Highlight(s.scanner.work)
	case ScanDone:
		switch s.scanner.work {
		case "button_mine":
			s.cosmoPanels.rightMines()
			s.cosmoPanels.activeRight(true)
		case "button_landing":
			s.cosmoPanels.rightLanding()
			s.cosmoPanels.left.Disable()
		case "button_scan":
			s.doneScan()
			s.scanner.stateZero()
		}
	default:
		Log(LVL_ERROR, "Unknown scan state ", scanState)
	}
}

//todo: show signature and name
func (s *cosmoScene) doneScan(){
	id:=s.scanner.obj.ID
	gp, ok := Data.Galaxy.Points[id]
	if !ok {
		return
	}
	key := "scan"
	msg := gp.ScanData
	if msg==""{
		msg = "id: "+id
	}
	if gp.Type == BUILDING_BLACKBOX {
		RequestRemoveBuilding(Client, id)
		key = "blackbox"
	}
	ClientLogGame(Client, key, "SCANNED ", msg)
}

func (s *cosmoScene) doneMine(corp string) {
	//todo: checks and magic detector
	AddMine(Data, Client, s.scanner.obj.ID, corp)
	msg:="planet: "+s.scanner.obj.ID+" corp: "+CompanyNameByOwner(corp)
	ClientLogGame(Client, "mine", msg)
	for i,c:=range Data.NaviData.Mines{
		if c== corp{
			Data.NaviData.Mines = append(Data.NaviData.Mines[:i], Data.NaviData.Mines[i+1:]...)
			return
		}
	}
	Log(LVL_ERROR,"we placed mine that we had not on board")
}