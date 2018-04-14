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

type cosmoScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	objects []*CosmoPoint
	idMap   map[string]*CosmoPoint

	//control
	thrustLevel      float64
	maneurLevel      float64
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

	ship := graph.NewSprite(GetAtlasTex("ship"), cam, true)
	ship.SetSize(50, 50)

	res := cosmoScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		objects: make([]*CosmoPoint, 0),
		idMap:   make(map[string]*CosmoPoint),
	}

	arrowTex := GetAtlasTex("arrow")
	res.thrustLevelHUD = graph.NewSprite(arrowTex, nil, false)
	res.thrustLevelHUD.SetSize(50, 50)
	res.thrustLevelHUD.SetAng(180)
	res.thrustLevelHUD.SetAlpha(0.7)
	res.thrustControlHUD = graph.NewSprite(arrowTex, nil, false)
	res.thrustControlHUD.SetSize(50, 50)
	res.thrustControlHUD.SetAlpha(0.5)
	res.turnLevelHUD = graph.NewSprite(arrowTex, nil, false)
	res.turnLevelHUD.SetSize(50, 50)
	res.turnLevelHUD.SetAng(-90)
	res.turnLevelHUD.SetAlpha(0.7)
	res.turnControlHUD = graph.NewSprite(arrowTex, nil, false)
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
		Data.PilotData.Ship.Pos = v2.V2{100, 100}
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

	s.camRecalc()
}
func (s *cosmoScene) camRecalc() {
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.AngleDeg = Data.PilotData.Ship.Ang
	s.cam.Recalc()
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	s.caption.Draw(image)

	for _, co := range s.objects {
		co.Draw(image)
	}

	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	s.ship.Draw(image)

	s.thrustLevelHUD.Draw(image)
	s.thrustControlHUD.Draw(image)
	s.turnLevelHUD.Draw(image)
	s.turnControlHUD.Draw(image)
}

func (s *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}
