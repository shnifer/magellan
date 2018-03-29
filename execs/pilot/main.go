package main

import (
	"time"
	"github.com/hajimehoshi/ebiten"
	"log"
	"fmt"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/Shnifer/magellan/input"
)

const resPath = "res/pilot/"

var last time.Time

func  mainLoop(window *ebiten.Image) error {
	dt:=time.Since(last)
	last = time.Now()

	input.Update()

	window.Clear()

	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %v\ndt = %.2f\n", fps, dt.Seconds())
	if input.Get("forward") {
		msg = msg + "forward!\n"
	}
	msg = msg + fmt.Sprint("Turn = ",input.GetF("turn"))
	ebitenutil.DebugPrint(window, msg)

	return nil
}

func main(){
	startClient()
	input.LoadConf(resPath)

	ebiten.SetFullscreen(DEFVAL.FullScreen)
	ebiten.SetRunnableInBackground(true)
	last=time.Now()
	if err := ebiten.Run(mainLoop, DEFVAL.WinW, DEFVAL.WinH, 1, "PILOT"); err != nil {
		log.Fatal(err)
	}
}

