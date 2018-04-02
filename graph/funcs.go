package graph

var winH, winW float64

func SetScreenSize(W, H int) {
	winW = float64(W)
	winH = float64(H)
}

func ScrP(x, y float64) Point {
	return Point{x * winW, y * winH}
}
