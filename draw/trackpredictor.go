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

	GPS *GravityPredictorSource

	//in S
	UpdT     float64
	NumInSec int
	GravEach int
	TrackLen int
	DrawMaxP int
}

type TrackPredictor struct {
	opts TrackPredictorOpts

	mu sync.Mutex

	//reset by update
	startCalcTime time.Time
	sessionTime   float64
	accel         v2.V2
	ship          RBData

	//created once, recalced pos before goroutine run
	gps *GravityPredictorSource

	isRunning bool
	points    []v2.V2
	calcTime  time.Time

	lastT time.Time
}

func NewTrackPredictor(opts TrackPredictorOpts) *TrackPredictor {
	if opts.GravEach == 0 {
		opts.GravEach = 1
	}
	if opts.UpdT == 0 {
		opts.UpdT = 1
	}
	if opts.NumInSec == 0 {
		opts.NumInSec = 1
	}
	if opts.DrawMaxP == 0 {
		opts.DrawMaxP = 1
	}
	return &TrackPredictor{
		opts:   opts,
		gps:    opts.GPS,
		points: make([]v2.V2, 0),
	}
}

func (tp *TrackPredictor) Req(Q *graph.DrawQueue) {
	// real time in s to redraw TrackPredictior
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if time.Since(tp.lastT).Seconds() > tp.opts.UpdT && !tp.isRunning {
		tp.lastT = time.Now()
		tp.isRunning = true
		//tp.gravGalaxy.loadPos(tp.opts.Galaxy)
		go tp.recalcPoints()
	}
	tp.drawPoints(Q)
}

//run under mutex
func (tp *TrackPredictor) drawPoints(Q *graph.DrawQueue) {
	//in s, must be a round part of minute
	const markEach = 1

	if tp.points == nil || len(tp.points) == 0 || (tp.calcTime == time.Time{}) {
		return
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

	drawCount := len(tp.points) - 1
	if tp.opts.Cam != nil {
		for i, p := range tp.points {
			if !tp.opts.Cam.PointInSpace(p) {
				drawCount = i
				break
			}
		}
	}
	if drawCount == 0 {
		return
	}
	var drawEach = 1
	if tp.opts.DrawMaxP > 0 {
		drawEach := drawCount/tp.opts.DrawMaxP + 1
		if drawEach > 10 {
			drawEach = 10
		}
		dt *= float64(drawEach)
	}

	var prev v2.V2
	var p v2.V2
	for i := 0; i < drawCount/drawEach+1; i++ {
		if i*drawEach >= len(tp.points) {
			break
		}
		p = tp.points[i*drawEach]
		if i > 0 && cutTime > 0 {
			graph.Line(Q, tp.opts.Cam, prev, p, tp.opts.Clr, tp.opts.Layer)
			for timeOffset >= markEach {
				timeOffset -= markEach
				if timeOffset < markEach {
					k := timeOffset / dt
					markP := p.Mul(1-k).AddMul(prev, k)
					tp.opts.Sprite.SetPos(markP)
					Q.Add(tp.opts.Sprite, tp.opts.Layer+1)
				}
			}
		}
		prev = p
		timeOffset += dt
		cutTime += dt
	}
}

func (tp *TrackPredictor) SetAccelSessionTimeShipPos(accel v2.V2, sessionTime float64, ship RBData) {
	tp.mu.Lock()
	tp.accel = accel
	tp.sessionTime = sessionTime
	tp.startCalcTime = time.Now()
	tp.ship = ship
	tp.mu.Unlock()
}
