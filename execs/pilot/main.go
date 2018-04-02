package main

import (
	"fmt"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/scene"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const DEBUG = true
const resPath = "res/pilot/"
const fontPath = "res/fonts/"

var (
	WinW int
	WinH int
)

var last time.Time

func mainLoop(window *ebiten.Image) error {
	dt := time.Since(last).Seconds()
	last = time.Now()

	input.Update()

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

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
	startProfile()

	WinW = DEFVAL.WinW
	WinH = DEFVAL.WinH

	graph.SetScreenSize(WinW, WinH)

	initClient()
	input.LoadConf(resPath)

	Scenes = scene.NewManager()

	face, err := graph.GetFace(fontPath+"phantom.ttf", 20)
	if err != nil {
		panic(err)
	}
	pauseScene := scene.NewPauseScene(face, Client.PauseReason)
	loginScene := NewLoginScene(face)
	Scenes.Install(scene_main, pauseScene, false)
	Scenes.Install(scene_pause, pauseScene, true)
	Scenes.Install(scene_login, loginScene, false)
	Scenes.SetOnPauseScene(scene_pause)
	Scenes.Activate(scene_pause, false)

	Client.Start()
	ebiten.SetFullscreen(DEFVAL.FullScreen)
	ebiten.SetRunnableInBackground(true)
	last = time.Now()
	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "PILOT"); err != nil {
		log.Fatal(err)
	}

	stopProfile()
}

func startProfile() {
	cpufn := DEFVAL.CpuProfFileName
	if cpufn != "" {
		f, err := os.Create(cpufn)
		if err != nil {
			log.Panicln("can't create cpu profile", cpufn, err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Panicln("can't start CPU profile ", err)
		}
	}
}

func stopProfile() {
	if DEFVAL.CpuProfFileName != "" {
		pprof.StopCPUProfile()
	}

	memfn := DEFVAL.MemProfFileName
	if memfn != "" {
		f, err := os.Create(memfn)
		if err != nil {
			log.Panicln("can't create mem profile", memfn)
		}
		runtime.GC()
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			log.Panicln("can't start mem profile", err)
		}
		f.Close()
	}
}
