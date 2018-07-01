package main

import (
	"github.com/hajimehoshi/ebiten"
	."github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"time"
	"github.com/Shnifer/magellan/graph"
	"golang.org/x/image/colornames"
	"log"
)

var ships [3]RBData
var last time.Time
var sprite *graph.Sprite
var img *ebiten.Image

const littleT = 0.00001
const maxT = 0.01
const minT = 0.000001

func calcShip(n int,dt float64) {
	ship:=ships[n]
	for dt>0 {
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

		len2:= ship.Pos.LenSqr()
		grav := Gravity(1, len2, 0)
		lt:=littleT
		if dt<lt{
			break
		}
		dt-=lt

		accel:=ship.Pos.Normed().Mul(-grav)

		ship.Vel = ship.Vel.AddMul(accel, lt)
		ship.Pos = ship.Pos.AddMul(ship.Vel, lt)
	}
	ships[n] = ship
}

func run(window *ebiten.Image) error{
	t:=time.Now()
	dt:=t.Sub(last).Seconds()
	last = t

	for n:=0;n<3;n++{
		calcShip(n, dt)
	}

	sprite.SetPos(ships[0].Pos)
	sprite.SetColor(colornames.Red)
	sprite.Draw(img)
	sprite.SetPos(ships[1].Pos)
	sprite.SetColor(colornames.Green)
	sprite.Draw(img)
	sprite.SetPos(ships[2].Pos)
	sprite.SetColor(colornames.Blue)
	sprite.Draw(img)
	window.DrawImage(img,&ebiten.DrawImageOptions{})

	return nil
}

func main(){
	SetGravityConsts(1,1)

	start:=v2.V2{X:0.3, Y:0}

	ships[0].Pos=start
	ships[1].Pos=start
	ships[2].Pos=start

	ships[0].Vel=v2.V2{X:0, Y:0.3}
	ships[1].Vel=v2.V2{X:0.2, Y:0.5}
	ships[2].Vel=v2.V2{X:-0.2, Y:1.5}

	last=time.Now()
	cam:=graph.NewCamera()
	cam.Center = v2.V2{450,450}
	cam.Scale = 1200
	cam.Recalc()

	var err error
	sprite,err=graph.NewSpriteFromFile("res/textures/particle.png",true,0,0,1,cam.Deny())
	if err!=nil{
		panic(err)
	}
	sprite.SetSize(10,10)

	img,_=ebiten.NewImage(900,900,ebiten.FilterDefault)

	log.Println("start")
	ebiten.Run(run, 900,900,1,"Gravity test")
}
