package flow

import (
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/v2"
)

type updDrawPointer interface {
	update(dt float64)
	drawPoint(p Point) *graph.DrawQueue
}

type Point struct {
	lifeTime float64
	maxTime  float64
	pos      V2
	updDraw  updDrawPointer
	attr     map[string]float64
}

func (p Point) Req() *graph.DrawQueue {
	return p.updDraw.drawPoint(p)
}

type AttrF = func(p Point) float64

func NewAttrFs() map[string]AttrF {
	return make(map[string]AttrF)
}

type VelocityF func(pos V2) V2
type SpawnPosF func() (pos V2)

type Params struct {
	SpawnPeriod float64

	SpawnPos       SpawnPosF
	SpawnLife      func() float64
	SpawnUpdDrawer func() updDrawPointer

	VelocityF VelocityF
	AttrFs    map[string]AttrF
}

type Flow struct {
	params Params
	points []Point
	spawnT float64

	isActiveSpawn bool
	isEmpty       bool
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
		fp.AttrFs = make(map[string]func(p Point) float64)
	}

	return &Flow{
		params:        fp,
		points:        []Point{},
		isActiveSpawn: true,
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

	if f.isActiveSpawn {
		f.isEmpty = false
		//spawn new
		f.spawnT += dt
		for f.spawnT >= f.params.SpawnPeriod {
			f.spawnT -= f.params.SpawnPeriod
			f.newPoint()
		}
	} else if l == 0 {
		f.isEmpty = true
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
	p := Point{
		maxTime: f.params.SpawnLife(),
		pos:     f.params.SpawnPos(),
		updDraw: f.params.SpawnUpdDrawer(),
		attr:    make(map[string]float64),
	}
	f.points = append(f.points, p)
}

func (f *Flow) SetActive(activeSpawn bool) {
	f.isActiveSpawn = activeSpawn
}

func (f *Flow) IsEmpty() bool {
	return f.isEmpty
}
