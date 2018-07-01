package main

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"time"
)

var start time.Time
var last time.Time
var tick <-chan time.Time
var n, m int

func mainLoop(window *ebiten.Image) error {
	now := time.Now()
	waited := now.Sub(last).Seconds() * 1000
	started := now.Sub(start).Seconds() * 1000
	last = now

	//update short time
	time.Sleep(time.Millisecond)

	if !ebiten.IsRunningSlowly() {
		//drawing long time
		//time.Sleep(time.Second / 20)
		m++
	}

	now = time.Now()
	ended := now.Sub(start).Seconds() * 1000
	worked := now.Sub(last).Seconds() * 1000
	last = now
	n++
	log.Printf("waited %10.3f / started at %10.3f / worked for %10.3f / ended at %10.3f / slow:%v", waited, started, worked, ended, ebiten.IsRunningSlowly())
	if waited+worked < 0 {
		add := 10 - waited - worked
		time.Sleep(time.Duration(add) * time.Millisecond)
	}
	select {
	case <-tick:
		log.Printf("drawed %v/%v", m, n)
		n = 0
		m = 0
	default:
	}
	return nil
}

func main() {
	start = time.Now()
	last = start
	tick = time.Tick(time.Second)
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(mainLoop, 200, 200, 1, "Dt test")
}
