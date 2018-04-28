package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
)

type engiScene struct {
	ship       *graph.Sprite
	caption    *graph.Text
	background *graph.Sprite
}

func newEngiScene() *engiScene {
	caption := graph.NewText("Engi scene", Fonts[Face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1, 0.1), graph.TopLeft())

	back := NewAtlasSpriteHUD("engibackground")
	back.SetSize(float64(WinW), float64(WinH))

	return &engiScene{
		caption:    caption,
		background: back,
	}
}

func (*engiScene) Init() {
	defer LogFunc("engiScene.Init")()
}

func (scene *engiScene) Update(dt float64) {
	defer LogFunc("engiScene.Update")()
}

func (scene *engiScene) Draw(image *ebiten.Image) {
	defer LogFunc("engiScene.Draw")()

	scene.background.Draw(image)
	scene.caption.Draw(image)
}

func (scene *engiScene) OnCommand(command string) {
}

func (*engiScene) Destroy() {
}
