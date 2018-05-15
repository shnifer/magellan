package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"time"
)

const resPath = "res/navi/"
const texPath = "res/textures/"

var (
	WinW int
	WinH int
)

var last time.Time
var Data commons.TData
var dt float64

func mainLoop(window *ebiten.Image) error {
	input.Update()

	Data.Update(DEFVAL.Role)

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	if commons.LOG_LEVEL <= commons.LVL_ERROR {
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %v\ndt = %.2f\n", fps, dt)
		if ebiten.IsRunningSlowly() {
			msg = msg + "is running SLOWLY!\n"
		}
		ebitenutil.DebugPrint(window, msg)
	}
	t := time.Now()
	dt = t.Sub(last).Seconds()
	last = t
	return nil
}

func main() {
	if DEFVAL.DoProf {
		commons.StartProfile(DEFVAL.CpuProfFileName)
		defer commons.StopProfile(DEFVAL.CpuProfFileName, DEFVAL.MemProfFileName)
	}

	if DEFVAL.FullScreen {
		ebiten.SetFullscreen(true)
		WinW, WinH = ebiten.MonitorSize()
	} else {
		WinW = DEFVAL.WinW
		WinH = DEFVAL.WinH
	}

	graph.SetScreenSize(WinW, WinH)

	Data = commons.NewData()

	initClient()
	input.LoadConf(resPath)

	draw.InitTexAtlas(texPath)

	createScenes()

	Client.Start()
	ebiten.SetRunnableInBackground(true)
	last = time.Now()
	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "NAVIGATOR"); err != nil {
		log.Fatal(err)
	}
}
