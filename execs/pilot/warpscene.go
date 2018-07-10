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
	"sort"
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

	stateData := Data.GetStateData()

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

	UpdateWarpAndShip(Data, dt, DEFVAL.DT)

	for _, co := range s.objects {
		co.Update(dt)
	}
	s.updateShipControl(dt)

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

	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)

	//for display draw calls only
	s.gravityAcc, s.gravityReport = SumGravityAccWithReport(Data.PilotData.Ship.Pos, Data.Galaxy,
		0.02)

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

	Q.Append(s.hud)

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
