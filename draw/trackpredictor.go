package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"image/color"
	"sync"
	"time"
)

type TrackPredictorOpts struct {
	Cam    *graph.Camera
	Sprite *graph.Sprite
	Clr    color.Color
	Layer  int
	Galaxy *Galaxy

	//in S
	UpdT     float64
	NumInSec int
	TrackLen int
}

type TrackPredictor struct {
	opts TrackPredictorOpts

	mu sync.Mutex

	//reset by update
	sessionTime float64
	accel       v2.V2
	ship        RBData

	//created once, recalced pos before goroutine run
	gravGalaxy gravGalaxyT

	isRunning bool
	points    []v2.V2
	calcTime  time.Time

	lastT time.Time
}

func NewTrackPredictor(opts TrackPredictorOpts) *TrackPredictor {
	return &TrackPredictor{
		opts:       opts,
		gravGalaxy: newGravGalaxy(opts.Galaxy),
		points:     make([]v2.V2, 0),
	}
}

func (tp *TrackPredictor) Req() *graph.DrawQueue {
	// real time in s to redraw TrackPredictior
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if time.Since(tp.lastT).Seconds() > tp.opts.UpdT && !tp.isRunning {
		tp.lastT = time.Now()
		tp.isRunning = true
		tp.gravGalaxy.loadPos(tp.opts.Galaxy)
		go tp.recalcPoints()
	}
	return tp.drawPoints()
}

//run under mutex
func (tp *TrackPredictor) drawPoints() *graph.DrawQueue {
	//in ms, must be a round part of minute
	const markEach = 1

	Q := graph.NewDrawQueue()
	if tp.points == nil {
		return Q
	}

	cutTime := -time.Since(tp.calcTime).Seconds()

	//ms within last minute 0 -- 59999
	t := tp.calcTime
	timeOffset := float64(t.Second()) + float64(t.Nanosecond())/1000000000
	for timeOffset >= markEach {
		timeOffset -= markEach
	}
	//in ms
	dt := 1 / float64(tp.opts.NumInSec)

	var prev v2.V2
	for i, p := range tp.points {
		if i > 0 && cutTime > 0 {
			Q.Add(graph.Line(tp.opts.Cam, prev, p, tp.opts.Clr), tp.opts.Layer)
			if timeOffset >= markEach {
				timeOffset -= markEach
				k := timeOffset / dt
				markP := p.Mul(1-k).AddMul(prev, k)
				tp.opts.Sprite.SetPos(markP)
				Q.Add(tp.opts.Sprite, tp.opts.Layer+1)
			}
		}
		prev = p
		timeOffset += dt
		cutTime += dt
	}

	return Q
}

func (tp *TrackPredictor) SetAccelSessionTimeShipPos(accel v2.V2, sessionTime float64, ship RBData) {
	tp.mu.Lock()
	tp.accel = accel
	tp.sessionTime = sessionTime
	tp.ship = ship
	tp.mu.Unlock()
}
