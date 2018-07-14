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
	"sort"
	"time"
)

type warpScene struct {
	ship *graph.Sprite
	cam  *graph.Camera

	objects map[string]*CosmoPoint
	objIDs  []string

	//sonar
	sonarSector *graph.Sector

	//control
	thrustLevel float64
	maneurLevel float64

	//setParams eachUpdate
	gravityAcc    v2.V2
	gravityReport []v2.V2

	predictor *WarpPredictor

	distTravaled float64
	timerStart   time.Time

	hud warpSceneHUD
	q   *graph.DrawQueue
}

func newWarpScene() *warpScene {
	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Scale = 100
	cam.Recalc()

	ship := NewAtlasSprite(ShipAN, cam.FixS())
	ship.SetSize(50, 50)

	sonarSector := graph.NewSector(cam.Phys())
	sonarSector.SetColor(colornames.Forestgreen)

	res := warpScene{
		ship:        ship,
		cam:         cam,
		objects:     make(map[string]*CosmoPoint),
		objIDs:      make([]string, 0),
		sonarSector: sonarSector,
		q:           graph.NewDrawQueue(),
		hud:         newWarpSceneHUD(cam),
	}

	return &res
}

func (s *warpScene) Init() {
	defer LogFunc("warpScene.Init")()

	s.objects = make(map[string]*CosmoPoint)
	s.thrustLevel = 0
	s.maneurLevel = 0
	s.distTravaled = 0
	s.timerStart = time.Now()

	stateData := Data.GetStateData()

	predictorSprite := NewAtlasSprite(PredictorAN, s.cam.Deny())
	predictorSprite.SetSize(20, 20)
	opts := WarpPredictorOpts{
		Cam:      s.cam,
		Sprite:   predictorSprite,
		Clr:      colornames.Palevioletred,
		Layer:    graph.Z_ABOVE_OBJECT + 1,
		Galaxy:   stateData.Galaxy,
		UpdT:     0.1,
		NumInSec: 10,
		TrackLen: 120,
		DrawMaxP: 30,
		PowN:     DEFVAL.WarpGravPowN,
	}

	s.predictor = NewWarpPredictor(opts)

	for id, pd := range stateData.Galaxy.Points {
		if pd.IsVirtual {
			continue
		}
		cosmoPoint := NewCosmoPoint(pd, s.cam.Phys())
		s.objects[id] = cosmoPoint
	}
}

func (s *warpScene) Update(dt float64) {
	defer LogFunc("warpScene.Update")()

	ppos := Data.PilotData.Ship.Pos
	UpdateWarpAndShip(Data, dt, DEFVAL.DT)
	s.distTravaled += Data.PilotData.Ship.Pos.Sub(ppos).Len()
	UpdateWarpAndShip(Data, dt, DEFVAL.DT, DEFVAL.WarpGravPowN)

	for _, co := range s.objects {
		co.Update(dt)
	}
	s.updateShipControl(dt)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		Data.PilotData.Ship.Pos = v2.V2{}
		Data.PilotData.Ship.Vel = v2.V2{}
		Data.PilotData.Ship.AngVel = 0
		Data.PilotData.Distortion = DEFVAL.MinDistortion
		Data.PilotData.Dir = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.timerStart = time.Now()
		s.distTravaled = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
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
			s.toCosmo(systemID)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= (1 + dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= (1 + dt)
	}

	s.updateHUD()

	//for display draw calls only
	s.gravityAcc, s.gravityReport = SumGravityAccWithReport(Data.PilotData.Ship.Pos, Data.Galaxy,
		0.02)

	s.sonarSector.SetCenter(Data.PilotData.Ship.Pos)
	s.sonarSector.SetRadius(Data.NaviData.SonarRange)
	s.sonarSector.SetAngles(
		Data.NaviData.SonarDir-Data.NaviData.SonarWide/2,
		Data.NaviData.SonarDir+Data.NaviData.SonarWide/2)

	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	s.predictor.SetPosDistDir(
		Data.PilotData.Ship.Pos, Data.PilotData.Distortion, Data.PilotData.Dir)
	s.camRecalc()
}
func (s *warpScene) camRecalc() {
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.AngleDeg = Data.PilotData.Ship.Ang
	s.cam.Recalc()
}

func (s *warpScene) Draw(image *ebiten.Image) {
	defer LogFunc("warpScene.Draw")()

	Q := s.q
	Q.Clear()

	Q.Add(s.sonarSector, graph.Z_UNDER_OBJECT)

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

	Q.Add(s.ship, graph.Z_GAME_OBJECT)

	Q.Append(s.predictor)
	Q.Append(s.hud)

	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %.0f\nTravel: %.0f\nTimer: %.1f\nDraws: %v", fps,
		s.distTravaled, time.Since(s.timerStart).Seconds(), Q.Len())
	fpsText = graph.NewText(msg, Fonts[Face_list], colornames.Cyan)
	fpsText.SetPosPivot(graph.ScrP(0.1, 0.1), v2.ZV)
	Q.Add(fpsText, graph.Z_STAT_HUD+10)

	Q.Run(image)
}

func (s *warpScene) toCosmo(systemID string) {
	state := Data.State
	state.StateID = STATE_cosmo
	state.GalaxyID = systemID
	Client.RequestNewState(state.Encode(), false)
}

func (*warpScene) Destroy() {
}
