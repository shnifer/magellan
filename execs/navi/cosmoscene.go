package main

import (
	"github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
)

type cosmoScene struct {
	ship         *graph.Sprite
	shipRB       *commons.RBFollower
	sessionTime  *commons.SessionTime
	lastPilotMsg int

	caption *graph.Text
	cam     *graph.Camera

	objects map[string]*CosmoPoint

	scanner *scanner

	naviMarkerT float64
}

func newCosmoScene() *cosmoScene {
	caption := graph.NewText("Navi scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite("ship", cam.Phys())
	ship.SetSize(50, 50)

	return &cosmoScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		objects: make(map[string]*CosmoPoint),
	}
}

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	stateData := Data.GetStateData()

	s.objects = make(map[string]*CosmoPoint)
	s.naviMarkerT = 0
	s.scanner = newScanner(s.cam)
	s.shipRB = commons.NewRBFollower(float64(DEFVAL.PingPeriod) / 1000)
	s.sessionTime = commons.NewSessionTime(Data.PilotData.SessionTime)

	for id, pd := range stateData.Galaxy.Points {
		cosmoPoint := NewCosmoPoint(pd, s.cam.Phys())
		s.objects[id] = cosmoPoint
	}
}

func (s *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()
	//PilotData Rigid Body emulation
	if Data.PilotData.MsgID != s.lastPilotMsg {
		s.shipRB.MoveTo(Data.PilotData.Ship)
		s.sessionTime.MoveTo(Data.PilotData.SessionTime)
		s.lastPilotMsg = Data.PilotData.MsgID
	}
	s.sessionTime.Update(dt)
	Data.Galaxy.Update(s.sessionTime.Get())

	s.shipRB.Update(dt)
	ship := s.shipRB.RB()

	s.cam.Pos = ship.Pos
	s.cam.Recalc()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mousex, mousey := ebiten.CursorPosition()
		s.procMouseClick(v2.V2{X: float64(mousex), Y: float64(mousey)})
	}
	s.naviMarkerT -= dt
	if s.naviMarkerT < 0 {
		s.naviMarkerT = 0
		Data.NaviData.ActiveMarker = false
	}

	for id, co := range s.objects {
		if gp, ok := Data.Galaxy.Points[id]; ok {
			s.objects[id].Pos = gp.Pos
		}
		co.Update(dt)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= 1 + dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= 1 + dt
	}

	s.ship.SetPosAng(ship.Pos, ship.Ang)
	s.scanner.update(ship.Pos, dt)
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	Q := graph.NewDrawQueue()

	Q.Append(s.scanner)

	for _, co := range s.objects {
		Q.Append(co)
	}

	//Q.Add(s.caption, graph.Z_STAT_HUD)
	Q.Add(s.ship, graph.Z_HUD)

	Q.Run(image)
}

func (s *cosmoScene) procMouseClick(scrPos v2.V2) {
	worldPos := s.cam.UnApply(scrPos)
	for id, obj := range Data.Galaxy.Points {
		if worldPos.Sub(obj.Pos).LenSqr() < (obj.Size * obj.Size) {
			s.scanner.clicked(s.objects[id])
			return
		}
	}
	Data.NaviData.ActiveMarker = true
	Data.NaviData.MarkerPos = worldPos
	s.naviMarkerT = DEFVAL.NaviMarketDuration
}

func (s *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}
