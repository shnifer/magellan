package main

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"strconv"
	"github.com/pkg/errors"
)

type RestoreRec struct {
	State      State
	CommonData CommonData
}

func (s *roomServer) saveRestorePoint(roomName string) {
	var rec RestoreRec
	s.stateMu.RLock()
	s.commonMu.RLock()
	defer s.stateMu.RUnlock()
	defer s.commonMu.RUnlock()

	state := s.curState[roomName]
	if state.ShipID=="" || state.GalaxyID=="" {
		return
	}
	rec.State = state
	rec.CommonData = s.commonData[roomName].Copy()
	rec.CommonData.ServerData.MsgID=0
	rec.CommonData.ServerData.OtherShips = nil

	go saveRec(s, rec)
}

func saveRec(s *roomServer, rec RestoreRec) {
	ship := rec.State.ShipID
	i, ok :=1, false
	for !ok{
		ch:=s.restore.KeysPrefix(ship+" - "+strconv.Itoa(i),nil)
		_,exist:=<-ch
		if exist{
			i++
		} else {
			ok = true
		}
	}
	key := ship+" - "+strconv.Itoa(i)+" - "+rec.State.GalaxyID
	s.restore.Write(key, rec.Encode())
}

//todo: any interface to run
func (s *roomServer) loadRestorePoint(roomName string, ship string, n int) error{
	cancel:=make(chan struct{})
	defer close(cancel)

	ch:=s.restore.KeysPrefix(ship+" - "+strconv.Itoa(n),cancel)
	key, ok:=<-ch
	if !ok {
		return errors.New("no such file")
	}

	dat,err:=s.restore.Read(key)
	if err!=nil{
		return err
	}
	rec, err:=RestoreRec{}.Decode(dat)
	if err!=nil{
		return err
	}

	s.stateMu.Lock()
	s.curState[roomName] = rec.State
	s.stateMu.Unlock()

	stateData, subscribe := s.loadStateData(rec.State)

	s.stateDataMu.Lock()
	s.stateData[roomName] = stateData
	s.stateDataMu.Unlock()

	s.subsMu.Lock()
	if s.subscribes[roomName] != nil {
		s.storage.Unsubscribe(s.subscribes[roomName])
	}
	s.subscribes[roomName] = subscribe
	s.subsMu.Unlock()

	//todo:check for command queues
	s.commonMu.Lock()
	s.commonData[roomName] = rec.CommonData
	s.commonMu.Unlock()

	return nil
}

func (r RestoreRec) Encode() []byte {
	buf, err := json.Marshal(r)
	if err != nil {
		Log(LVL_ERROR, "Can't marshal RestoreRec ", err)
		return nil
	}
	return buf
}

func (RestoreRec) Decode(buf []byte) (r RestoreRec, err error) {
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return RestoreRec{}, err
	}
	return r, nil
}
