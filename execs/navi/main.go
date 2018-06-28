package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"time"
)

const roleName = "navi"

var (
	WinW int
	WinH int
)

var last time.Time
var Data commons.TData
var dt float64

var fpsText *graph.Text
var showFps <-chan time.Time

func mainLoop(window *ebiten.Image) error {
	input.Update()

	Data.Update(DEFVAL.Role)

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	select {
	case <-showFps:
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %.0f\n", fps)
		fpsText = graph.NewText(msg, draw.Fonts[draw.Face_list], color.White)
		fpsText.SetPosPivot(graph.ScrP(0.1, 0.1), v2.ZV)
	default:
	}
	if fpsText != nil {
		fpsText.Draw(window)
	}

	t := time.Now()
	dt = t.Sub(last).Seconds()
	last = t
	return nil
}

func main() {
	log.Start(time.Duration(DEFVAL.LogTimeoutMs)*time.Millisecond,
		time.Duration(DEFVAL.LogRetryMinMs)*time.Millisecond,
		time.Duration(DEFVAL.LogRetryMaxMs)*time.Millisecond,
		DEFVAL.LogIP, DEFVAL.LogHostName)

	if DEFVAL.DoProf {
		commons.StartProfile(roleName)
		defer commons.StopProfile(roleName)
	}

	if DEFVAL.FullScreen {
		ebiten.SetFullscreen(true)
		WinW, WinH = ebiten.MonitorSize()
		if DEFVAL.HalfResolution {
			WinW, WinH = WinW/2, WinH/2
		}
	} else {
		WinW = DEFVAL.WinW
		WinH = DEFVAL.WinH
	}

	graph.SetScreenSize(WinW, WinH)
	draw.LowQualityCosmoPoint(DEFVAL.LowQ)

	initClient()
	input.LoadConf("input_" + roleName + ".json")
	Data = commons.NewData()

	draw.InitFonts()
	draw.InitTexAtlas()
	commons.InitSignatureAtlas()

	createScenes()

	Client.Start()

	ebiten.SetRunnableInBackground(true)

	last = time.Now()

	showFps = time.Tick(time.Second)

	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "NAVIGATOR"); err != nil {
		log.Log(log.LVL_FATAL, err)
	}
}
