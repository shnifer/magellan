package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
)

const trailPeriod = 0.25
const trailLifeTime = 10

type cosmoScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	naviMarker *graph.Sprite

	objects []*CosmoPoint
	idMap   map[string]*CosmoPoint

	//control
	thrustLevel float64
	maneurLevel float64

	//trail
	trailT float64
	trail  *graph.FadingArray

	background *graph.Sprite
	compass    *graph.Sprite

	//hud
	thrustLevelHUD   *graph.Sprite
	thrustControlHUD *graph.Sprite
	turnLevelHUD     *graph.Sprite
	turnControlHUD   *graph.Sprite
}

func newCosmoScene() *cosmoScene {
	caption := graph.NewText("Fly scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite("ship", cam, true, false)
	ship.SetSize(50, 50)

	marker := NewAtlasSprite("marker", cam, true, true)
	marker.SetPivot(graph.MiddleBottom())

	background := NewAtlasSpriteHUD("background")
	background.SetSize(float64(WinW), float64(WinH))
	background.SetPivot(graph.TopLeft())
	background.SetColor(colornames.Dimgrey)

	compass := NewAtlasSprite("compass", cam, true, false)
	compassSize := float64(WinH) * 0.8
	compass.SetSize(compassSize, compassSize)
	compass.SetAlpha(0.5)

	res := cosmoScene{
		caption:    caption,
		ship:       ship,
		cam:        cam,
		naviMarker: marker,
		objects:    make([]*CosmoPoint, 0),
		idMap:      make(map[string]*CosmoPoint),
		background: background,
		compass:    compass,
	}

	res.trail = graph.NewFadingArray(GetAtlasTex("trail"), trailLifeTime/trailPeriod, cam, true, true)

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

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	s.objects = s.objects[:0]
	s.idMap = make(map[string]*CosmoPoint)
	s.thrustLevel = 0
	s.maneurLevel = 0
	s.trailT = 0
	s.trail.Clear()

	stateData := Data.GetStateData()

	for _, pd := range stateData.Galaxy.Points {
		cosmoPoint := NewCosmoPoint(pd, s.cam)
		s.objects = append(s.objects, cosmoPoint)

		if pd.ID != "" {
			s.idMap[pd.ID] = cosmoPoint
		}
		if pd.ParentID != "" {
			cosmoPoint.Parent = s.idMap[pd.ParentID]
		}
	}
}

func (s *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()
	Data.PilotData.SessionTime += dt
	sessionTime := Data.PilotData.SessionTime
	for _, co := range s.objects {
		co.Update(sessionTime)
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
		Data.PilotData.Ship.Pos = s.idMap["magellan"].Pos
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		s.toWarp()
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
	s.compass.SetPos(Data.PilotData.Ship.Pos)

	s.camRecalc()
}
func (s *cosmoScene) camRecalc() {
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.AngleDeg = Data.PilotData.Ship.Ang
	s.cam.Recalc()
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	s.background.Draw(image)
	s.compass.Draw(image)

	for _, co := range s.objects {
		co.Draw(image)
	}
	s.trail.Draw(image)

	if Data.NaviData.ActiveMarker {
		s.naviMarker.SetPos(Data.NaviData.MarkerPos)
		s.naviMarker.Draw(image)
	}

	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	s.ship.Draw(image)

	s.caption.Draw(image)

	s.thrustLevelHUD.Draw(image)
	s.thrustControlHUD.Draw(image)
	s.turnLevelHUD.Draw(image)
	s.turnControlHUD.Draw(image)
}

func (s *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}

func (s *cosmoScene) toWarp() {
	state := Data.State
	state.StateID = STATE_warp
	state.GalaxyID = WARP_Galaxy_ID
	Client.RequestNewState(state.Encode())
}
