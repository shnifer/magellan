package draw

import (
	. "github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/v2"
	"image/color"
	"math"
	"runtime"
	"sync"
	"time"
)

type WarpPredictorOpts struct {
	Cam    *graph.Camera
	Sprite *graph.Sprite
	Clr    color.Color
	Layer  int

	Galaxy *Galaxy

	//in S
	UpdT     float64
	NumInSec int
	TrackLen int

	DrawMaxP int

	PowN float64
}

type WarpPredictor struct {
	opts WarpPredictorOpts

	mu sync.Mutex

	//reset by update
	startCalcTime time.Time
	pos           v2.V2
	distortion    float64
	dir           float64

	image gravImage

	isRunning bool
	points    []v2.V2
	calcTime  time.Time

	lastT time.Time
}

func NewWarpPredictor(opts WarpPredictorOpts) *WarpPredictor {
	if opts.UpdT == 0 {
		opts.UpdT = 1
	}
	if opts.NumInSec == 0 {
		opts.NumInSec = 1
	}

	image := warpGravImage(opts.Galaxy)

	opts.Galaxy = nil

	return &WarpPredictor{
		opts:   opts,
		image:  image,
		points: make([]v2.V2, 0),
	}
}

func (wp *WarpPredictor) Req(Q *graph.DrawQueue) {
	// real time in s to redraw TrackPredictior
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if time.Since(wp.lastT).Seconds() > wp.opts.UpdT && !wp.isRunning {
		wp.lastT = time.Now()
		wp.isRunning = true
		go wp.recalcPoints()
	}
	wp.drawPoints(Q)
}

//run under mutex
func (wp *WarpPredictor) drawPoints(Q *graph.DrawQueue) {
	//in ms, must be a round part of minute
	const markEach = 1

	if wp.points == nil || len(wp.points) == 0 || (wp.calcTime == time.Time{}) {
		return
	}

	cutTime := -time.Since(wp.calcTime).Seconds()

	//ms within last minute 0 -- 59999
	t := wp.calcTime
	timeOffset := float64(t.Second()) + float64(t.Nanosecond())/1000000000
	for timeOffset >= markEach {
		timeOffset -= markEach
	}
	//in ms
	dt := 1 / float64(wp.opts.NumInSec)

	drawCount := len(wp.points) - 1
	if wp.opts.Cam != nil {
		for i, p := range wp.points {
			if !wp.opts.Cam.PointInSpace(p) {
				drawCount = i
				break
			}
		}
	}
	if drawCount == 0 {
		return
	}

	var drawEach = 1
	if wp.opts.DrawMaxP > 0 {
		drawEach := drawCount/wp.opts.DrawMaxP + 1
		if drawEach > 10 {
			drawEach = 10
		}
		dt *= float64(drawEach)
	}

	var prev v2.V2
	var p v2.V2
	for i := 0; i < drawCount/drawEach+1; i++ {
		if i*drawEach >= len(wp.points) {
			break
		}
		p = wp.points[i*drawEach]
		if i > 0 && cutTime > 0 {
			graph.Line(Q, wp.opts.Cam, prev, p, wp.opts.Clr, wp.opts.Layer)
			if timeOffset >= markEach {
				timeOffset -= markEach
				k := timeOffset / dt
				markP := p.Mul(1-k).AddMul(prev, k)
				wp.opts.Sprite.SetPos(markP)
				Q.Add(wp.opts.Sprite, wp.opts.Layer+1)
			}
		}
		prev = p
		timeOffset += dt
		cutTime += dt
	}
}

func (wp *WarpPredictor) SetPosDistDir(shipPos v2.V2, distortion, dir float64) {
	wp.mu.Lock()
	wp.pos = shipPos
	wp.distortion = distortion
	wp.dir = dir
	wp.startCalcTime = time.Now()
	wp.mu.Unlock()
}

func warpGravImage(galaxy *Galaxy) gravImage {
	res := make(gravImage, 0)
	var imP gravP
	for _, p := range galaxy.Ordered {
		if p.IsVirtual || p.Mass == 0 {
			continue
		}
		imP = gravP{
			mass:   p.Mass,
			gDepth: p.GDepth,
			pos:    p.Pos,
		}
		res = append(res, imP)
	}
	return res
}

func (wp *WarpPredictor) recalcPoints() {
	recalcPointsMu.Lock()
	defer recalcPointsMu.Unlock()

	wp.mu.Lock()
	distortion := wp.distortion
	pos := wp.pos
	dir := wp.dir
	calcTime := wp.startCalcTime
	wp.mu.Unlock()

	count := int(wp.opts.TrackLen*wp.opts.NumInSec) + 1
	points := make([]v2.V2, count)
	dt := 1 / float64(wp.opts.NumInSec)

	points[0] = pos
	//warp update COPYPASTE
	gravK := math.Pow(distortion, wp.opts.PowN)
	vel := VelDistWarpK * distortion

	V := v2.InDir(dir).Mul(vel)
	var grav v2.V2
	for i := 1; i < count; i++ {
		grav = wp.image.sumWarpGrav(pos).Mul(gravK)
		V.DoAddMul(grav, dt)
		V = V.Normed().Mul(vel)
		pos.DoAddMul(V, dt)

		points[i] = pos
		if i%100 == 0 {
			runtime.Gosched()
		}
	}

	wp.mu.Lock()
	wp.points = points
	wp.isRunning = false
	wp.calcTime = calcTime
	wp.mu.Unlock()
}

//duplicate of SumGravityAcc for gravGalaxyT case
func (gi gravImage) sumWarpGrav(pos v2.V2) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range gi {
		v = obj.pos.Sub(pos)
		len2 = v.LenSqr()
		G = WarpGravity(obj.mass, len2, obj.gDepth)
		sumF.DoAddMul(v.Normed(), G)
	}
	return sumF
}
