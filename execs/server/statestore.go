package main

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/peterbourgon/diskv"
	"github.com/pkg/errors"
	"strconv"
)

type RestoreRec struct {
	State      State
	CommonData CommonData
}

func (rh *roomHolder) saveRestorePoint(restore *diskv.Diskv) {
	var rec RestoreRec
	rh.stateMu.RLock()
	rh.commonMu.RLock()
	defer rh.stateMu.RUnlock()
	defer rh.commonMu.RUnlock()

	state := rh.curState
	if state.ShipID == "" || state.GalaxyID == "" {
		return
	}
	rec.State = state
	rec.CommonData = rh.commonData.Copy()
	rec.CommonData.ServerData.MsgID = 0
	rec.CommonData.ServerData.OtherShips = nil
	go saveRec(restore, rec)
}

func saveRec(restore *diskv.Diskv, rec RestoreRec) {
	ship := rec.State.ShipID
	i, ok := 1, false
	for !ok {
		ch := restore.KeysPrefix(ship+" - "+strconv.Itoa(i)+" - ", nil)
		_, exist := <-ch
		if exist {
			i++
		} else {
			ok = true
		}
	}
	key := ship + " - " + strconv.Itoa(i) + " - " + rec.State.GalaxyID
	restore.Write(key, rec.Encode())
}

//todo: any interface to run
func (rs *roomServer) loadRestorePoint(roomName string, ship string, n int) error {
	cancel := make(chan struct{})
	defer close(cancel)

	ch := rs.restore.KeysPrefix(ship+" - "+strconv.Itoa(n)+" - ", cancel)
	key, ok := <-ch
	if !ok {
		return errors.New("no such file")
	}

	dat, err := rs.restore.Read(key)
	if err != nil {
		return err
	}
	rec, err := RestoreRec{}.Decode(dat)
	if err != nil {
		return err
	}

	holder := rs.getHolder(roomName)
	state := rec.State
	if oldSub := holder.getSubscribe(); oldSub != nil {
		rs.storage.Unsubscribe(oldSub)
	}

	stateData, subscribe := rs.loadStateData(state)
	holder.rdyStateData(state, stateData, subscribe, false)
	holder.setCommon(rec.CommonData)
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
