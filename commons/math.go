package commons

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

var gravityConst float64

func SetGravityConst(G float64) {
	gravityConst = G
}

//
func Gravity(mass, lenSqr, zdist float64) float64 {
	d2 := lenSqr + zdist*zdist
	d2 = d2 * d2
	if d2 == 0 {
		return 0
	}
	return gravityConst * mass * lenSqr / d2
}
