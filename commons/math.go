package commons

import (
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
