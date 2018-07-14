package main

import (
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/draw"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/log"
	"github.com/hajimehoshi/ebiten"
	"time"
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
var showFps <-chan time.Time

func mainLoop(window *ebiten.Image) error {
	go mainStart()

	input.Update()
	//Pilot data must not be overwriten by other clients each tick, cz of Ship.Pos.Extrapolate
	if Data.CommonData.PilotData != nil {
		Data.CommonData.PilotData.MsgID++
	}
	Data.Update(DEFVAL.Role)

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	t := time.Now()
	dt = t.Sub(last).Seconds()
	last = t

	go mainStop()
	return nil
}

func mainStart() {
	time.Sleep(time.Nanosecond)
}
func mainStop() {
	time.Sleep(time.Nanosecond)
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

	commons.SetGravityConsts(DEFVAL.GravityConst, DEFVAL.WarpGravityConst)
	commons.SetVelDistWarpK(DEFVAL.VelDistWarpK)
	commons.SetWarpGravThreshold(DEFVAL.WarpGravThreshold)

	createScenes()

	Client.Start()

	ebiten.SetRunnableInBackground(true)

	last = time.Now()

	showFps = time.Tick(time.Second)

	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "PILOT"); err != nil {
		log.Log(log.LVL_FATAL, err)
	}
}
