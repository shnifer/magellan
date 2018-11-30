package main

import (
	. "github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/draw"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/scene"
)

const (
	scene_cosmo = "cosmo"
	scene_warp  = "warp"
	scene_pause = "pause"
	scene_login = "login"
)

var Scenes *scene.Manager

func createScenes() {
	defer LogFunc("createScenes")()

	Scenes = scene.NewManager()

	pauseScene := scene.NewPauseScene(Fonts[Face_cap], Client.PauseReason)
	loginScene := NewLoginScene()
	cosmoScene := newCosmoScene()
	warpScene := newWarpScene()
	Scenes.Install(scene_pause, pauseScene, true)
	Scenes.Install(scene_login, loginScene, false)
	Scenes.Install(scene_cosmo, cosmoScene, false)
	Scenes.Install(scene_warp, warpScene, false)
	Scenes.SetAsPauseScene(scene_pause)
	Scenes.Activate(scene_pause, false)
	Scenes.WaitDone()
}

func stateChanged(wanted string) {
	defer LogFunc("state.stateChanged " + wanted)()

	state := State{}.Decode(wanted)

	logKeys := make(map[string]string, 5)
	logKeys["Room"] = DEFVAL.Room
	logKeys["Role"] = DEFVAL.Role
	logKeys["Galaxy"] = state.GalaxyID
	logKeys["Ship"] = state.ShipID
	logKeys["State"] = state.StateID

	SetLogFields(logKeys)

	Data.SetState(state)

	switch state.StateID {
	case STATE_login:
		Scenes.Activate(scene_login, true)
	case STATE_cosmo:
		Scenes.Activate(scene_cosmo, true)
	case STATE_warp:
		Scenes.Activate(scene_warp, true)
	}
}

func initSceneState() {
	defer LogFunc("state.initSceneState")()

	stateID := Data.GetState().StateID

	var sceneName string

	switch stateID {
	case STATE_login:
		sceneName = scene_login
	case STATE_cosmo:
		sceneName = scene_cosmo
	case STATE_warp:
		sceneName = scene_warp
	}

	if sceneName != "" {
		Scenes.Init(sceneName)
		Scenes.WaitDone()
	} else {
		Log(LVL_ERROR, "unknown scene to init for state = ", stateID)
	}
}

func onCommand(command string) {
	Scenes.OnCommand(command)
}

func pause() {
	defer LogFunc("state.pause")()
	Log(LVL_WARN, "pause")
	Scenes.SetPaused(true)
}

func unpause() {
	defer LogFunc("state.unpause")()
	Log(LVL_WARN, "upause")
	Scenes.SetPaused(false)
}

func discon() {
	Log(LVL_WARN, "lost connect")
}

func recon() {
	Log(LVL_WARN, "recon!")
}
