package graph

import "github.com/hajimehoshi/ebiten"

func CircleTex(radius int) Tex {
	d := radius*2 + 1
	img, _ := ebiten.NewImage(d, d, ebiten.FilterDefault)

	p := make([]byte, d*d*4)
	dw := d * 4
	r2 := radius * radius

	for x := -d; x <= d; x++ {
		for y := -d; y <= d; y++ {
			if x*x+y*y <= r2 {
				ix := x + d
				iy := y + d
				for i := 0; i < 4; i++ {
					p[4*ix+dw*iy+i] = 255
				}
			}
		}
	}

	return TexFromImage(img, ebiten.FilterDefault, 0, 0, 0)
}
