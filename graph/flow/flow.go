package flow

import (
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/v2"
)

type updDrawPointer interface {
	update(dt float64)
	drawPoint(p point) *graph.DrawQueue
}

type point struct {
	lifeTime float64
	maxTime  float64
	pos      V2
	updDraw  updDrawPointer
	attr     map[string]float64
}

func (p point) Req() *graph.DrawQueue {
	return p.updDraw.drawPoint(p)
}

type attrF = func(p point) float64

func NewAttrFs() map[string]attrF {
	return make(map[string]attrF)
}

type Params struct {
	SpawnPeriod float64

	SpawnPos       func() (pos V2)
	SpawnLife      func() float64
	SpawnUpdDrawer func() updDrawPointer

	VelocityF func(pos V2) V2
	AttrFs    map[string]attrF
}

type Flow struct {
	params Params
	points []point
	spawnT float64
}

func (fp Params) New() *Flow {
	if fp.SpawnPeriod == 0 {
		panic("NewFlow: zero SpawnPeriod")
	}
	if fp.VelocityF == nil {
		fp.VelocityF = func(V2) V2 { return ZV }
	}
	if fp.SpawnPos == nil {
		fp.SpawnPos = func() V2 { return ZV }
	}
	if fp.SpawnLife == nil {
		fp.SpawnLife = func() float64 { return 1 }
	}
	if fp.AttrFs == nil {
		fp.AttrFs = make(map[string]func(p point) float64)
	}

	return &Flow{
		params: fp,
		points: []point{},
	}
}

func (f *Flow) Update(dt float64) {
	//check for life time
	l := len(f.points)
	for i := 0; i < l; i++ {
		f.points[i].lifeTime += dt
		if f.points[i].lifeTime > f.points[i].maxTime {
			f.points[i] = f.points[l-1]
			f.points = f.points[:l-1]
			l--
		}
	}

	//spawn new
	f.spawnT += dt
	for f.spawnT >= f.params.SpawnPeriod {
		f.spawnT -= f.params.SpawnPeriod
		f.newPoint()
	}

	//move
	for i, p := range f.points {
		vel := f.params.VelocityF(p.pos)
		p.pos.DoAddMul(vel, dt)
		f.points[i] = p
	}

	//attr update
	for i, p := range f.points {
		for name, F := range f.params.AttrFs {
			f.points[i].attr[name] = F(p)
		}
	}

	//draw update
	for _, p := range f.points {
		p.updDraw.update(dt)
	}
}

func (f *Flow) Req() *graph.DrawQueue {
	res := graph.NewDrawQueue()
	for _, p := range f.points {
		res.Append(p)
	}
	return res
}

func (f *Flow) newPoint() {
	p := point{
		maxTime: f.params.SpawnLife(),
		pos:     f.params.SpawnPos(),
		updDraw: f.params.SpawnUpdDrawer(),
		attr:    make(map[string]float64),
	}
	f.points = append(f.points, p)
}