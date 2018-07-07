package main

import (
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
	"os"
	"runtime/trace"
	"time"
)

const shortWorkMs = 12
const longWorkMs = 22

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

func startSome() {
}

func startALot() {
}

func endSome() {
}

func endALot() {
}

func doSomeWork(win *ebiten.Image) {
	go startSome()
	t := time.After(shortWorkMs * time.Millisecond)
	var s int

loop:
	for {
		select {
		case <-t:
			break loop
		default:
			s = (s + 1) % 100
		}

	}
	win.Fill(color.Black)
	go endSome()
}

func doALotofWork(win *ebiten.Image) {
	go startALot()
	t := time.After(longWorkMs * time.Millisecond)
	var s int

loop:
	for {
		select {
		case <-t:
			break loop
		default:
			s = (s + 1) % 100
		}

	}
	win.Fill(color.White)
	go endALot()
}

func mainLoop2(win *ebiten.Image) error {
	slow := ebiten.IsRunningSlowly()
	ch <- slow
	if slow {
		doSomeWork(win)
	} else {
		doALotofWork(win)
	}
	return nil
}

var ch chan bool

func monitor() {
	t := 0
	f := 0
	tick := time.Tick(time.Second)
	for {
		select {
		case isSlo := <-ch:
			if isSlo {
				t++
			} else {
				f++
			}
		case <-tick:
			log.Printf("%v / %v", f, t+f)
			t = 0
			f = 0
		}
	}
}

func main() {
	f, err := os.Create("dttrace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	trace.Start(f)
	defer trace.Stop()
	ch = make(chan bool, 1)
	go monitor()
	start = time.Now()
	last = start
	tick = time.Tick(time.Second)
	ebiten.SetRunnableInBackground(true)
	ebiten.Run(mainLoop2, 200, 200, 1, "Dt test")
}
