package main

import (
	"fmt"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
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

	//setParams eachUpdate
	gravityAcc    v2.V2
	gravityReport []v2.V2

	q *graph.DrawQueue
}

func newWarpScene() *warpScene {
	caption := graph.NewText("Warp scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := NewAtlasSprite(ShipAN, cam.FixS())
	ship.SetSize(50, 50)

	sonarSector := graph.NewSector(cam.Phys())
	sonarSector.SetColor(colornames.Forestgreen)

	res := warpScene{
		caption:     caption,
		ship:        ship,
		cam:         cam,
		objects:     make(map[string]*CosmoPoint),
		sonarSector: sonarSector,
		q:           graph.NewDrawQueue(),
	}

	res.trail = graph.NewFadingArray(GetAtlasTex(TrailAN), trailLifeTime/trailPeriod,
		cam.Deny())

	arrowTex := GetAtlasTex(ThrustArrowAN)
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

	s.drawScale(Q)
	s.drawGravity(Q)

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

func (s *warpScene) drawScale(Q *graph.DrawQueue) {
	//Scale factor hud
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

	p := func(i int) v2.V2 {
		return s.cam.Center.AddMul(v2.InDir(float64(360/32)*float64(i)), circleRadPx)
	}
	for i := 0; i <= 32; i++ {
		graph.LineScr(Q, p(i), p(i+1), colornames.Oldlace, graph.Z_STAT_HUD+10)
	}

	msg = fmt.Sprintf("circle radius: %f", physRad)
	physRadText := graph.NewText(msg, Fonts[Face_mono], colornames.Oldlace)
	physRadText.SetPosPivot(graph.ScrP(0.5, 0.4), graph.TopMiddle())
	Q.Add(physRadText, graph.Z_STAT_HUD+10)
}

func (s *warpScene) drawGravity(Q *graph.DrawQueue) {
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
