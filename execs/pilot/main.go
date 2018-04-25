package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/hajimehoshi/ebiten"
	"log"
	"time"
)

const resPath = "res/pilot/"
const texPath = "res/textures/"

var (
	WinW int
	WinH int
)

var last time.Time
var Data commons.TData

var ticker <-chan time.Time

func mainLoop(window *ebiten.Image) error {
	t := time.Now()
	dt := t.Sub(last).Seconds()
	last = t

	input.Update()

	Data.Update(DEFVAL.Role)

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	select {
	case <-ticker:
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %v\ndt = %.2f\n", fps, dt)
		commons.Log(commons.LVL_ERROR, msg)
	default:
	}

	return nil
}

func main() {
	if DEFVAL.DoProf {
		startProfile()
		defer stopProfile()
	}

	WinW = DEFVAL.WinW
	WinH = DEFVAL.WinH

	graph.SetScreenSize(WinW, WinH)

	initClient()
	input.LoadConf(resPath)

	Data = commons.NewData()

	commons.InitTexAtlas(texPath)
	commons.SetGravityConsts(DEFVAL.GravityConst, DEFVAL.WarpGravityConst)

	createScenes()

	Client.Start()
	ebiten.SetFullscreen(DEFVAL.FullScreen)
	ebiten.SetRunnableInBackground(true)

	ticker = time.Tick(time.Second)
	last = time.Now()

	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "PILOT"); err != nil {
		log.Fatal(err)
	}
}
