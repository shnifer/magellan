package draw

import "github.com/Shnifer/magellan/graph"

const (
	Mark_size   = 25
	sprite_size = 40
)

func MarkAlpha(size float64, cam *graph.Camera) (alphaMark, alphaSprite float64) {
	if cam != nil {
		size *= cam.Scale
	}
	if size <= Mark_size {
		alphaMark = 1
		alphaSprite = 0
	} else if size >= sprite_size {
		alphaMark = 0
		alphaSprite = 1
	} else {
		k := (size - Mark_size) / (sprite_size - Mark_size)
		alphaMark = 1 - k
		alphaSprite = k
	}
	return alphaMark, alphaSprite
}

func MarkScaleLevel(level int) float64 {
	if level <= 1 {
		return 1.0 / 1
	} else if level == 2 {
		return 1.0 / 2
	} else if level == 3 {
		return 1.0 / 3
	} else {
		return 1.0 / 4
	}

}
