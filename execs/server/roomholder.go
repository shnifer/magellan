package main

import (
	. "github.com/shnifer/magellan/commons"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/storage"
	"sync"
)

type roomHolder struct {
	roomName string

	stateMu  sync.RWMutex
	curState State

	stateDataMu sync.RWMutex
	stateData   StateData

	subsMu    sync.RWMutex
	subscribe chan storage.Event

	commonMu   sync.RWMutex
	commonData CommonData
}

func newRoomHolder(name string) *roomHolder {
	return &roomHolder{
		roomName:   name,
		commonData: CommonData{}.Empty(),
	}
}

func (rh *roomHolder) getState() State {
	rh.stateMu.RLock()
	defer rh.stateMu.RUnlock()
	return rh.curState
}

func (rh *roomHolder) getCommon() (commonData CommonData) {
	rh.commonMu.RLock()
	defer rh.commonMu.RUnlock()
	return rh.commonData.Copy()
}

func (rh *roomHolder) setCommon(cd CommonData) {
	rh.commonMu.Lock()
	defer rh.commonMu.Unlock()

	cd.FillNotNil(&rh.commonData)
}

func (rh *roomHolder) getCommonServerData() ServerData {
	rh.commonMu.RLock()
	defer rh.commonMu.RUnlock()
	return *rh.commonData.Copy().ServerData
}

func (rh *roomHolder) setCommonServerData(sd ServerData) {
	rh.commonMu.Lock()
	defer rh.commonMu.Unlock()

	rh.commonData.ServerData = &sd
}

func (rh *roomHolder) getSubscribe() chan storage.Event {
	rh.subsMu.RLock()
	defer rh.subsMu.RUnlock()

	return rh.subscribe
}

func (rh *roomHolder) rdyStateData(state State, stateData StateData, subscribe chan storage.Event, generateCommon bool) {
	rh.stateMu.Lock()
	defer rh.stateMu.Unlock()
	rh.stateDataMu.Lock()
	defer rh.stateDataMu.Unlock()
	rh.subsMu.Lock()
	defer rh.subsMu.Unlock()
	rh.commonMu.Lock()
	defer rh.commonMu.Unlock()

	prevState := rh.curState
	rh.curState = state
	rh.stateData = stateData
	rh.subscribe = subscribe

	if generateCommon {
		commonData := rh.commonData.Copy()
		genData := generateCommonData(commonData, stateData, state, prevState)
		rh.commonData = genData
	}
}

func (rh *roomHolder) getStateData() StateData {
	rh.stateDataMu.RLock()
	defer rh.stateDataMu.RUnlock()

	return rh.stateData.Copy()
}

func (rh *roomHolder) doUpdateSubscribes() {
	rh.subsMu.RLock()
	defer rh.subsMu.RUnlock()

	if rh.subscribe == nil {
		return
	}

	rh.stateDataMu.Lock()
	defer rh.stateDataMu.Unlock()

loop:
	for {
		select {
		case event, ok := <-rh.subscribe:
			if !ok {
				break loop
			}

			server.AddCommand(rh.roomName, EventToCommand(event))

			fk := event.Key.FullKey()
			switch event.Type {
			case storage.Add:
				b, err := Building{}.Decode([]byte(event.Data))
				if err != nil {
					Log(LVL_ERROR, "doUpdateSubscribes can't Decode building ", event.Data, " ", err)
					continue
				}
				b.FullKey = fk

				rh.stateData.Buildings[b.FullKey] = b
			case storage.Remove:
				delete(rh.stateData.Buildings, fk)
			default:
				Log(LVL_ERROR, "doUpdateSubscribes unknown event type ", event.Type)
			}
		default:
			break loop
		}
	}
}
