package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
)

type warpScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	objects map[string]*CosmoPoint

	//trail
	trailT float64
	trail  *graph.FadingArray

	//sonar
	sonarSector *graph.Sector

	//control
	thrustLevel      float64
	maneurLevel      float64
	thrustLevelHUD   *graph.Sprite
	thrustControlHUD *graph.Sprite
	turnLevelHUD     *graph.Sprite
	turnControlHUD   *graph.Sprite
}

func newWarpScene() *warpScene {
	caption := graph.NewText("Warp scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite("ship", cam.FixS())
	ship.SetSize(50, 50)

	sonarSector := graph.NewSector(cam.Phys())
	sonarSector.SetColor(colornames.Forestgreen)

	res := warpScene{
		caption:     caption,
		ship:        ship,
		cam:         cam,
		objects:     make(map[string]*CosmoPoint),
		sonarSector: sonarSector,
	}

	res.trail = graph.NewFadingArray(GetAtlasTex("trail"), trailLifeTime/trailPeriod,
		cam.Deny())

	arrowTex := GetAtlasTex("arrow")
	res.thrustLevelHUD = graph.NewSpriteHUD(arrowTex)
	res.thrustLevelHUD.SetSize(50, 50)
	res.thrustLevelHUD.SetAng(180)
	res.thrustLevelHUD.SetAlpha(0.7)
	res.thrustControlHUD = graph.NewSpriteHUD(arrowTex)
	res.thrustControlHUD.SetSize(50, 50)
	res.thrustControlHUD.SetAlpha(0.5)
	res.turnLevelHUD = graph.NewSpriteHUD(arrowTex)
	res.turnLevelHUD.SetSize(50, 50)
	res.turnLevelHUD.SetAng(-90)
	res.turnLevelHUD.SetAlpha(0.7)
	res.turnControlHUD = graph.NewSpriteHUD(arrowTex)
	res.turnControlHUD.SetSize(50, 50)
	res.turnControlHUD.SetAng(90)
	res.turnControlHUD.SetAlpha(0.5)

	return &res
}

func (s *warpScene) Init() {
	defer LogFunc("warpScene.Init")()

	s.objects = make(map[string]*CosmoPoint)
	s.thrustLevel = 0
	s.maneurLevel = 0

	stateData := Data.GetStateData()

	for id, pd := range stateData.Galaxy.Points {
		cosmoPoint := NewCosmoPoint(pd, s.cam.Phys())
		s.objects[id] = cosmoPoint
	}
}

func (s *warpScene) Update(dt float64) {
	defer LogFunc("warpScene.Update")()
	Data.PilotData.SessionTime += dt

	for _, co := range s.objects {
		co.Update(dt)
	}

	s.updateShipControl(dt)
	s.procShipGravity(dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
		Data.PilotData.Ship.Pos = v2.V2{X: 100, Y: 100}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		s.toCosmo()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= (1 + dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= (1 + dt)
	}

	s.thrustLevelHUD.SetPos(graph.ScrP(0.15, 0.5-0.4*s.thrustLevel))
	s.thrustControlHUD.SetPos(graph.ScrP(0.1, 0.5-0.4*input.GetF("forward")))

	s.turnLevelHUD.SetPos(graph.ScrP(0.5-0.4*s.maneurLevel, 0.15))
	s.turnControlHUD.SetPos(graph.ScrP(0.5-0.4*input.GetF("turn"), 0.1))

	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)

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

	s.sonarSector.SetCenter(Data.PilotData.Ship.Pos)
	s.sonarSector.SetRadius(Data.NaviData.SonarRange)
	s.sonarSector.SetAngles(
		Data.NaviData.SonarDir-Data.NaviData.SonarWide/2,
		Data.NaviData.SonarDir+Data.NaviData.SonarWide/2)

	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	s.camRecalc()
}
func (s *warpScene) camRecalc() {
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.AngleDeg = Data.PilotData.Ship.Ang
	s.cam.Recalc()
}

func (s *warpScene) Draw(image *ebiten.Image) {
	defer LogFunc("warpScene.Draw")()

	Q := graph.NewDrawQueue()

	Q.Add(s.sonarSector, graph.Z_UNDER_OBJECT)

	for _, co := range s.objects {
		Q.Append(co)
	}

	Q.Add(s.trail, graph.Z_UNDER_OBJECT)

	Q.Add(s.ship, graph.Z_GAME_OBJECT)

	Q.Add(s.thrustLevelHUD, graph.Z_HUD)
	Q.Add(s.thrustControlHUD, graph.Z_HUD)
	Q.Add(s.turnLevelHUD, graph.Z_HUD)
	Q.Add(s.turnControlHUD, graph.Z_HUD)

	Q.Add(s.caption, graph.Z_STAT_HUD)

	Q.Run(image)
}

func (s *warpScene) toCosmo() {
	state := Data.State
	state.StateID = STATE_cosmo

	//spawn to closest
	minDist := -1.0
	systemID := ""
	for _, v := range Data.Galaxy.Points {
		dist := Data.PilotData.Ship.Pos.Sub(v.Pos).LenSqr()
		if minDist < 0 || dist < minDist {
			systemID = v.ID
			minDist = dist
		}
	}

	if systemID != "" {
		state.GalaxyID = systemID
		Client.RequestNewState(state.Encode(), false)
	}
}

func (s *warpScene) OnCommand(command string) {
}

func (*warpScene) Destroy() {
}
