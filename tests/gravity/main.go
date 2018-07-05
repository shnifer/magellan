package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"time"
)

var lastP [3]v2.V2
var colors [3]color.Color
var ships [3]RBData
var last time.Time
var sprite *graph.Sprite
var img *ebiten.Image
var cam *graph.Camera
var startT time.Time
var crossed bool

const littleT = 0.00001
const maxT = 0.01
const minT = 0.000001

func calcShip(n int, dt float64) {
	ship := ships[n]
	for dt > 0 {
		/*
			grav := Gravity(1, len2, 0)
			gravDev:=Gravity(1, (ln+0.00001)*(ln+0.00001), 0) - grav
			gravDev/=-0.00001
			gravDev/=grav

			moveMax:=0.001/gravDev
			timeMax:=moveMax/ship.Vel.Len()

			lt=timeMax
			if lt>maxT{
				lt=maxT
			} else if lt<minT {
				lt=minT
			}
		*/

		len2 := ship.Pos.LenSqr()
		grav := Gravity(1, len2, 0.0)
		lt := littleT
		if dt < lt {
			break
		}
		dt -= lt

		accel := ship.Pos.Normed().Mul(-grav)

		ship.Vel = ship.Vel.AddMul(accel, lt)
		ship.Pos = ship.Pos.AddMul(ship.Vel, lt)
	}
	ships[n] = ship
}

func run(window *ebiten.Image) error {
	t := time.Now()
	dt := t.Sub(last).Seconds()
	last = t

	for n := 0; n < 3; n++ {
		calcShip(n, dt)
	}

	for n := 0; n < 3; n++ {
		graph.Line(cam, lastP[n], ships[n].Pos, colors[n]).Draw(img)
		sprite.SetPos(ships[n].Pos)
		lastP[n] = ships[n].Pos
		sprite.SetColor(colors[n])
		sprite.Draw(img)
	}

	window.DrawImage(img, &ebiten.DrawImageOptions{})

	if !crossed && ships[0].Pos.Y < 0 {
		crossed = true
	}
	if crossed && ships[0].Pos.Y > 0 {
		//	log.Println(time.Now().Sub(startT).Seconds()*1000)
	}

	return nil
}

func main() {
	SetGravityConsts(1, 1)

	start := v2.V2{X: 1, Y: 0}

	for n := 0; n < 3; n++ {
		ships[n].Pos = start
		lastP[n] = start
	}

	ships[0].Vel = v2.V2{X: 0, Y: 0.3}
	ships[1].Vel = v2.V2{X: 0.5, Y: 0.5}
	ships[2].Vel = v2.V2{X: -0.3, Y: 1}

	last = time.Now()
	cam = graph.NewCamera()
	cam.Center = v2.V2{500, 500}
	cam.Scale = 400
	cam.Recalc()

	colors[0] = colornames.Red
	colors[1] = colornames.Green
	colors[2] = colornames.Blue

	var err error
	sprite, err = graph.NewSpriteFromFile("res/textures/particle.png", true, 0, 0, 1, cam.Deny())
	if err != nil {
		panic(err)
	}
	sprite.SetSize(10, 10)

	img, _ = ebiten.NewImage(1000, 1000, ebiten.FilterDefault)

	circle := graph.NewSprite(graph.CircleTex(), cam.Phys())
	circle.SetSize(0.03, 0.03)
	circle.Draw(img)

	log.Println("start")
	startT = time.Now()
	ebiten.Run(run, 1000, 1000, 1, "Gravity test")
}
