package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"image/color"
	"time"
)

type TrackPredictor struct {
	cam    *graph.Camera
	sprite *graph.Sprite
	data   *TData
	mode   int
	clr    color.Color
	layer  int

	q *graph.DrawQueue

	lastT time.Time
}

const (
	Track_CurrentThrust int = iota
	Track_ZeroThrust
)

func NewTrackPredictor(cam *graph.Camera, sprite *graph.Sprite, data *TData, mode int, clr color.Color, layer int) *TrackPredictor {
	return &TrackPredictor{
		sprite: sprite,
		data:   data,
		cam:    cam,
		clr:    clr,
		layer:  layer,
		mode:   mode,
	}
}

func (tp *TrackPredictor) Req() *graph.DrawQueue {
	//flytime step
	const dt = 1.0 / 10
	const recalcGravEach = 3
	//graph mark eack Nth
	const markEach = 1 / dt
	//len of prediction in seconds
	const trackLen = 8

	// real time in s to redraw TrackPredictior
	const updT = 1.0 / 10

	if tp.q != nil && time.Since(tp.lastT).Seconds() < updT {
		return tp.q
	}
	tp.lastT = time.Now()

	tp.q = graph.NewDrawQueue()

	var accel v2.V2
	switch tp.mode {
	case Track_CurrentThrust:
		accel = tp.data.PilotData.ThrustVector
	case Track_ZeroThrust:
		accel = v2.ZV
	}

	ship := tp.data.PilotData.Ship
	var grav v2.V2
	for n := 0; n <= trackLen/dt; n++ {
		if (n%recalcGravEach) ==0 {
			grav = SumGravityF(ship.Pos, tp.data.Galaxy).Mul(1/tp.data.BSP.Mass)
		}
		ship.Vel.DoAddMul(v2.Add(grav, accel), dt)
		prevPos := ship.Pos
		ship = ship.Extrapolate(dt)
		tp.q.Add(graph.Line(tp.cam, prevPos, ship.Pos, tp.clr), tp.layer)
		if (n % markEach) == 0 && n!=0{
			tp.sprite.SetPos(ship.Pos)
			tp.q.Add(tp.sprite, tp.layer+1)
		}
	}

	return tp.q
}
