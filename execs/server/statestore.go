package main

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	. "github.com/Shnifer/magellan/log"
	"github.com/peterbourgon/diskv"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"strings"
	"time"
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

func (rs *roomServer) getRestorePoint(ship string, restoreN int) (RestoreRec, error) {
	cancel := make(chan struct{})
	defer close(cancel)

	ch := rs.restore.KeysPrefix(ship+" - "+strconv.Itoa(restoreN)+" - ", cancel)

	key, ok := <-ch
	if !ok {
		return RestoreRec{}, errors.New("no such file")
	}

	dat, err := rs.restore.Read(key)
	if err != nil {
		return RestoreRec{}, err
	}
	rec, err := RestoreRec{}.Decode(dat)
	if err != nil {
		return RestoreRec{}, err
	}
	return rec, nil
}

func (rs *roomServer) loadRestorePoint(roomName string, ship string, restoreN int) error {
	rs.RLock()
	defer rs.RUnlock()

	rec, err := rs.getRestorePoint(ship, restoreN)
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

type restorePoint struct {
	restN int
	memo  string
}

func (rs *roomServer) getShipRestoreList(ship string) []restorePoint {
	cancel := make(chan struct{})
	defer close(cancel)

	ArestN := make([]int, 0)
	Memos := make(map[int]string)
	res := make([]restorePoint, 0)

	ch := rs.restore.KeysPrefix(ship+" - ", cancel)

	for key := range ch {
		restS := strings.TrimPrefix(key, ship+" - ")
		ind := strings.Index(restS, " - ")
		if ind < 0 {
			continue
		}
		restS = restS[:ind]
		restN, err := strconv.Atoi(restS)
		if err != nil {
			Log(LVL_ERROR, "Wrong restore point key: ", key)
			continue
		}
		dat, err := rs.restore.Read(key)
		if err != nil {
			Log(LVL_ERROR, "Can't read restore point key: ", key)
			continue
		}
		rec, err := RestoreRec{}.Decode(dat)
		if err != nil {
			Log(LVL_ERROR, "Can't decode restore point rec: ", dat)
			continue
		}
		st := rec.CommonData.PilotData.SessionTime
		t := StartDateTime.Add(time.Duration(st) * time.Second)
		Memos[restN] = rec.State.GalaxyID + " " + t.Format(time.ANSIC)
		ArestN = append(ArestN, restN)
	}
	sort.Sort(sort.IntSlice(ArestN))

	for _, n := range ArestN {
		res = append(res, restorePoint{
			restN: n,
			memo:  Memos[n],
		})
	}

	return res
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

func (rs *roomServer) saveRoom(roomName string) {
	rs.RLock()
	defer rs.RUnlock()

	holder, ok := rs.holders[roomName]
	if !ok {
		Log(LVL_ERROR, "saveRoom runned for not found holder room ", roomName)
		return
	}
	holder.saveRestorePoint(rs.restore)
}

//runned as goroutine from console request
func (rs *roomServer) DoLoadRestore(shipId string, restoreN int, roomName string) {

	rec, err := rs.getRestorePoint(shipId, restoreN)
	if err != nil {
		return
	}

	//prev save
	rs.saveRoom(roomName)
	rs.loadMu.Lock()
	rs.loadPlans[roomName] = loadPlan{
		timeout:  time.Now().Add(time.Duration(DEFVAL.RestoreTimeoutS) * time.Second),
		state:    rec.State,
		restoreN: restoreN,
		shipId:   shipId,
	}
	rs.loadMu.Unlock()

	server.KillRoom(roomName)
	go loadDaemon(rs, roomName, rec.State)
}

func loadDaemon(rs *roomServer, roomName string, state State) {
	started := time.Now()
	tick := time.Tick(time.Second / 4)
	for range tick {
		if time.Since(started) > time.Duration(DEFVAL.RestoreTimeoutS)*time.Second {
			return
		}

		if doLoadDaemonCheck(rs, roomName, state) {
			return
		}
	}
}

func doLoadDaemonCheck(rs *roomServer, roomName string, state State) bool {
	rs.RLock()
	defer rs.RUnlock()

	holder, ok := rs.holders[roomName]
	if !ok {
		return false
	}

	if holder.getState() != startState {
		return false
	}

	return server.SetNewState(roomName, state.Encode(), true)
}
