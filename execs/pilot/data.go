package main

import (
	. "github.com/Shnifer/magellan/commons"
	"sync"
)

type pilotData struct {
	mu sync.RWMutex

	//state data
	state  State
	bsp    BSP
	galaxy Galaxy

	//common data
	ship ShipPos
}

var Data pilotData

func (pd pilotData) getStateData(data []byte) chan struct{} {
	done := make(chan struct{})

	go func() {
		//anyway done, even with error
		defer close(done)

		//get state data
		md, err := MapData{}.Decode(data)
		if err != nil {
			panic("Weird state data:")
		}

		pd.mu.Lock()
		if bspDat, ok := md[PART_BSP]; ok {
			pd.bsp = BSP{}.Decode([]byte(bspDat))
		} else {
			pd.bsp = BSP{}
		}

		if galaxyDat, ok := md[PART_Galaxy]; ok {
			pd.galaxy = Galaxy{}.Decode([]byte(galaxyDat))
		} else {
			pd.galaxy = Galaxy{}
		}
		pd.mu.Unlock()

		initSceneState()
	}()

	return done
}

func (pd pilotData) commonSend() []byte {
	pd.mu.RLock()
	defer pd.mu.RUnlock()

	dat := pd.ship.Encode()
	md := make(MapData)
	md[PART_ShipPos] = dat
	res, err := md.Encode()
	if err != nil {
		panic("CommonSend " + err.Error())
	}
	return []byte(res)
}

func (pd pilotData) commonRecv(buf []byte) {
	pd.mu.Lock()
	defer pd.mu.Unlock()

	md, err := MapData{}.Decode(buf)
	if err != nil {
		panic("pilotData.commonRecv Can't decode mapData " + err.Error())
	}

	if part, ok := md[PART_ShipPos]; ok {
		pd.ship = ShipPos{}.Decode([]byte(part))
	}
}
