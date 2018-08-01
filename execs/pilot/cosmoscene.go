package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"math"
	"sort"
	"time"
)

const trailPeriod = 0.25
const trailLifeTime = 10

type cosmoScene struct {
	ship     *graph.Sprite
	shipMark *graph.Sprite

	caption *graph.Text
	cam     *graph.Camera

	naviMarker *WayPoint

	objects map[string]*CosmoPoint
	objIDs  []string

	lastServerID int
	otherShips   map[string]*OtherShip

	//control
	cruiseOn     bool
	cruiseInput  float64
	maneurDetail bool
	thrustLevel  float64
	maneurLevel  float64

	//trail
	trailT float64
	trail  *graph.FadingArray

	warpEngine *cosmoSceneWarpEngine
	hud        cosmoSceneHUD
	predictors predictors

	//setParams eachUpdate
	gravityAcc    v2.V2
	gravityReport []v2.V2

	//show from Navi
	scanRange  *graph.Sprite
	distCircle *graph.CircleLine

	timerStart   time.Time
	distTravaled float64
	heatBucket   float64

	q *graph.DrawQueue
}

func newCosmoScene() *cosmoScene {
	caption := graph.NewText("Fly scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite(ShipAN, cam.Phys())
	ship.SetSize(ShipSize, ShipSize)

	shipMark := NewAtlasSprite(MARKShipAN, cam.FixS())
	//shipMark.SetSize(50,50)

	marker := NewWayPoint(cam, colornames.Green, true)

	hud := newCosmoSceneHUD(cam)

	scanRange := graph.NewSprite(graph.CircleTex(), cam.Phys())
	scanRange.SetAlpha(0.5)
	scanRange.SetColor(colornames.Indigo)

	clo := graph.CircleLineOpts{
		Layer:  graph.Z_STAT_HUD + 10,
		Clr:    colornames.Oldlace,
		PCount: 64,
	}
	distCircle := graph.NewCircleLine(cam.Center, float64(WinH)*0.3, clo)

	res := cosmoScene{
		caption:    caption,
		ship:       ship,
		shipMark:   shipMark,
		cam:        cam,
		naviMarker: marker,
		hud:        hud,
		objects:    make(map[string]*CosmoPoint),
		objIDs:     make([]string, 0),
		otherShips: make(map[string]*OtherShip),
		scanRange:  scanRange,
		distCircle: distCircle,
		q:          graph.NewDrawQueue(),
	}

	res.trail = graph.NewFadingArray(GetAtlasTex(TrailAN), trailLifeTime/trailPeriod, cam.Deny())

	return &res
}

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	s.objects = make(map[string]*CosmoPoint)
	s.otherShips = make(map[string]*OtherShip)
	s.objIDs = make([]string, 0)
	s.warpEngine = newCosmoSceneWarpEngine()
	s.thrustLevel = 0
	s.cruiseOn = false
	s.maneurLevel = 0
	s.maneurDetail = false
	s.trailT = 0
	s.lastServerID = 0
	s.distTravaled = 0
	s.heatBucket = 0
	s.timerStart = time.Now()
	s.trail.Clear()

	stateData := Data.GetStateData()

	for _, pd := range stateData.Galaxy.Ordered {
		if pd.IsVirtual {
			continue
		}
		cosmoPoint := NewCosmoPoint(pd, s.cam.Phys())
		s.objects[pd.ID] = cosmoPoint
	}

	s.predictors.init(s.cam)

	graph.ClearCache()
}

func (s *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()

	Data.PilotData.FlightTime+=dt
	//received new data about otherShips
	if Data.ServerData.MsgID != s.lastServerID {
		s.actualizeOtherShips()
	}

	//setParams actual otherShips
	for id := range s.otherShips {
		s.otherShips[id].Update(dt)
	}

	//setParams galaxy now to calc right gravity
	//Data.PilotData.SessionTime += dt
	//sessionTime := Data.PilotData.SessionTime
	//Data.Galaxy.Update(sessionTime)
	//s.procShipGravity(dt)
	//Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)
	if !Data.NaviData.IsOrbiting {
		spos := Data.PilotData.Ship.Pos
		UpdateGalaxyAndShip(Data, dt, DEFVAL.DT)
		dpos := Data.PilotData.Ship.Pos.Sub(spos).Len()
		s.distTravaled += dpos
		starheat := 0.0
		heat := s.thrustLevel*Data.SP.March_engine.Thrust_max*Data.SP.March_engine.Heat_prod/100 + starheat - Data.SP.Shields.Heat_sink
		s.heatBucket += heat * dt
		if s.heatBucket < 0 {
			s.heatBucket = 0
		}
	} else {
		Data.PilotData.SessionTime += dt
		Data.Galaxy.Update(Data.PilotData.SessionTime)
		obj := Data.Galaxy.Points[Data.NaviData.OrbitObjectID]
		Data.PilotData.Ship.Pos = obj.Pos
		Data.PilotData.Ship.Vel = v2.ZV
		Data.PilotData.Ship.AngVel = 0
	}

	for id, co := range s.objects {
		if gp, ok := Data.Galaxy.Points[id]; ok {
			s.objects[id].Pos = gp.Pos
		}
		co.Update(dt)
	}

	//we ignore ship control while orbiting
	//same with emissions
	if !Data.NaviData.IsOrbiting {
		s.updateShipControl(dt)
		s.procEmissions(dt)
	}

	//for display draw calls only
	s.gravityAcc, s.gravityReport = SumGravityAccWithReport(Data.PilotData.Ship.Pos, Data.Galaxy,
		0.02)
	s.warpEngine.gravityAcc = s.gravityAcc

	if DEFVAL.DebugControl {
		s.updateDebugControl(dt)
	}

	s.trailUpdate(dt)

	if s.thrustLevel > 0 {
		Data.PilotData.HeatProduction = Data.SP.March_engine.Heat_prod * s.thrustLevel
	} else {
		Data.PilotData.HeatProduction = 0
	}
	s.warpEngine.update(dt)
	s.predictors.setParams()
}
func (s *cosmoScene) camRecalc() {
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.AngleDeg = Data.PilotData.Ship.Ang
	s.cam.Recalc()
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	s.camRecalc()
	s.UpdateHUD()
	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	s.shipMark.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)

	Q := s.q
	Q.Clear()

	Q.Append(s.hud)
	Q.Append(s.warpEngine)

	if len(s.objIDs) != len(s.objects) {
		s.objIDs = make([]string, len(s.objects))
		var i int
		for id := range s.objects {
			s.objIDs[i] = id
			i++
		}
		sort.Strings(s.objIDs)
	}

	for _, id := range s.objIDs {
		Q.Append(s.objects[id])
	}
	Q.Add(s.trail, graph.Z_UNDER_OBJECT)

	s.naviMarker.SetActive(Data.NaviData.ActiveMarker)
	if Data.NaviData.ActiveMarker {
		s.naviMarker.SetShipPoint(Data.PilotData.Ship.Pos, Data.NaviData.MarkerPos)
		Q.Append(s.naviMarker)
	}

	alphaMark, alphaSprite := MarkAlpha(ShipSize/2.0, s.cam)
	if alphaMark > 0 && s.shipMark != nil {
		s.shipMark.SetAlpha(alphaMark)
		Q.Add(s.shipMark, graph.Z_HUD)
	}
	if alphaSprite > 0 && s.ship != nil {
		s.ship.SetAlpha(alphaSprite)
		Q.Add(s.ship, graph.Z_HUD)
	}

	//Q.Add(s.caption, graph.Z_STAT_HUD)

	for _, os := range s.otherShips {
		Q.Append(os)
	}

	if !Data.NaviData.IsOrbiting {
		Q.Append(s.predictors)
	}

	if Data.NaviData.IsScanning {
		var Range float64
		if Data.NaviData.IsDrop {
			Range = Data.SP.Scanner.DropRange * 2
		} else {
			Range = Data.SP.Scanner.ScanRange * 2
		}
		if p, ok := Data.Galaxy.Points[Data.NaviData.ScanObjectID]; ok {
			s.scanRange.SetPos(p.Pos)
			s.scanRange.SetSize(Range, Range)
			Q.Add(s.scanRange, graph.Z_UNDER_OBJECT)
		}
	}

	s.drawScale(Q)
	s.drawGravity(Q)

	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %.0f\nHeat: %.0f/%v\nTravel: %.0f\nTimer: %.1f\nDraws: %v", fps, s.heatBucket/1000,
		Data.SP.Shields.Heat_capacity, s.distTravaled, time.Since(s.timerStart).Seconds(), Q.Len())
	fpsText = graph.NewText(msg, Fonts[Face_list], colornames.Cyan)
	fpsText.SetPosPivot(graph.ScrP(0.1, 0.1), v2.ZV)
	Q.Add(fpsText, graph.Z_STAT_HUD+10)
	Q.Run(image)
}

func (s *cosmoScene) trailUpdate(dt float64) {
	s.trailT += dt
	if s.trailT > trailPeriod {
		s.trailT -= trailPeriod

		s.trail.Add(graph.ArrayElem{
			Size:     5,
			Pos:      Data.PilotData.Ship.Pos,
			LifeTime: trailLifeTime,
		})
	}
	s.trail.Update(dt)
}

func (s *cosmoScene) updateDebugControl(dt float64) {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
		Data.PilotData.Ship.Pos = Data.Galaxy.Points["saturn"].Pos
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		s.predictors.show = !s.predictors.show
	}

	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		v := !ebiten.IsVsyncEnabled()
		ebiten.SetVsyncEnabled(v)
		log.Println("VSync set to ", v)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.timerStart = time.Now()
		s.distTravaled = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		AddBeacon(Data, Client, "just a test beacon")
		ClientLogGame(Client, "ADD BEACON KEY", "just a test beacon")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		toWarp()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= 1 + dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= 1 + dt
	}
	s.cam.Scale = Clamp(s.cam.Scale, 0.0001, 10000)
}

func (*cosmoScene) Destroy() {
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

func (s *cosmoScene) drawScale(Q *graph.DrawQueue) {
	//PosScale factor hud
	camScale := s.cam.Scale * graph.GS()
	maxLen := float64(WinW) * 0.8
	order := math.Floor(math.Log10(maxLen / camScale))
	val := math.Pow10(int(order))
	l := camScale * val

	from := graph.ScrP(0.1, 0.9)
	to := from.AddMul(v2.V2{X: 1, Y: 0}, l)
	mid := from.AddMul(v2.V2{X: 1, Y: 0}, l/2)
	mid.Y += 10

	tick := v2.V2{X: 0, Y: 5}

	graph.LineScr(Q, from, to, colornames.White, graph.Z_STAT_HUD+10)
	graph.LineScr(Q, from.Sub(tick), from.Add(tick), colornames.White, graph.Z_STAT_HUD+10)
	graph.LineScr(Q, to.Sub(tick), to.Add(tick), colornames.White, graph.Z_STAT_HUD+10)

	msg := fmt.Sprintf("%v", val)
	scaleText := graph.NewText(msg, Fonts[Face_mono], colornames.White)
	scaleText.SetPosPivot(mid, graph.TopMiddle())
	Q.Add(scaleText, graph.Z_STAT_HUD+10)

	circleRadPx := float64(WinH) * 0.3
	physRad := circleRadPx / s.cam.Scale / graph.GS()

	Q.Append(s.distCircle)

	msg = fmt.Sprintf("circle radius: %f", physRad)
	physRadText := graph.NewText(msg, Fonts[Face_mono], colornames.Oldlace)
	physRadText.SetPosPivot(graph.ScrP(0.5, 0.4), graph.TopMiddle())
	Q.Add(physRadText, graph.Z_STAT_HUD+10)
}

func (s *cosmoScene) drawGravity(Q *graph.DrawQueue) {
	scale := float64(WinH) * 0.3 / (s.cam.Scale * graph.GS())
	ship := Data.PilotData.Ship.Pos
	thrust := Data.PilotData.ThrustVector
	drawv := func(v v2.V2, clr color.Color) {
		graph.Line(Q, s.cam, ship, ship.AddMul(v, scale), clr, graph.Z_STAT_HUD+10)
	}
	for _, v := range s.gravityReport {
		drawv(v, colornames.Deepskyblue)
	}
	drawv(s.gravityAcc, colornames.Lightblue)
	drawv(Data.PilotData.ThrustVector, colornames.Darkolivegreen)
	drawv(thrust.Add(s.gravityAcc), colornames.White)
}
