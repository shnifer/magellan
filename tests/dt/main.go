package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"log"
	"sync"
	"time"
)

var mu sync.Mutex
var n int

func ID() int {
	mu.Lock()
	defer mu.Unlock()
	n++
	return n
}

var last time.Time

func mainLoop(window *ebiten.Image) error {
	t := time.Now()
	dt := t.Sub(last).Seconds()
	last = t
	msg := fmt.Sprintf("%0.10f", dt)
	if dt < 0.001 {
		msg = msg + "!"
	}
	if ebiten.IsRunningSlowly() {
		msg = msg + " slow"
	}
	log.Println(msg)
	time.Sleep(time.Microsecond)
	return nil
}

func main() {
	ebiten.Run(mainLoop, 200, 200, 1, "Dt test")
}
