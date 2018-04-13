package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
)

type cosmoScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	objects []*CosmoPoint
	idMap   map[string]*CosmoPoint
}

func newCosmoScene() *cosmoScene {
	caption := graph.NewText("Fly scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := graph.NewSprite(GetAtlasTex("ship"), cam, true)
	ship.SetSize(50, 50)

	return &cosmoScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		objects: make([]*CosmoPoint, 0),
		idMap:   make(map[string]*CosmoPoint),
	}
}

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	s.objects = s.objects[:0]
	s.idMap = make(map[string]*CosmoPoint)

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
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyW):
		Data.PilotData.Ship.Vel.Y += 100
	case ebiten.IsKeyPressed(ebiten.KeyS):
		Data.PilotData.Ship.Vel.Y -= 100
	case ebiten.IsKeyPressed(ebiten.KeyA):
		Data.PilotData.Ship.AngVel += 1
	case ebiten.IsKeyPressed(ebiten.KeyD):
		Data.PilotData.Ship.AngVel -= 1
	}

	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)
	for _, co := range s.objects {
		co.Update(sessionTime)
	}
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
}

func (s *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}
