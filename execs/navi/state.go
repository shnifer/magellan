package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/scene"
	"golang.org/x/image/colornames"
	"log"
)

const (
	scene_main  = "main"
	scene_pause = "pause"
	scene_login = "login"
)

var Scenes *scene.Manager

func createScenes() {
	defer LogFunc("createScenes")()

	Scenes = scene.NewManager()

	pauseScene := scene.NewPauseScene(Fonts[Face_cap], Client.PauseReason)
	loginScene := scene.NewCaptionSceneString(Fonts[Face_cap], colornames.Goldenrod,
		"waiting for login on other terminal")
	Scenes.Install(scene_main, pauseScene, false)
	Scenes.Install(scene_pause, pauseScene, true)
	Scenes.Install(scene_login, loginScene, true)
	Scenes.SetAsPauseScene(scene_pause)
	Scenes.Activate(scene_pause, false)
	Scenes.WaitDone()
}

//Network cycle - direct handler
func stateChanged(wanted string) {
	defer LogFunc("state.stateChanged " + wanted)()

	state := State{}.Decode(wanted)

	Data.SetState(state)

	switch state.StateID {
	case STATE_login:
		Scenes.Activate(scene_login, true)
	case STATE_cosmo:
		newScene := newCosmoScene()
		Scenes.Install(scene_main, newScene, false)
		Scenes.Activate(scene_main, true)
	case STATE_warp:
	}
}

//Network cycle - handler in goroutine
func initSceneState() {
	defer LogFunc("state.initSceneState")()

	stateID := Data.GetState().StateID

	var sceneName string

	switch stateID {
	case STATE_login:
		sceneName = scene_login
	case STATE_cosmo:
		sceneName = scene_main
	case STATE_warp:
		sceneName = scene_main
	}
	if sceneName != "" {
		Scenes.Init(sceneName)
		Scenes.WaitDone()
	} else {
		log.Println("unknown scene to init for state = ", stateID)
	}
}

//Network cycle - direct handler
func onCommand(command string) {
	Scenes.OnCommand(command)
}

//Network cycle - direct handler
func pause() {
	defer LogFunc("state.pause")()
	Log(LVL_WARNING, "pause")
	Scenes.SetPaused(true)
}

//Network cycle - direct handler
func unpause() {
	defer LogFunc("state.unpause")()
	Log(LVL_WARNING, "upause")
	Scenes.SetPaused(false)
}

//Network cycle - direct handler
func discon() {
	Log(LVL_WARNING, "lost connect")
}

//Network cycle - direct handler
func recon() {
	Log(LVL_WARNING, "recon!")
}
