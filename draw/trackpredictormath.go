package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"time"
)

type gravP struct {
	id        string
	parentInd int
	pos       v2.V2
	orbit     float64
	period    float64
	mass      float64
	gDepth    float64
}
type gravGalaxyT []gravP

func (tp *TrackPredictor) recalcPoints() {
	tp.mu.Lock()
	accel := tp.accel
	ss := tp.sessionTime
	ship := tp.ship
	tp.calcTime = time.Now()
	tp.mu.Unlock()

	count := int(tp.opts.TrackLen*tp.opts.NumInSec) + 1
	points := make([]v2.V2, count)
	dt := 1 / float64(tp.opts.NumInSec)

	points[0] = ship.Pos
	var grav v2.V2

	for i := 1; i < count; i++ {
		ss += dt
		tp.gravGalaxy.update(ss)
		grav = tp.gravGalaxy.sumGrav(ship.Pos)
		ship.Vel.DoAddMul(v2.Add(grav, accel), dt)
		ship.Pos.DoAddMul(ship.Vel, dt)
		points[i] = ship.Pos
	}

	tp.mu.Lock()
	tp.points = points
	tp.isRunning = false
	tp.mu.Unlock()
}

func newGravGalaxy(galaxy *Galaxy) gravGalaxyT {
	res := make(gravGalaxyT, 0, len(galaxy.Ordered))

	ord := make(map[string]int)
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		p := gravP{
			id:     obj.ID,
			pos:    obj.Pos,
			orbit:  obj.Orbit,
			period: obj.Period,
			mass:   obj.Mass,
			gDepth: obj.GDepth,
		}
		if obj.ParentID == "" {
			p.parentInd = -1
		} else {
			p.parentInd = ord[obj.ParentID]
		}
		res = append(res, p)
		ord[obj.ID] = len(res) - 1
	}

	return res
}

func (gg gravGalaxyT) loadPos(galaxy *Galaxy) {
	for i, p := range gg {
		gg[i].pos = galaxy.Points[p.id].Pos
	}
}

func (gg gravGalaxyT) update(sessionTime float64) {

	var parent v2.V2

	//skip lvl 0 objects, they do not move
	for i, p := range gg {
		if p.parentInd == -1 {
			continue
		}

		parent = gg[p.parentInd].pos

		angle := (360 / p.period) * sessionTime
		gg[i].pos = parent.AddMul(v2.InDir(angle), p.orbit)
	}
}

//duplicate of SumGravityAcc for gravGalaxyT case
func (gg gravGalaxyT) sumGrav(pos v2.V2) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range gg {
		v = obj.pos.Sub(pos)
		len2 = v.LenSqr()
		G = Gravity(obj.mass, len2, obj.gDepth)
		sumF.DoAddMul(v.Normed(), G)
	}
	return sumF
}
