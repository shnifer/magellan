package graph

import (
	"github.com/Shnifer/magellan/v2"
)

var winH, winW float64

func SetScreenSize(W, H int) {
	winW = float64(W)
	winH = float64(H)
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
