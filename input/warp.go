package input

const warpLevelThreshold = 0.25
const historyLen = 6

var warpInputHistory []float64

func WarpLevel(inputname string) float64 {
	v := GetF(inputname)

	warpInputHistory = append(warpInputHistory, v)
	cut := len(warpInputHistory) - historyLen
	if cut > 0 {
		warpInputHistory = warpInputHistory[cut:]
		var s float64
		for _, v := range warpInputHistory {
			s += v
		}
		v = s / historyLen
	} else {
		v = 0
	}

	v = (v + 1) / 2
	if v < warpLevelThreshold {
		return 0
	}
	if v > 1 {
		return 1
	}
	return (v - warpLevelThreshold) / (1 - warpLevelThreshold)
}
