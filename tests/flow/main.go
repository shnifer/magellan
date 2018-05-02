package main

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/graph/flow"
	. "github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"math"
	"time"
)

func update(window *ebiten.Image) error {
	now := time.Now()
	dt := now.Sub(last).Seconds()
	last = now

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		flowN = (flowN + 1) % len(flows)
		fl = flows[flowN].New()
	}

	fl.Update(dt)

	if ebiten.IsRunningSlowly() {
		return nil

	}

	Q := graph.NewDrawQueue()
	Q.Append(fl)
	Q.Run(window)

	return nil
}

var sprite *graph.Sprite
var fl *flow.Flow
var last time.Time

var flows []flow.FlowParams
var flowN int

const screenSize = 600

func main() {
	cam := graph.NewCamera()
	cam.Center = V2{X: screenSize / 2, Y: screenSize / 2}
	cam.Scale = screenSize / 2
	cam.Recalc()

	var err error
	sprite, err = graph.NewSpriteFromFile(
		"res/textures/flame_ani.png", true, 192, 192, 19, cam, false, false)
	if err != nil {
		sprite, err = graph.NewSpriteFromFile(
			"flame_ani.png", true, 192, 192, 19, cam, false, false)
		if err != nil {
			panic(err)
		}
	}
	sprite.SetSize(0.2, 0.2)

	drawer := flow.SpriteDrawerParams{
		Sprite:       sprite,
		DoRandomLine: true,
		FPS:          10,
		CycleType:    graph.Cycle_PingPong,
		Layer:        graph.Z_GAME_OBJECT,
	}.New

	const (
		medLife     = 5
		devLife     = 30
		spawnPeriod = 0.2
	)

	velFs := make(map[string]func(V2) V2)

	velFs["rotation"] = flow.ComposeRadial(flow.ConstC(1), flow.ConstC(0.1))

	sinx := func(x, y float64) float64 {
		return math.Sin(y*5) / 2
	}
	velFs["sinUp"] = flow.ComposeDecart(sinx, flow.ConstC(0.3))

	flows = append(flows, flow.FlowParams{
		SpawnPeriod:    spawnPeriod,
		SpawnPos:       flow.RandomInCirc(1),
		SpawnLife:      flow.NormRand(medLife, devLife),
		SpawnUpdDrawer: drawer,
		VelocityF:      velFs["rotation"],
	})

	flows = append(flows, flow.FlowParams{
		SpawnPeriod:    spawnPeriod,
		SpawnPos:       flow.RandomOnSide(V2{0, -1}, 0.5),
		SpawnLife:      flow.NormRand(medLife, devLife),
		SpawnUpdDrawer: drawer,
		VelocityF:      velFs["sinUp"],
	})

	fl = flows[flowN].New()

	last = time.Now()
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(update, screenSize, screenSize, 1, "test")
}
