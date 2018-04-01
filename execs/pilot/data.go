package main

import (
	. "github.com/Shnifer/magellan/commons"
	"sync"
)

const GalaxyPath = "res/galaxy/"

type pilotData struct {
	mu sync.RWMutex

	//state data
	state State
	bsp BSP
	galaxy Galaxy

	//common data
	ship ShipPos
}

var Data pilotData

func (pd pilotData) getStateData(data []byte) chan struct{} {
	done:=make(chan struct{})

	go func(){
		//anyway done, even with error
		defer close(done)

		md, err:=MapData{}.Decode(data)
		if err!=nil{
			panic("Weird state data:")
		}

		pd.mu.Lock()
		defer pd.mu.Unlock()
		if bspDat,ok:=md[PART_BSP];ok{
			pd.bsp = BSP{}.Decode([]byte(bspDat))
		} else {
			pd.bsp = BSP{}
		}

		if galaxyDat,ok:=md[PART_Galaxy];ok{
			pd.galaxy = Galaxy{}.Decode([]byte(galaxyDat))
		} else {
			pd.galaxy = Galaxy{}
		}
	}()

	return done
}

func (pd pilotData) stateChanged(wanted string) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	pd.state = State{}.Decode(wanted)
}


func (pd pilotData) commonSend() []byte {
	pd.mu.RLock()
	defer pd.mu.RUnlock()

	dat := pd.ship.Encode()
	md := make(MapData)
	md[PART_ShipPos] = dat
	res, err:=md.Encode()
	if err!=nil{
		panic("CommonSend "+err.Error())
	}
	return []byte(res)
}


func (pd pilotData) commonRecv(buf []byte) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	md,err:=MapData{}.Decode(buf)
	if err!=nil{
		panic("pilotData.commonRecv Can't decode mapData "+err.Error())
	}

	if part,ok:=md[PART_ShipPos]; ok{
		pd.ship = ShipPos{}.Decode([]byte(part))
	}
}