package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
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

	objects map[string]*CosmoPoint

	lastServerID int
	otherShips   map[string]*OtherShip

	//control
	thrustLevel float64
	maneurLevel float64

	//trail
	trailT float64
	trail  *graph.FadingArray

	hud cosmoSceneHUD

	showPredictor   bool
	predictorZero   *TrackPredictor
	predictorThrust *TrackPredictor
}

func newCosmoScene() *cosmoScene {
	caption := graph.NewText("Fly scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite("ship", cam.FixS())
	ship.SetSize(50, 50)

	marker := NewAtlasSprite("marker", cam.Deny())
	marker.SetPivot(graph.MidBottom())

	predictorSprite := NewAtlasSprite("trail", cam.Deny())
	predictorSprite.SetSize(20, 20)
	predictorThrust := NewTrackPredictor(cam, predictorSprite, &Data, Track_CurrentThrust, colornames.Palevioletred, graph.Z_ABOVE_OBJECT+1)

	predictor2Sprite := NewAtlasSprite("trail", cam.Deny())
	predictor2Sprite.SetSize(15, 15)
	predictor2Sprite.SetColor(colornames.Darkgray)

	predictorZero := NewTrackPredictor(cam, predictor2Sprite, &Data, Track_ZeroThrust, colornames.Cadetblue, graph.Z_ABOVE_OBJECT)

	hud := newCosmoSceneHUD(cam)

	res := cosmoScene{
		caption:         caption,
		ship:            ship,
		cam:             cam,
		naviMarker:      marker,
		hud:             hud,
		objects:         make(map[string]*CosmoPoint),
		otherShips:      make(map[string]*OtherShip),
		showPredictor:   true,
		predictorThrust: predictorThrust,
		predictorZero:   predictorZero,
	}

	res.trail = graph.NewFadingArray(GetAtlasTex("trail"), trailLifeTime/trailPeriod, cam.Deny())

	return &res
}

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	s.objects = make(map[string]*CosmoPoint)
	s.otherShips = make(map[string]*OtherShip)
	s.thrustLevel = 0
	s.maneurLevel = 0
	s.trailT = 0
	s.lastServerID = 0
	s.trail.Clear()

	stateData := Data.GetStateData()

	for _, pd := range stateData.Galaxy.Ordered {
		cosmoPoint := NewCosmoPoint(pd, s.cam.Phys())
		s.objects[pd.ID] = cosmoPoint
	}

	for _, b := range stateData.Buildings {
		s.addBuilding(b)
	}
}

func (s *cosmoScene) addBuilding(b Building) {

}

func (s *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()

	if Data.ServerData.MsgID != s.lastServerID {
		s.lastServerID = Data.ServerData.MsgID

		for _, otherData := range Data.ServerData.OtherShips {
			otherShip, ok := s.otherShips[otherData.Id]
			if !ok {
				otherShip = NewOtherShip(s.cam.FixS(), otherData.Name, float64(DEFVAL.OtherShipElastic)/1000)
				s.otherShips[otherData.Id] = otherShip
			}
			otherShip.SetRB(otherData.Ship)
		}

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

	for id := range s.otherShips {
		s.otherShips[id].Update(dt)
	}

	Data.PilotData.SessionTime += dt
	sessionTime := Data.PilotData.SessionTime
	Data.Galaxy.Update(sessionTime)

	for id, co := range s.objects {
		if gp, ok := Data.Galaxy.Points[id]; ok {
			s.objects[id].Pos = gp.Pos
		}
		co.Update(dt)
	}

	s.updateShipControl(dt)
	s.procShipGravity(dt)
	s.procEmissions(dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
		Data.PilotData.Ship.Pos = Data.Galaxy.Points["magellan"].Pos
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		s.toWarp()
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		s.showPredictor = !s.showPredictor
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		AddBeacon("just a test beacon")
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= 1 + dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= 1 + dt
	}

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
	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)

	if s.thrustLevel > 0 {
		Data.PilotData.HeatProduction = Data.SP.Thrust_heat_prod * s.thrustLevel
	} else {
		Data.PilotData.HeatProduction = 0
	}
	s.UpdateHUD()
	s.camRecalc()
}
func (s *cosmoScene) camRecalc() {
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.AngleDeg = Data.PilotData.Ship.Ang
	s.cam.Recalc()
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	Q := graph.NewDrawQueue()

	Q.Append(s.hud)

	for _, co := range s.objects {
		Q.Append(co)
	}
	Q.Add(s.trail, graph.Z_UNDER_OBJECT)

	if Data.NaviData.ActiveMarker {
		s.naviMarker.SetPos(Data.NaviData.MarkerPos)
		Q.Add(s.naviMarker, graph.Z_ABOVE_OBJECT)
	}

	Q.Add(s.ship, graph.Z_HUD)

	//Q.Add(s.caption, graph.Z_STAT_HUD)

	for _, os := range s.otherShips {
		Q.Append(os)
	}

	if s.showPredictor {
		Q.Append(s.predictorThrust)
		Q.Append(s.predictorZero)
	}

	Q.Run(image)
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
