package input

const warpLevelThreshold = 0.1

func WarpLevel(inputname string) float64 {
	v := GetF(inputname)
	v = (v + 1) / 2
	if v < warpLevelThreshold {
		return 0
	}
	if v > 1 {
		return 1
	}
	return (v - warpLevelThreshold) / (1 - warpLevelThreshold)
}
