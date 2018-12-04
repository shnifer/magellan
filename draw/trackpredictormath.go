package draw

import (
	. "github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/v2"
	"log"
	"runtime"
	"sync"
)

//global mutex for 1 worker for prediction calculation and cache
var recalcPointsMu sync.Mutex

func (tp *TrackPredictor) recalcPoints() {
	recalcPointsMu.Lock()
	defer recalcPointsMu.Unlock()

	tp.mu.Lock()
	accel := tp.accel
	ss := tp.sessionTime
	ship := tp.ship
	calcTime := tp.startCalcTime
	tp.mu.Unlock()

	count := int(tp.opts.TrackLen*tp.opts.NumInSec) + 1
	points := make([]v2.V2, count)
	dt := 1 / float64(tp.opts.NumInSec)

	points[0] = ship.Pos
	var grav v2.V2

	gravGalaxy := tp.gps.get(ss)

	for i := 1; i < count; i++ {
		ss += dt
		if i%tp.opts.GravEach == 0 {
			gravGalaxy = tp.gps.get(ss)
		}
		grav = gravGalaxy.sumGrav(ship.Pos)
		ship.Vel.DoAddMul(v2.Add(grav, accel), dt)
		ship.Pos.DoAddMul(ship.Vel, dt)
		points[i] = ship.Pos
		if i%10 == 0 {
			runtime.Gosched()
		}
	}

	tp.mu.Lock()
	tp.points = points
	tp.isRunning = false
	tp.calcTime = calcTime
	tp.mu.Unlock()
}

//duplicate of SumGravityAcc for gravGalaxyT case
func (gi gravImage) sumGrav(pos v2.V2) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range gi {
		v = obj.pos.Sub(pos)
		len2 = v.LenSqr()
		G = Gravity(obj.mass, len2, obj.gDepth)
		sumF.DoAddMul(v.Normed(), G)
	}
	return sumF
}
