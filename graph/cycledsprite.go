package graph

const (
	Cycle_OneTime int = iota
	Cycle_Loop
	Cycle_PingPong
)

type CycledSprite struct {
	Sprite

	cycleType int
	curT      float64
	periodT   float64

	//for pingpong
	dir int

	//spriterange
	min, max int

	paused bool
}

func NewCycledSprite(sprite *Sprite, cycleType int, fps float64) *CycledSprite {
	return NewCycledSpriteRange(sprite, cycleType, fps, 0, sprite.SpritesCount()-1)
}

func NewCycledSpriteRange(sprite *Sprite, cycleType int, fps float64, min, max int) *CycledSprite {
	if fps == 0 {
		panic("zero fps!")
	}
	if min > max {
		panic("NewCycledSpriteRange: max < min")
	}
	if min < 0 {
		panic("NewCycledSpriteRange: min < 0")
	}
	if max > sprite.SpritesCount()-1 {
		panic("NewCycledSpriteRange: max>SpritesCount()-1")
	}

	var s Sprite
	s = *sprite
	return &CycledSprite{
		Sprite:    s,
		cycleType: cycleType,
		curT:      0,
		periodT:   1 / fps,
		dir:       1,
		min:       min,
		max:       max,
	}
}

func (cs *CycledSprite) Reset() {
	cs.curT = 0
	cs.dir = 1
	cs.paused = false
}

func (cs *CycledSprite) SetPause(paused bool) {
	cs.paused = paused
}

func (cs *CycledSprite) Update(dt float64) {
	if !cs.paused {
		cs.curT += dt
		for cs.curT >= cs.periodT {
			cs.curT -= cs.periodT
			cs.nextSprite()
		}
	}
}

func (cs *CycledSprite) nextSprite() {
	n := cs.SpriteN()
	n += cs.dir
	switch cs.cycleType {
	case Cycle_OneTime:
		if n > cs.max {
			n = cs.max
		}
	case Cycle_Loop:
		if n > cs.max {
			n = cs.min
		}
	case Cycle_PingPong:
		if n > cs.max {
			n = cs.max - 1
			cs.dir = -1
		}
		if n < cs.min {
			n = cs.min + 1
			cs.dir = 1
		}
	}
	if n < cs.min {
		n = cs.min
	} else if n > cs.max {
		n = cs.max
	}
	cs.SetSpriteN(n)
}
