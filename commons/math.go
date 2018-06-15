package commons

import (
	"github.com/Shnifer/magellan/v2"
	"math/rand"
)

func Clamp(x, min, max float64) float64 {
	switch {
	case x > max:
		return max
	case x < min:
		return min
	default:
		return x
	}
}

//add delta to X, but X can't be more than 1
func Add1(x *float64, delta float64) {
	*x += delta
	if *x > 1 {
		*x = 1
	}
	return
}

var gravityConst float64
var warpGravityConst float64

func SetGravityConsts(G, W float64) {
	gravityConst = G
	warpGravityConst = W
}

//gravity accelleration (g) from planet with given mass at given range
func Gravity(mass, lenSqr, zDist float64) float64 {
	d2 := lenSqr + zDist*zDist

	if d2 == 0 {
		return 0
	}

	return gravityConst * mass / d2

	//d2 = d2 * d2
	//return gravityConst * mass * lenSqr / d2
}

func SumGravityAcc(pos v2.V2, galaxy *Galaxy) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		G = Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), G)
	}
	return sumF
}

func WarpGravity(mass, lenSqr, velSqr, zDist float64) float64 {

	d2 := lenSqr + zDist*zDist
	d2 = d2 * d2
	if d2 == 0 {
		return 0
	}

	return warpGravityConst * mass * lenSqr / d2 * (1 + velSqr)
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
