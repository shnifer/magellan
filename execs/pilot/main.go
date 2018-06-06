package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/log"
	"github.com/hajimehoshi/ebiten"
	"time"
	"image/color"
	"github.com/Shnifer/magellan/v2"
)

const roleName = "pilot"

var (
	WinW int
	WinH int
)

var last time.Time
var Data commons.TData
var dt float64

var fpsText *graph.Text
var showFps  <-chan time.Time

func mainLoop(window *ebiten.Image) error {
	input.Update()

	//Pilot data must not be overwriten by other clients each tick, cz of Ship.Pos.Extrapolate
	if Data.CommonData.PilotData != nil {
		Data.CommonData.PilotData.MsgID++
	}
	Data.Update(DEFVAL.Role)

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	select {
	case <-showFps:
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %.0f\nALT-F4 to close\nWASD to control\nQ-E scale\nSPACE - stop\nENTER - reset position", fps)
		fpsText=graph.NewText(msg,draw.Fonts[draw.Face_list],color.White)
		fpsText.SetPosPivot(graph.ScrP(0.1,0.1),v2.ZV)
	default:
	}
	if fpsText!=nil {
		fpsText.Draw(window)
	}

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
		if DEFVAL.HalfResolution {
			WinW,WinH = WinW/2, WinH/2
		}
	} else {
		WinW = DEFVAL.WinW
		WinH = DEFVAL.WinH
	}

	graph.SetScreenSize(WinW, WinH)

	initClient()
	input.LoadConf("input_" + roleName + ".json")
	Data = commons.NewData()

	draw.InitFonts()
	draw.InitTexAtlas()
	commons.InitSignatureAtlas()

	commons.SetGravityConsts(DEFVAL.GravityConst, DEFVAL.WarpGravityConst)

	createScenes()

	Client.Start()

	ebiten.SetRunnableInBackground(true)

	last = time.Now()

	showFps=time.Tick(time.Second)

	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "PILOT"); err != nil {
		log.Log(log.LVL_FATAL, err)
	}
}
