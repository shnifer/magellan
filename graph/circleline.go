package graph

import (
	"github.com/Shnifer/magellan/v2"
	"image/color"
)

type CircleLineOpts struct {
	Params CamParams
	Clr    color.Color
	Layer  int
	PCount int
}

type CircleLine struct {
	opts   CircleLineOpts
	center v2.V2
	radius float64
	points []v2.V2
}

var clPoints map[int][]v2.V2

func init() {
	clPoints = make(map[int][]v2.V2)
}
func getPoints(n int) []v2.V2 {
	res, ok := clPoints[n]
	if ok {
		return res
	}
	res = make([]v2.V2, n)
	for i := 0; i < n; i++ {
		res[i] = v2.InDir(float64(i) * 360 / float64(n))
	}
	clPoints[n] = res
	return res
}

func NewCircleLine(center v2.V2, radius float64, opts CircleLineOpts) *CircleLine {
	res := &CircleLine{
		opts:   opts,
		center: center,
		radius: radius,
		points: getPoints(opts.PCount),
	}
	return res
}

func (cl *CircleLine) SetColor(clr color.Color) {
	cl.opts.Clr = clr
}
func (cl *CircleLine) SetPos(pos v2.V2) {
	cl.center = pos
}
func (cl *CircleLine) SetRadius(r float64) {
	cl.radius = r
}

func (cl *CircleLine) p(base v2.V2, r float64, i int) v2.V2 {
	return base.AddMul(cl.points[i], r)
}

func (cl *CircleLine) Req(Q *DrawQueue) {
	cam := cl.opts.Params.Cam
	base := cl.center
	if cam != nil {
		base = cam.Apply(base)
	}
	r := cl.radius
	if !cl.opts.Params.DenyScale && cam != nil {
		r *= cam.Scale * GS()
	}
	if r<1{
		return
	}
	from := cl.p(base, r, 0)
	var to v2.V2
	for i := 0; i < cl.opts.PCount; i++ {
		to = cl.p(base, r, (i+1)%cl.opts.PCount)
		LineScr(Q, from, to, cl.opts.Clr, cl.opts.Layer)
		from = to
	}
}
