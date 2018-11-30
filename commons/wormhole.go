package commons

import (
	"encoding/json"
	"github.com/pkg/errors"
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/static"
	"strconv"
	"time"
)

type WormHole struct {
	ID       int
	System   string
	TimePlan []int
}

const WormHoleYouDIE = "DIE!"

const WormHoleFN = "wormholes.json"

//in seconds
const WormHolePeriod = 60 * 15

var wormHoles map[int]*WormHole
var whBySystem map[string]*WormHole

type WormHoleDirSys struct {
	Src, Dest string
}

type WormHoleDirN struct {
	Src, Dest int
}

func InitWormHoles() {
	wormHoles = make(map[int]*WormHole)
	whBySystem = make(map[string]*WormHole)

	dat, err := static.Load("DB", WormHoleFN)
	if err != nil {
		Log(LVL_ERROR, err)
	}
	var whs []*WormHole
	err = json.Unmarshal(dat, &whs)
	if err != nil {
		Log(LVL_ERROR, err)
	}
	for _, v := range whs {
		wormHoles[v.ID] = v
		whBySystem[v.System] = v
	}
}

func GetWormHoleTarget(src string) (string, error) {
	wh, ok := whBySystem[src]
	if !ok {
		return "", errors.New("Not found source wormhole in system " + src)
	}
	targetId := wh.getTarget()
	if targetId <= 0 {
		return WormHoleYouDIE, nil
	}
	target, ok := wormHoles[targetId]
	if !ok {
		return "", errors.New("Not found target wormhole id " + strconv.Itoa(targetId))
	}
	return target.System, nil
}

func GetCurrentWormHoleDirectionSys() []WormHoleDirSys {
	res := make([]WormHoleDirSys, 0)
	for src, wh := range whBySystem {
		dest := wh.getTarget()
		if dest > 0 {
			res = append(res, WormHoleDirSys{Src: src, Dest: wormHoles[dest].System})
		}
	}
	return res
}

func GetCurrentWormHoleDirectionN() map[string]WormHoleDirN {
	res := make(map[string]WormHoleDirN)
	for src, wh := range wormHoles {
		dest := wh.getTarget()
		res[wh.System] = WormHoleDirN{Src: src, Dest: dest}
	}
	return res
}

func (wh WormHole) getTarget() int {
	l := len(wh.TimePlan)
	if l == 0 {
		return 0
	}
	n := int(time.Now().Unix()/WormHolePeriod) % l
	return wh.TimePlan[n]
}

func GetWormHolesNs() map[int]struct{} {
	res := make(map[int]struct{})
	for i := range wormHoles {
		res[i] = struct{}{}
	}
	return res
}
