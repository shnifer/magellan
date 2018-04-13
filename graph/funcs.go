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
