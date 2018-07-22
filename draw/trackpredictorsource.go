package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"math"
)

type gravGalaxyP struct {
	id        string
	parentInd int
	hasChild  bool
	pos       v2.V2
	orbit     float64
	period    float64
	angPhase float64
	mass      float64
	gDepth    float64
}
type gravGalaxyT []gravGalaxyP

type GravityPredictorSource struct {
	dt      float64
	sample  gravGalaxyT
	garantL int

	startT float64
	cache  []gravImage
}

type gravP struct {
	pos    v2.V2
	mass   float64
	gDepth float64
}

type gravImage []gravP

func NewGravityPredictorSource(galaxy *Galaxy, gravDt float64, garantLen int) *GravityPredictorSource {
	res := &GravityPredictorSource{
		sample:  newGravGalaxy(galaxy),
		dt:      gravDt,
		garantL: garantLen,
		cache:   make([]gravImage, 0, garantLen),
	}
	return res
}

func (gps *GravityPredictorSource) get(ss float64) gravImage {
	//empty or start
	if len(gps.cache) == 0 {
		gps.startT = ss
		res := gps.calc(ss)
		gps.cache = append(gps.cache, res)
		return res
	}

	//if we need earlier data -- just calc and return, do not want to cache
	if ss < gps.startT {
		return gps.calc(ss)
	}

	//we have this time cached
	reqInd := int(math.Round((ss - gps.startT) / gps.dt))
	if reqInd < len(gps.cache) {
		return gps.cache[reqInd]
	}

	//we have something far forward, so drop chache and start a new
	if reqInd > len(gps.cache)+gps.garantL {
		gps.cache = gps.cache[:0]
		return gps.get(ss)
	}

	//well we don't have requested, but we can hold in in current cache
	//let's calc all data up to
	var t float64
	for i := len(gps.cache); i <= reqInd; i++ {
		t = gps.startT + float64(i)*gps.dt
		gps.cache = append(gps.cache, gps.calc(t))
	}
	res := gps.cache[reqInd]

	//cut if needed
	if len(gps.cache) > 2*gps.garantL {
		startInd := len(gps.cache) - gps.garantL
		gps.startT = gps.startT + float64(startInd)*gps.dt
		a := gps.cache[startInd:]
		copy(gps.cache, a)
		gps.cache = gps.cache[:len(a)]
	}

	return res
}

//allocate copy data from sample, adjust time -- return
func (gps *GravityPredictorSource) calc(ss float64) gravImage {
	return gps.sample.calc(ss)
}

func newGravGalaxy(galaxy *Galaxy) gravGalaxyT {
	res := make(gravGalaxyT, 0, len(galaxy.Ordered))

	//map[id]index
	ord := make(map[string]int)
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 || !obj.IsVirtual {
			continue
		}
		p := gravGalaxyP{
			id:     obj.ID,
			pos:    obj.Pos,
			orbit:  obj.Orbit,
			period: obj.Period,
			angPhase: obj.AngPhase,
			mass:   obj.Mass,
			gDepth: obj.GDepth,
		}
		if obj.ParentID == "" {
			p.parentInd = -1
		} else {
			p.parentInd = ord[obj.ParentID]
			res[p.parentInd].hasChild = true
		}
		res = append(res, p)
		ord[obj.ID] = len(res) - 1
	}

	return res
}

//func (gg gravGalaxyT) loadPos(galaxy *Galaxy) {
//	for i, p := range gg {
//		gg[i].pos = galaxy.Points[p.id].Pos
//	}
//}

func (gg gravGalaxyT) calc(sessionTime float64) gravImage {
	res := make(gravImage, len(gg))

	var parent v2.V2
	var pos v2.V2
	//skip lvl 0 objects, they do not move
	for i, p := range gg {
		if p.parentInd == -1 {
			continue
		}
		parent = gg[p.parentInd].pos
		angle := (360 / p.period) * sessionTime + p.angPhase
		pos = parent.AddMul(v2.InDir(angle), p.orbit)
		gg[i].pos = pos
		res[i] = gravP{
			pos:    pos,
			mass:   p.mass,
			gDepth: p.gDepth,
		}
	}

	return res
}
