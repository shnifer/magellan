package graph

import (
	"github.com/Shnifer/magellan/v2"
)

var winH, winW float64
var globalScale float64

func init() {
	globalScale = 1
}

func SetScreenSize(W, H int) {
	winW = float64(W)
	winH = float64(H)

	globalScale = calcGlobalScale(H)
}

func calcGlobalScale(H int) float64 {
	bounds := []float64{0.75, 0.8, 0.9, 1.0, 1.1, 1.2, 1.25, 1.333333, 1.4, 1.5, 1.6, 1.666666,
		1.75, 1.8, 2.0, 2.25, 2.5, 2.75, 3.0, 4.0, 5.0}
	res := float64(H) / 1000
	l := len(bounds)
	if res <= bounds[0] {
		return bounds[0]
	}
	if res >= bounds[l-1] {
		return bounds[l-1]
	}
	for i, v := range bounds {
		if i == 0 || res > v {
			continue
		}
		if res/bounds[i-1] < v/res {
			return bounds[i-1]
		} else {
			return v
		}
	}
	return bounds[l-1]
}

func ScrP(x, y float64) v2.V2 {
	return v2.V2{X: x * winW, Y: y * winH}
}

//put angle in degs in [0;360) range
func NormAng(angle float64) float64 {
	for angle < 0 {
		angle += 360
	}
	for angle >= 360 {
		angle -= 360
	}
	return angle
}

//normalize start angle in [0;360) and end in [start; start+360)
//so always end > start. End value itself may be more than 360
func NormAngRange(start, end float64) (float64, float64) {
	if start > end {
		start, end = end, start
	}

	for start < 0 {
		start += 360
	}
	for start >= 360 {
		start -= 360
	}
	for end < start {
		end += 360
	}
	for end >= start+360 {
		end -= 360
	}
	return start, end
}

//global scale factor
func GS() float64 {
	return globalScale
}
