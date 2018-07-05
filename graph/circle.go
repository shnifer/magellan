package graph

import "github.com/hajimehoshi/ebiten"

const defaultTexRadius = 512

var defaultCircle Tex

func init() {
	defaultCircle = circleTex(defaultTexRadius)
}

func CircleTex() Tex {
	return defaultCircle
}

func circleTex(radius int) Tex {
	d := radius*2 + 1
	img, _ := ebiten.NewImage(d, d, ebiten.FilterLinear)

	p := make([]byte, d*d*4)
	dw := d * 4
	r2 := radius * radius

	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			if x*x+y*y <= r2 {
				ix := x + radius
				iy := y + radius
				for i := 0; i < 4; i++ {
					p[4*ix+dw*iy+i] = 255
				}
			}
		}
	}

	img.ReplacePixels(p)

	return TexFromImage(img, ebiten.FilterLinear, 0, 0, 0, "~circle")
}
