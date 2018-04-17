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
	caption := graph.NewText("Navi scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := graph.NewSprite(GetAtlasTex("ship"), cam, false)
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

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		s.cam.Scale *= (1 + dt)
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		s.cam.Scale /= (1 + dt)
	}

	//PilotData Rigid Body emulation
	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)
	s.cam.Pos = Data.PilotData.Ship.Pos
	s.cam.Recalc()
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	for _, co := range s.objects {
		co.Draw(image)
	}

	s.caption.Draw(image)
	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	img, op := s.ship.ImageOp()
	image.DrawImage(img, op)
}

func (s *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}
