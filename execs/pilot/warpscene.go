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
	"log"
	"math"
	"sort"
	"time"
)

const warpingOutTime = 3
const warpingOutAnnounceP = 0.2

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
	fuelConsumed float64

	warpingOutT float64
	at          *AnnounceText
	atTime      float64
	alreadyWarp bool

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

	at := NewAnnounceText(graph.ScrP(0.5, 0.3), graph.Center(),
		Fonts[Face_cap], graph.Z_STAT_HUD)

	res := warpScene{
		ship:        ship,
		cam:         cam,
		objects:     make(map[string]*CosmoPoint),
		objIDs:      make([]string, 0),
		sonarSector: sonarSector,
		q:           graph.NewDrawQueue(),
		hud:         newWarpSceneHUD(cam),
		at:          at,
	}

	return &res
}

func (s *warpScene) Init() {
	defer LogFunc("warpScene.Init")()

	s.objects = make(map[string]*CosmoPoint)
	s.objIDs = make([]string, 0)
	s.thrustLevel = 0.01
	s.maneurLevel = 0
	s.distTravaled = 0
	s.timerStart = time.Now()
	s.fuelConsumed = 0
	s.warpingOutT = 0
	s.alreadyWarp = false
	s.cam.Scale = 1

	stateData := Data.GetStateData()

	predictorSprite := NewAtlasSprite(PredictorAN, s.cam.Deny())
	predictorSprite.SetSize(20, 20)
	opts := WarpPredictorOpts{
		Cam:      s.cam,
		Sprite:   predictorSprite,
		Clr:      colornames.Palevioletred,
		Layer:    graph.Z_ABOVE_OBJECT + 1,
		Galaxy:   stateData.Galaxy,
		UpdT:     DEFVAL.WarpPredictorUpdT,
		NumInSec: DEFVAL.WarpPredictorNumInSec,
		TrackLen: DEFVAL.WarpPredictorTrackLen,
		DrawMaxP: DEFVAL.WarpPredictorDrawMaxP,
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

	Data.PilotData.FlightTime+=dt

	ppos := Data.PilotData.Ship.Pos
	UpdateWarpAndShip(Data, dt, DEFVAL.DT, DEFVAL.WarpGravPowN)
	s.distTravaled += Data.PilotData.Ship.Pos.Sub(ppos).Len()
	s.fuelConsumed += Data.SP.Warp_engine.Consumption *
		(s.thrustLevel + math.Abs(s.maneurLevel)) * dt

	Data.PilotData.WarpPos = Data.PilotData.Ship.Pos
	s.at.Update(dt)

	s.warpChecks()

	for _, co := range s.objects {
		co.Update(dt)
	}
	s.updateShipControl(dt)

	if DEFVAL.DebugControl {
		s.debugControl()
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

	Q.Append(s.at)

	Q.Append(s.predictor)
	Q.Append(s.hud)

	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %.0f\nTravel: %.0f\nFuel used: %.0f\nTimer: %.1f\nDraws: %v",
		fps, s.distTravaled, s.fuelConsumed, time.Since(s.timerStart).Seconds(), Q.Len())
	fpsText = graph.NewText(msg, Fonts[Face_list], colornames.Cyan)
	fpsText.SetPosPivot(graph.ScrP(0.1, 0.1), v2.ZV)
	Q.Add(fpsText, graph.Z_STAT_HUD+10)

	Q.Run(image)
}

func (s *warpScene) warpedOut() {
	if s.alreadyWarp {
		return
	}
	systemID := ""
	for _, gp := range Data.Galaxy.Ordered {
		if Data.PilotData.Ship.Pos.Sub(gp.Pos).LenSqr() <
			gp.WarpYellowOutDist*gp.WarpYellowOutDist {
			systemID = gp.ID
			break
		}
	}
	if systemID == "" {
		log.Println("warp to zero system")
		s.toCosmo(ZERO_Galaxy_ID)
		return
	}
	ship := Data.PilotData.Ship.Pos
	sys := Data.Galaxy.Points[systemID]
	dist := ship.Sub(sys.Pos).Len()
	if dist <= sys.WarpRedOutDist {
		log.Println("hard damage")
	} else if dist < sys.WarpGreenInDist || dist > sys.WarpGreenOutDist {
		log.Println("medium damage")
	}
	s.toCosmo(systemID)
}

func (s *warpScene) toCosmo(systemID string) {
	s.alreadyWarp = true
	state := Data.State
	state.StateID = STATE_cosmo
	state.GalaxyID = systemID
	Client.RequestNewState(state.Encode(), false)
}

func (*warpScene) Destroy() {
}

func (s *warpScene) debugControl() {
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
		s.fuelConsumed = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		s.toCosmo("solar")
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= (1 + dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= (1 + dt)
	}
	s.cam.Scale = Clamp(s.cam.Scale, 0.0001, 10000)
}

func (s *warpScene) warpChecks() {
	if Data.PilotData.Distortion == 0 {
		s.warpingOutT += dt
		s.atTime -= dt

		if s.atTime <= 0 {
			s.atTime = warpingOutAnnounceP

			msg := fmt.Sprintf("warping out in %1.1f", warpingOutTime-s.warpingOutT)
			clr := colornames.Red
			if s.warpingOutT < warpingOutTime {
				clr = colornames.Yellow
			}
			s.at.AddMsg(msg, clr, warpingOutAnnounceP)
		}
	} else {
		s.warpingOutT = 0
		s.atTime = 0
	}
	if s.warpingOutT > warpingOutTime {
		s.warpedOut()
	}
	for _, gp := range Data.Galaxy.Ordered {
		if Data.PilotData.Ship.Pos.Sub(gp.Pos).LenSqr() <
			gp.WarpRedOutDist*gp.WarpRedOutDist {
			s.warpedOut()
		}
	}
}
