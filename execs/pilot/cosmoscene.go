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
	objects []*cosmoPoint
}

type cosmoPoint struct {
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
	}
}

func (*cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()

}

func (scene *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()

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

func (scene *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	scene.caption.Draw(image)
	scene.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	scene.ship.Draw(image)
}

func (scene *cosmoScene) OnCommand(command string) {
}

func (*cosmoScene) Destroy() {
}
