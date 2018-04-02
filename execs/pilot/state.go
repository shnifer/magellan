package main

import (
	. "github.com/Shnifer/magellan/commons"
	"log"
)

const (
	scene_main  = "main"
	scene_pause = "pause"
	scene_login = "login"
)

func stateChanged(wanted string) {
	log.Println("state changed : ",wanted)
	state := State{}.Decode(wanted)

	Data.mu.Lock()
	Data.state = state
	Data.mu.Unlock()

	switch state.Special {
	case STATE_login:
		Scenes.Activate(scene_login, true)
	case STATE_cosmo:
		scene:=newCosmoScene()
		Scenes.Install(scene_main, scene, false)
		Scenes.Activate(scene_main, true)
	case STATE_warp:
	}
}

//called within Data.Mutex
func initSceneState(){
	var sceneName string

	switch Data.state.Special {
	case STATE_login:
		sceneName = scene_login
	case STATE_cosmo:
		sceneName = scene_main
	case STATE_warp:
		sceneName = scene_main
	}
	if sceneName!="" {
		Scenes.Init(sceneName)
	} else {
		log.Println("unknown scene to init for state = ",Data.state.Special)
	}
}


func pause() {
	log.Println("pause")
	Scenes.SetPaused(true)
}

func unpause() {
	log.Println("upause")
	Scenes.SetPaused(false)
}
