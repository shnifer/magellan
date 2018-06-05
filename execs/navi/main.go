package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/log"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"time"
	"image/color"
	"github.com/Shnifer/magellan/v2"
)

const roleName = "navi"

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

	fps := ebiten.CurrentFPS()
	msg := fmt.Sprintf("FPS: %.0f\nALT-F4 to close\nWASD to control\nQ-E scale\nSPACE - stop\nENTER - reset position", fps)
	fpsText:=graph.NewText(msg,draw.Fonts[draw.Face_list],color.White)
	fpsText.SetPosPivot(graph.ScrP(0.1,0.1),v2.ZV)
	fpsText.Draw(window)

	t := time.Now()
	dt = t.Sub(last).Seconds()
	last = t
	return nil
}

func main() {
	log.Start(time.Duration(DEFVAL.LogLogTimeoutMs)*time.Millisecond,
		time.Duration(DEFVAL.LogRetryMinMs)*time.Millisecond,
		time.Duration(DEFVAL.LogRetryMaxMs)*time.Millisecond,
		DEFVAL.LogIP)

	if DEFVAL.DoProf {
		commons.StartProfile(roleName)
		defer commons.StopProfile(roleName)
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
	input.LoadConf("input_" + roleName + ".json")

	draw.InitFonts()
	draw.InitTexAtlas()
	commons.InitSignatureAtlas()

	createScenes()

	Client.Start()
	ebiten.SetRunnableInBackground(true)
	last = time.Now()
	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "NAVIGATOR"); err != nil {
		log.Log(log.LVL_FATAL, err)
	}
}
