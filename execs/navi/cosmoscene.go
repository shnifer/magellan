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

const (
	inputMain = 0
	inputText = 1
)

type cosmoScene struct {
	ship     *graph.Sprite
	shipMark *graph.Sprite
	shipRB   *commons.RBFollower

	sessionTime  *commons.SessionTime
	lastPilotMsg int

	lastServerID int
	otherShips   map[string]*OtherShip

	caption *graph.Text
	cam     *graph.Camera

	isCamToShip bool

	objects map[string]*CosmoPoint

	scanner *scanner

	predictors  predictors
	cosmoPanels *cosmoPanels

	naviMarker *WayPoint
	shipMarker *WayPoint

	announce *AnnounceText

	background *graph.Sprite
	f9         *graph.Frame9HUD

	inputFocus int
	textInput  *TextInput

	q *graph.DrawQueue
}

func newCosmoScene() *cosmoScene {
	var background *graph.Sprite

	if !DEFVAL.LowQ {
		background = NewAtlasSpriteHUD(commons.DefaultBackgroundAN)
		background.SetSize(float64(WinW), float64(WinH))
		background.SetPivot(graph.TopLeft())
		background.SetColor(colornames.Dimgrey)
	}

	caption := graph.NewText("Navi scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite(commons.ShipAN, cam.Phys())
	ship.SetSize(commons.ShipSize, commons.ShipSize)

	shipMark := NewAtlasSprite(commons.MARKShipAN, cam.FixS())

	cosmoPanels := newCosmoPanels()

	at := NewAnnounceText(graph.ScrP(0.5, 0.3), graph.Center(),
		Fonts[Face_cap], graph.Z_STAT_HUD)
	textPanel := NewAtlasSprite(commons.TextPanelAN, graph.NoCam)
	textPanel.SetPos(graph.ScrP(0.5, 0))
	textPanel.SetPivot(graph.TopMiddle())
	size := graph.ScrP(0.6, 0.1)
	textPanel.SetSize(size.X, size.Y)

	naviMarker:=NewWayPoint(cam, colornames.Green, true)
	shipMarker:=NewWayPoint(cam, colornames.Yellow, false)
	shipMarker.SetActive(true)

	f9 := NewAtlasFrame9HUD(commons.Frame9AN, WinW, WinH, graph.Z_HUD-1)

	scene := &cosmoScene{
		caption:     caption,
		ship:        ship,
		shipMark:    shipMark,
		cam:         cam,
		naviMarker: naviMarker,
		shipMarker: shipMarker,
		cosmoPanels: cosmoPanels,
		objects:     make(map[string]*CosmoPoint),
		otherShips:  make(map[string]*OtherShip),
		announce:    at,
		background:  background,
		f9:          f9,
		q:           graph.NewDrawQueue(),
	}
	scene.textInput = NewTextInput(textPanel, Fonts[Face_cap], colornames.White, graph.Z_HUD+1, scene.onBeaconTextInput)

	return scene
}

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	stateData := Data.GetStateData()

	s.objects = make(map[string]*CosmoPoint)
	s.otherShips = make(map[string]*OtherShip)
	s.lastServerID = 0
	s.isCamToShip = true
	s.scanner = newScanner(s.cam, s.scanState)
	s.shipRB = commons.NewRBFollower(float64(DEFVAL.PingPeriod) / 1000)
	s.sessionTime = commons.NewSessionTime(Data.PilotData.SessionTime)

	s.predictors.init(s.cam)

	for id, pd := range stateData.Galaxy.Points {
		if pd.IsVirtual {
			continue
		}
		cosmoPoint := NewCosmoPoint(pd, s.cam.Phys())
		s.objects[id] = cosmoPoint
	}

	graph.ClearCache()
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
	sessionTime := s.sessionTime.Get()
	Data.Galaxy.Update(sessionTime)

	if Data.ServerData.MsgID != s.lastServerID {
		s.actualizeOtherShips()
	}

	s.shipRB.Update(dt)
	ship := s.shipRB.RB()
	if Data.NaviData.IsOrbiting {
		ship.Pos = Data.Galaxy.Points[Data.NaviData.OrbitObjectID].Pos
		ship.Vel = v2.ZV
		ship.AngVel = 0
	}

	//update actual otherShips
	for id := range s.otherShips {
		s.otherShips[id].Update(dt)
	}

	for id, co := range s.objects {
		if gp, ok := Data.Galaxy.Points[id]; ok {
			s.objects[id].Pos = gp.Pos
		}
		co.Update(dt)
	}

	if s.inputFocus == inputMain {
		s.updateInputMain()
	}

	s.ship.SetPosAng(ship.Pos, ship.Ang)
	s.shipMark.SetPosAng(ship.Pos, ship.Ang)

	s.predictors.setParams(sessionTime, ship)
	s.updateControl(dt)
	s.announce.Update(dt)

	s.naviMarker.SetShipPoint(s.cam.Pos, Data.NaviData.MarkerPos)
	s.shipMarker.SetShipPoint(s.cam.Pos, ship.Pos)

	if s.isCamToShip {
		s.cam.Pos = ship.Pos
		s.cam.Recalc()
	}
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	Q := s.q
	Q.Clear()

	if s.background != nil {
		Q.Add(s.background, graph.Z_STAT_BACKGROUND)
		Q.Append(s.f9)
	}

	Q.Append(s.scanner)

	for _, co := range s.objects {
		Q.Append(co)
	}

	for _, os := range s.otherShips {
		Q.Append(os)
	}

	Q.Append(s.predictors)

	s.naviMarker.SetActive(Data.NaviData.ActiveMarker)
	if Data.NaviData.ActiveMarker {
		Q.Append(s.naviMarker)
	}
	Q.Append(s.shipMarker)

	alphaMark, alphaSprite := MarkAlpha(commons.ShipSize/2.0, s.cam)
	if alphaMark > 0 && s.shipMark != nil {
		s.shipMark.SetAlpha(alphaMark)
		Q.Add(s.shipMark, graph.Z_HUD)
	}
	if alphaSprite > 0 && s.ship != nil {
		s.ship.SetAlpha(alphaSprite)
		Q.Add(s.ship, graph.Z_HUD)
	}

	if s.inputFocus == inputText {
		Q.Append(s.textInput)
	}

	Q.Append(s.cosmoPanels)
	Q.Append(s.announce)

	Q.Run(image)
}

func (s *cosmoScene) actualizeOtherShips() {
	s.lastServerID = Data.ServerData.MsgID

	//Create new otherShip and move all to new positions
	for _, otherData := range Data.ServerData.OtherShips {
		otherShip, ok := s.otherShips[otherData.Id]
		if !ok {
			otherShip = NewOtherShip(s.cam.FixS(), otherData.Name, float64(DEFVAL.OtherShipElastic)/1000)
			s.otherShips[otherData.Id] = otherShip
		}
		otherShip.SetRB(otherData.Ship)
	}

	//check for lost otherShips to delete
	for id := range s.otherShips {
		found := false
		for _, otherData := range Data.ServerData.OtherShips {
			if otherData.Id == id {
				found = true
				break
			}
		}
		if !found {
			delete(s.otherShips, id)
		}
	}
}

func (*cosmoScene) Destroy() {
}

func (s *cosmoScene) updateInputMain() {
	moveScale := 10 / s.cam.Scale
	_,wy:=ebiten.MouseWheel()
	if wy>0{
		s.cam.Scale *= 1.41
		s.cam.Recalc()
	} else if wy<0{
		s.cam.Scale /= 1.41
		s.cam.Recalc()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ){
		s.cam.Scale *= 1 + dt
		s.cam.Recalc()
	}
	if ebiten.IsKeyPressed(ebiten.KeyE){
		s.cam.Scale /= 1 + dt
		s.cam.Recalc()
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		s.cam.Pos.DoAddMul(v2.V2{X: 0, Y: 1}, moveScale)
		s.isCamToShip = false
		s.cam.Recalc()
	}
	if ebiten.IsKeyPressed(ebiten.KeyS)|| ebiten.IsKeyPressed(ebiten.KeyDown){
		s.cam.Pos.DoAddMul(v2.V2{X: 0, Y: -1}, moveScale)
		s.isCamToShip = false
		s.cam.Recalc()
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft){
		s.cam.Pos.DoAddMul(v2.V2{X: -1, Y: 0}, moveScale)
		s.isCamToShip = false
		s.cam.Recalc()
	}
	if ebiten.IsKeyPressed(ebiten.KeyD)|| ebiten.IsKeyPressed(ebiten.KeyRight) {
		s.cam.Pos.DoAddMul(v2.V2{X: 1, Y: 0}, moveScale)
		s.isCamToShip = false
		s.cam.Recalc()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.isCamToShip = true
	}
}
