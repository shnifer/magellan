package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
)

type engiScene struct {
	ship    *graph.Sprite
	caption *graph.Text
	cam     *graph.Camera
}

func newCosmoScene() *engiScene {
	caption := graph.NewText("Engi scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	cam := graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship := graph.NewSprite(GetAtlasTex("ship"), cam, false)

	return &engiScene{
		caption: caption,
		ship:    ship,
		cam:     cam,
	}
}

func (*engiScene) Init() {
	defer LogFunc("engiScene.Init")()
}

func (scene *engiScene) Update(dt float64) {
	defer LogFunc("engiScene.Update")()

	//PilotData Rigid Body emulation
	Data.PilotData.Ship = Data.PilotData.Ship.Extrapolate(dt)
}

func (scene *engiScene) Draw(image *ebiten.Image) {
	defer LogFunc("engiScene.Draw")()

	scene.caption.Draw(image)
	scene.ship.SetPosAng(Data.PilotData.Ship.Pos, Data.PilotData.Ship.Ang)
	img, op := scene.ship.ImageOp()
	image.DrawImage(img, op)
}

func (scene *engiScene) OnCommand(command string) {
}

func (*engiScene) Destroy() {
}
