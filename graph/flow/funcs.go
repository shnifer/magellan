package flow

import (
	. "github.com/Shnifer/magellan/v2"
	"math"
	"math/rand"
)

func SinLifeTime(med, dev, period float64) func(p point) float64 {
	return func(p point) float64 {
		return med + math.Sin(p.lifeTime/period*2*math.Pi)*dev
	}
}
func SinMaxTime(med, dev, periods float64) func(p point) float64 {
	return func(p point) float64 {
		return med + math.Sin(p.lifeTime/p.maxTime*periods*2*math.Pi)*dev
	}
}

func ComposeRadial(tang, norm func(l, w float64) float64) func(V2) V2 {
	return func(pos V2) (vel V2) {
		l := pos.Len()
		w := pos.Dir()
		t := tang(l, w)
		n := norm(l, w)
		return Add(pos.Mul(n/l), pos.Rotate90().Mul(t/l))
	}
}

func ComposeDecart(vx, vy func(x, y float64) float64) func(V2) V2 {
	return func(pos V2) (vel V2) {
		return V2{X: vx(pos.X, pos.Y), Y: vy(pos.X, pos.Y)}
	}
}

func ConstC(val float64) func(float64, float64) float64 {
	return func(float64, float64) float64 {
		return val
	}
}

func LineRand(min, max float64) func() float64 {
	return func() float64 {
		return min + rand.Float64()*(max-min)
	}
}

func NormRand(center, devPercent float64) func() float64 {
	return func() float64 {
		return center * KDev(devPercent)
	}
}

func RandomInCirc(R float64) func() V2 {
	return func() V2 {
		return RandomInCircle(R)
	}
}

func RandomOnSide(sideOrt V2, wide float64) func() V2 {
	return func() V2 {
		l := 1 - rand.Float64()*wide
		wOrt := sideOrt.Rotate90()
		return sideOrt.Mul(l).AddMul(wOrt, rand.Float64()*2-1)
	}
}

//Возвращает коэффициент нормальной дистрибуций
//сигма в процентах devProcent
//68% попадут в (100-devProcent, 100+devProcent)
//95% попадут в (100-2*devProcent, 100+2*devProcent)
//Отклонения больше 3 сигма ограничиваются
func KDev(devProcent float64) float64 {
	r := rand.NormFloat64()
	if r > 3 {
		r = 3
	}
	if r < (-3) {
		r = -3
	}
	r = 1 + r*devProcent/100
	if r < 0 {
		r = 0.00001
	}
	return r
}
