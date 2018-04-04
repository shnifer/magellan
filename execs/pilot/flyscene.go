package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/Shnifer/magellan/graph"
	"golang.org/x/image/colornames"
	. "github.com/Shnifer/magellan/commons"
)

type cosmoScene struct{
	ship *graph.Sprite
	caption *graph.Text
	cam *graph.Camera
}

func newCosmoScene() *cosmoScene{
	caption:=graph.NewText("Fly scene", fonts[face_cap], colornames.Aliceblue)
	caption.SetPosPivot(graph.ScrP(0.1,0.1),graph.TopLeft())

	cam:=graph.NewCamera()
	cam.Center = graph.ScrP(0.5, 0.5)
	cam.Recalc()

	ship,err:=graph.NewSpriteFromFile(texPath+"ship.png", ebiten.FilterDefault, 0, 0, cam, false)
	if err!=nil{
		panic(err)
	}

	return &cosmoScene{
		caption: caption,
		ship:ship,
		cam:cam,
	}
}

func (*cosmoScene) Init() {
	defer LogFunc("cosmoScene.Init")()
}

func (scene *cosmoScene) Update(dt float64) {
	defer LogFunc("cosmoScene.Update")()

	switch {
	case ebiten.IsKeyPressed(ebiten.KeyW):
		Data.ship.Pos.Y+=10
	case ebiten.IsKeyPressed(ebiten.KeyS):
		Data.ship.Pos.Y-=10
	}
}

func (scene *cosmoScene) Draw(image *ebiten.Image) {
	defer LogFunc("cosmoScene.Draw")()

	scene.caption.Draw(image)
	scene.ship.SetPos(Data.ship.Pos)
	img, op:= scene.ship.ImageOp()
	image.DrawImage(img,op)
}

func (*cosmoScene) Destroy() {
}