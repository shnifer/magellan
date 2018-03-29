package main

import (
	"fmt"
	"github.com/Shnifer/magellan/input"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"time"
	"github.com/Shnifer/magellan/scene"
	"github.com/Shnifer/magellan/graph"
)

const DEBUG = true
const resPath = "res/pilot/"
const fontPath  = "res/fonts/"

var (
	WinW int
	WinH int
)

var last time.Time
func mainLoop(window *ebiten.Image) error {
	dt := time.Since(last).Seconds()
	last = time.Now()

	input.Update()

	window.Clear()
	Scenes.Update(dt)
	if !ebiten.IsRunningSlowly() {
		Scenes.Draw(window)
	}

	if DEBUG {
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %v\ndt = %.2f\n", fps, dt)
		if input.Get("forward") {
			msg = msg + "forward!\n"
		}

		if ebiten.IsRunningSlowly() {
			msg = msg + "is running SLOWLY!\n"
		}
		msg = msg + fmt.Sprint("Turn = ", input.GetF("turn"))
		ebitenutil.DebugPrint(window, msg)
	}

	return nil
}

var Scenes *scene.Manager
func main() {
	WinW = DEFVAL.WinW
	WinH = DEFVAL.WinH

	startClient()
	input.LoadConf(resPath)
	Scenes = scene.NewManager()

	face,err:=graph.GetFace(fontPath+"phantom.ttf",20)
	if err!=nil{
		panic(err)
	}
	pauseScene := scene.NewPauseScene(face, Client.PauseReason)
	Scenes.Install("pause", pauseScene)
	Scenes.Activate("pause")

	ebiten.SetFullscreen(DEFVAL.FullScreen)
	ebiten.SetRunnableInBackground(true)
	last = time.Now()
	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "PILOT"); err != nil {
		log.Fatal(err)
	}
}
