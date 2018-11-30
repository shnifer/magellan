package graph

import (
	"crypto/md5"
	"github.com/shnifer/magellan/v2"
	"math"
	"math/rand"
)

type WavedCircleOpts struct {
	Sprite  *Sprite
	Params  CamParams
	Layer   int
	PCount  int
	RandGen string
}

type waveParams struct {
	k1, p1, k2, p2 float64
	n1, n2         int
}

type WavedCircle struct {
	opts WavedCircleOpts

	center v2.V2
	radMin float64
	radMax float64
	points []v2.V2

	waveParams

	t float64
}

func defWaveParams(randgen string) waveParams {
	hash := md5.Sum([]byte(randgen))
	var seed int64
	seed += int64(hash[0])
	seed = seed << 8
	seed += int64(hash[1])
	seed = seed << 8
	seed += int64(hash[2])
	seed = seed << 8
	seed += int64(hash[3])
	rand.Seed(seed)
	rSign1 := float64(rand.Intn(2)*2 - 1)
	rSign2 := float64(rand.Intn(2)*2 - 1)
	return waveParams{
		k1: 1,
		k2: 0.3 + rand.Float64()*0.7,
		p1: rSign1 * (1 + rand.Float64()*3),
		p2: rSign2 * (1 + rand.Float64()*3),
		n1: rand.Intn(4) + 1,
		n2: rand.Intn(4) + 1,
	}
}

func NewWavedCircle(center v2.V2, radMin, radMax float64, opts WavedCircleOpts) *WavedCircle {
	res := &WavedCircle{
		opts:       opts,
		center:     center,
		radMin:     radMin,
		radMax:     radMax,
		points:     getPoints(opts.PCount),
		waveParams: defWaveParams(opts.RandGen),
	}
	return res
}

func (wc *WavedCircle) Update(dt float64) {
	wc.t += dt
}

func (wc *WavedCircle) SetPos(pos v2.V2) {
	wc.center = pos
}
func (wc *WavedCircle) SetRadius(min, max float64) {
	wc.radMin = min
	wc.radMax = max
}

func (wc *WavedCircle) p(base v2.V2, rMin, rMax float64, i int) v2.V2 {
	r := (rMin + rMax) / 2
	amp := (rMax - rMin) / 2
	prm := wc.waveParams
	pt := float64(i) / float64(wc.opts.PCount)

	wave1 := prm.k1 * math.Sin(wc.t*prm.p1+pt*2*math.Pi*float64(prm.n1))
	wave2 := prm.k2 * math.Sin(wc.t*prm.p2+pt*2*math.Pi*float64(prm.n2))

	r += amp * (wave1 + wave2) / (prm.k1 + prm.k2)

	return base.AddMul(wc.points[i], r)
}

func (wc *WavedCircle) Req(Q *DrawQueue) {
	base := wc.center
	rMin := wc.radMin
	rMax := wc.radMax
	spriteSize := rMin * 6 / float64(wc.opts.PCount)
	wc.opts.Sprite.SetSizeProportion(spriteSize)

	var pos v2.V2
	for i := 0; i < wc.opts.PCount; i++ {
		pos = wc.p(base, rMin, rMax, i)
		wc.opts.Sprite.SetPosAng(pos, float64(i)/float64(wc.opts.PCount)*360)
		Q.Add(wc.opts.Sprite, wc.opts.Layer)
	}
}
