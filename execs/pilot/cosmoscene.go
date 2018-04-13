package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/scene"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"math"
	"time"
)

type cosmoScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera

	objects []*cosmoPoint
	idMap   map[string]*cosmoPoint
}

type cosmoPoint struct {
	sprite *graph.Sprite

	pos  v2.V2
	size float64

	parent   *cosmoPoint
	orbit    float64
	angVel   float64
	angPhase float64
}

var texAssoc map[string]string

func getTexByType(typeName string) graph.Tex {
	//TODO: make texture atlas structure
	if texAssoc == nil {
		texAssoc = make(map[string]string)
		texAssoc["star"] = ""
		texAssoc["planet"] = ""
	}
	fn, ok := texAssoc[typeName]
	if !ok {
		panic("getTexByType: unknown type " + typeName)
	}

	tex, err := graph.GetTex(texPath+fn, ebiten.FilterDefault, 0, 0)
	if err != nil {
		panic(err)
	}
	return tex
}

func newCosmoPoint(pd GalaxyPoint, cam *graph.Camera) *cosmoPoint {
	tex := getTexByType(pd.Type)
	sprite := graph.NewSprite(tex, cam, false)
	return &cosmoPoint{
		sprite:   sprite,
		pos:      pd.Pos,
		size:     pd.Size,
		orbit:    pd.Orbit,
		angVel:   360 / pd.Period,
		angPhase: pd.DegStart,
	}
}

func newCosmoScene() *cosmoScene {
	caption := graph.NewText("Fly scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship, err := graph.NewSpriteFromFile(texPath+"ship.png", ebiten.FilterDefault, 0, 0, cam, false)
	if err != nil {
		panic(err)
	}

	return &cosmoScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
		objects: make([]*cosmoPoint, 0),
		idMap:   make(map[string]*cosmoPoint),
	}
}

func (s *cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

	s.objects = s.objects[:0]
	s.idMap = make(map[string]*cosmoPoint)

	stateData := Data.GetStateData()

	for _, pd := range stateData.Galaxy.Points {
		cosmoPoint := newCosmoPoint(pd, s.cam)
		s.objects = append(s.objects, cosmoPoint)

		if pd.ID != "" {
			s.idMap[pd.ID] = cosmoPoint
		}
		if pd.ParentID != "" {
			cosmoPoint.parent = s.idMap[pd.ParentID]
		}
	}
}

func (s *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()
	Data.PilotData.SessionTime = Data.PilotData.SessionTime.Add(time.Second * time.Duration(dt))

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
}

func (s *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	s.caption.Draw(image)
	s.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	s.ship.Draw(image)
}

func (s *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}
