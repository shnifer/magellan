package graph

import (
	"github.com/hajimehoshi/ebiten"
	"sort"
)

const (
	//Static back, usually one instance
	Z_STAT_BACKGROUND = iota * 100
	//Dynamic back, i.e. coordinate axises
	Z_BACKGROUND
	//under game object, i.e. selection glow
	Z_UNDER_OBJECT
	//main objects
	Z_GAME_OBJECT
	//Above objects. Gizmos and labels
	Z_ABOVE_OBJECT
	//non cam parts of UI
	Z_HUD
	//static HUD, for nice edges
	Z_STAT_HUD
)

type DrawF = func(image *ebiten.Image)

var DrawFZero = func(image *ebiten.Image) {}

type DrawReq struct {
	drawReq DrawF
	layer   int
	group   string
}

func NewReq(layer int, group string, f DrawF) DrawReq {
	return DrawReq{
		drawReq: f,
		layer:   layer,
		group:   group,
	}
}

type drawer interface {
	DrawF() (f DrawF, group string)
}

type drawQueuer interface {
	Req() *DrawQueue
}

type reqs []DrawReq

func (r reqs) Len() int      { return len(r) }
func (r reqs) Swap(x, y int) { r[x], r[y] = r[y], r[x] }
func (r reqs) Less(x, y int) bool {
	if r[x].layer != r[y].layer {
		return r[x].layer < r[y].layer
	}
	return r[x].group < r[y].group
}

type DrawQueue struct {
	reqs reqs
}

func NewDrawQueue() *DrawQueue {
	res := make(reqs, 0)
	return &DrawQueue{
		reqs: res,
	}
}

//For simple objects(like Sprite amd other graph primitives) that do not know it's layer
func (dq *DrawQueue) Add(drawer drawer, layer int) {
	f, group := drawer.DrawF()
	dq.reqs = append(dq.reqs, NewReq(layer, group, f))
}

//For game objects that create a set of requests with layers and groups
func (dq *DrawQueue) Append(drawQueuer drawQueuer) {
	dq.reqs = append(dq.reqs, drawQueuer.Req().reqs...)
}

func (dq *DrawQueue) Run(dest *ebiten.Image) {
	sort.Stable(dq.reqs)
	for _, req := range dq.reqs {
		req.drawReq(dest)
	}
}
