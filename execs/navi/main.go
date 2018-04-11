package main

import (
	"fmt"
	"github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/input"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

const resPath = "res/navi/"
const texPath = "res/textures/"

var (
	WinW int
	WinH int
)

var last time.Time
var LocalData commons.TData

func mainLoop(window *ebiten.Image) error {
	dt := time.Since(last).Seconds()
	last = time.Now()

	input.Update()

	loadLocalData()

	Scenes.UpdateAndDraw(dt, window, !ebiten.IsRunningSlowly())

	sendLocalData()

	if commons.LOG_LEVEL <= commons.LVL_ERROR {
		fps := ebiten.CurrentFPS()
		msg := fmt.Sprintf("FPS: %v\ndt = %.2f\n", fps, dt)
		if ebiten.IsRunningSlowly() {
			msg = msg + "is running SLOWLY!\n"
		}
		ebitenutil.DebugPrint(window, msg)
	}

	return nil
}

func main() {
	startProfile()

	WinW = DEFVAL.WinW
	WinH = DEFVAL.WinH

	graph.SetScreenSize(WinW, WinH)

	initClient()
	input.LoadConf(resPath)

	createScenes()

	Client.Start()
	ebiten.SetFullscreen(DEFVAL.FullScreen)
	ebiten.SetRunnableInBackground(true)
	last = time.Now()
	if err := ebiten.Run(mainLoop, WinW, WinH, 1, "NAVIGATOR"); err != nil {
		log.Fatal(err)
	}

	stopProfile()
}

func loadLocalData() {
	defer commons.LogFunc("loadLocalData")()

	NetData.Mu.RLock()
	if NetData.State != LocalData.State {
		LocalData.State = NetData.State
		LocalData.StateData = NetData.StateData.Copy()
	}
	LocalData.CommonData = NetData.CommonData.Copy()
	NetData.Mu.RUnlock()
}

func sendLocalData() {
	defer commons.LogFunc("sendLocalData")()

	NetData.Mu.Lock()
	LocalData.CommonData.Part(DEFVAL.Role).FillNotNil(&NetData.CommonData)
	NetData.Mu.Unlock()
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
