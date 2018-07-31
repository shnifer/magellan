package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"time"
)

var (
	WinW int
	WinH int
)

var last time.Time
var GalaxyName string
var CurGalaxy *commons.Galaxy
var Scene *scene
var dt float64

var fpsText *graph.Text
var showFps <-chan time.Time

var changeGalaxy chan string
var sessionTime float64

func mainLoop(window *ebiten.Image) error {
	sessionTime = time.Now().Sub(commons.StartDateTime).Seconds()

	updateNamesAndNotes()
	updateBuildings()

	select {
	case newGalaxy := <-changeGalaxy:
		changeState(newGalaxy)
	default:
	}

	Scene.update(dt)
	if !ebiten.IsRunningSlowly() {
		Scene.draw(window)
		if fpsText != nil {
			fpsText.Draw(window)
		}
	}

	select {
	case <-showFps:
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %.0f", fps)
		fpsText = graph.NewText(msg, draw.Fonts[draw.Face_list], color.White)
		fpsText.SetPosPivot(graph.ScrP(0.1, 0.1), v2.ZV)
	default:
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
		commons.StartProfile(roleName, DEFVAL.DebugPort)
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

	initFlightStorage()
	initNamesStorage()

	graph.SetScreenSize(WinW, WinH)
	draw.LowQualityCosmoPoint(DEFVAL.LowQ)
	ebiten.SetVsyncEnabled(DEFVAL.VSync)

	draw.InitFonts()
	draw.InitTexAtlas()
	commons.InitSignatureAtlas()

	ebiten.SetRunnableInBackground(true)

	changeGalaxy = make(chan string, 1)
	changeGalaxy <- commons.WARP_Galaxy_ID

	last = time.Now()
	showFps = time.Tick(time.Second)

	Scene = newScene()
	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "MISSION CONTROL CENTER"); err != nil {
		log.Log(log.LVL_FATAL, err)
	}
}

func changeState(newGalaxy string) {
	loadNewGalaxy(newGalaxy)
	GalaxyName = newGalaxy
	getNames(newGalaxy)
	getBuildings(newGalaxy)
	Scene.init()
}
