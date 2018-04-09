package main

import (
	. "github.com/Shnifer/magellan/commons"
	"log"
	"github.com/Shnifer/magellan/scene"
)

const (
	scene_main  = "main"
	scene_pause = "pause"
	scene_login = "login"
)

var Scenes *scene.Manager

func createScenes(){
	Scenes = scene.NewManager()

	pauseScene := scene.NewPauseScene(Fonts[Face_cap], Client.PauseReason)
	loginScene := NewLoginScene()
	Scenes.Install(scene_main, pauseScene, false)
	Scenes.Install(scene_pause, pauseScene, true)
	Scenes.Install(scene_login, loginScene, false)
	Scenes.SetAsPauseScene(scene_pause)
	Scenes.Activate(scene_pause, false)
}

func stateChanged(wanted string) {
	defer LogFunc("state.stateChanged")()
	log.Println("state changed : ", wanted)
	state := State{}.Decode(wanted)

	Data.setState(state)

	switch state.Special {
	case STATE_login:
		Scenes.Activate(scene_login, true)
	case STATE_cosmo:
		scene := newCosmoScene()
		Scenes.Install(scene_main, scene, false)
		Scenes.Activate(scene_main, true)
	case STATE_warp:
	}
}

//called within Data.Mutex
func initSceneState() {
	defer LogFunc("state.initSceneState")()

	var sceneName string

	switch Data.state.Special {
	case STATE_login:
		sceneName = scene_login
	case STATE_cosmo:
		sceneName = scene_main
	case STATE_warp:
		sceneName = scene_main
	}
	if sceneName != "" {
		Scenes.Init(sceneName)
	} else {
		log.Println("unknown scene to init for state = ", Data.state.Special)
	}
}

func pause() {
	defer LogFunc("state.pause")()
	Log(LVL_WARNING, "pause")
	Scenes.SetPaused(true)
}

func unpause() {
	defer LogFunc("state.unpause")()
	Log(LVL_WARNING, "upause")
	Scenes.SetPaused(false)
}

func discon() {
	Log(LVL_WARNING, "lost connect")
}

func recon() {
	Log(LVL_WARNING, "recon!")
}
