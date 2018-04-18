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
var warpGravityConst float64

func SetGravityConsts(G, W float64) {
	gravityConst = G
	warpGravityConst = W
}

//
func Gravity(mass, lenSqr, zDist float64) float64 {
	d2 := lenSqr + zDist*zDist
	d2 = d2 * d2
	if d2 == 0 {
		return 0
	}
	return gravityConst * mass * lenSqr / d2
}

func WarpGravity(mass, lenSqr, velSqr, zDist float64) float64 {

	d2 := lenSqr + zDist*zDist
	d2 = d2 * d2
	if d2 == 0 {
		return 0
	}

	return warpGravityConst * mass * lenSqr / d2 * (1 + velSqr)
}
