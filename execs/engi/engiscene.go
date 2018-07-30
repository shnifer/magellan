package main

import (
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type engiScene struct {
	background *graph.Sprite

	systemsMonitor *systemsMonitor

	q *graph.DrawQueue
}

func newEngiScene() *engiScene {
	back := NewAtlasSpriteHUD(EngiBackgroundAN)
	back.SetSize(float64(WinW), float64(WinH))
	back.SetPivot(graph.TopLeft())

	return &engiScene{
		background:     back,
		systemsMonitor: newSystemsMonitor(),
		q:              graph.NewDrawQueue(),
	}
}

func (*engiScene) Init() {
	defer LogFunc("engiScene.Init")()
}

func (scene *engiScene) Update(dt float64) {
	defer LogFunc("engiScene.Update")()

	//emissions := CalculateEmissions(Data.Galaxy, Data.PilotData.Ship.Pos)

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
	}
}

func (scene *engiScene) Draw(image *ebiten.Image) {
	defer LogFunc("engiScene.Draw")()
	Q := scene.q
	Q.Clear()

	Q.Append(scene.systemsMonitor)

	Q.Run(image)
}

func (scene *engiScene) OnCommand(command string) {
}

func (*engiScene) Destroy() {
}
