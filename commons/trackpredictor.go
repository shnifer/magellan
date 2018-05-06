package commons

import (
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
	const dt = 1.0 / 10
	const markEach = 1 / dt
	const trackLen = 10

	const updT = 0.1

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
	for n := 1; n <= trackLen/dt; n++ {
		grav := SumGravity(ship.Pos, tp.data.Galaxy.Points)
		ship.Vel.DoAddMul(v2.Add(grav, accel), dt)
		prevPos := ship.Pos
		ship = ship.Extrapolate(dt)
		tp.q.Add(graph.Line(tp.cam, prevPos, ship.Pos, tp.clr), tp.layer)
		if (n % markEach) == 0 {
			tp.sprite.SetPos(ship.Pos)
			tp.q.Add(tp.sprite, tp.layer+1)
		}
	}

	return tp.q
}
