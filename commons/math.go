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
