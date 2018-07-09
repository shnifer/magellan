package main

import (
	"github.com/hajimehoshi/ebiten"
	"time"
	"image/color"
	"log"
	"os"
	"runtime/trace"
)

const shortWorkMs = 5
const longWorkMs = 22

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
/*	t := time.After(shortWorkMs * time.Millisecond)
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
*/  time.Sleep(shortWorkMs * time.Millisecond)
	go endSome()
}

func doALotofWork(win *ebiten.Image) {
	go startALot()
	/*for {
	t := time.After(longWorkMs * time.Millisecond)
	var s int
loop:
		select {
		case <-t:
			break loop
		default:
			s = (s + 1) % 100
		}

	}*/
	time.Sleep(longWorkMs * time.Millisecond)
	win.Fill(color.White)
	go endALot()
}

func mainLoop(win *ebiten.Image) error {
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

func fpsMonitor() {
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
			log.Printf("Draw/Total: %v / %v", f, t+f)
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

/*	font,err:=truetype.Parse(fonts.ArcadeN_ttf)
	if err !=nil{
		panic(err)
	}
	face = truetype.NewFace(font, &truetype.Options{Size:20})
*/
	ch = make(chan bool, 1)
	go fpsMonitor()
	ebiten.SetRunnableInBackground(false)
	ebiten.SetFullscreen(true)
	ebiten.Run(mainLoop, 200, 200, 1, "Dt test")
}
