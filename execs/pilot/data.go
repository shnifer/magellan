package main

import (
	. "github.com/Shnifer/magellan/commons"
)

var state State

var bsp BSP
var galaxy Galaxy

func getStateData(data []byte) {
	md, err:=MapData{}.Decode(data)
	if err!=nil{
		panic("Weird state data:")
	}

	if bspDat,ok:=md[PART_BSP];ok{
		bsp = BSP{}.Decode([]byte(bspDat))
	} else {
		bsp = BSP{}
	}

	if galaxyDat,ok:=md[PART_Galaxy];ok{
		galaxy = Galaxy{}.Decode([]byte(galaxyDat))
	} else {
		galaxy = Galaxy{}
	}
}

func stateChanged(wanted string) {
	state = State{}.Decode(wanted)
}