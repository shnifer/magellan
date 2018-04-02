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

	switch Data.state.Special {
	case STATE_login:
		Scenes.Activate(scene_login, true)
	}
}


func initSceneState(){
	var sceneName string

	switch Data.state.Special {
	case STATE_login:
		sceneName = scene_login
	default:
		sceneName = scene_main
	}

	Scenes.Init(sceneName)
}


func pause() {
	log.Println("pause")
	Scenes.SetPaused(true)
}

func unpause() {
	log.Println("upause")
	Scenes.SetPaused(false)
}
