package commons

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/static"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type WormHole struct {
	ID       int
	System   string
	TimePlan []int
}

const WarmHoleYouDIE = "DIE!"

const WormHoleFN = "wormholes.json"

//in seconds
const WormHolePeriod = 60 * 15

var wormHoles map[int]*WormHole
var whBySystem map[string]*WormHole

func InitWormHoles() {
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
		return WarmHoleYouDIE, nil
	}
	target, ok := wormHoles[targetId]
	if !ok {
		return "", errors.New("Not found target wormhole id " + strconv.Itoa(targetId))
	}
	return target.System, nil
}

func (wh WormHole) getTarget() int {
	l := len(wh.TimePlan)
	if l == 0 {
		return 0
	}
	n := int(time.Now().Unix()/WormHolePeriod) % l
	return wh.TimePlan[n]
}
