package graph

const (
	Cycle_OneTime int = iota
	Cycle_Loop
	Cycle_PingPong
)

type CycledSprite struct {
	*Sprite

	cycleType int
	curT      float64
	periodT   float64

	//for pingpong
	dir int

	paused bool
}

func NewCycledSprite(sprite *Sprite, cycleType int, fps float64) *CycledSprite {
	if fps == 0 {
		panic("zero fps!")
	}

	return &CycledSprite{
		Sprite:    sprite,
		cycleType: cycleType,
		curT:      0,
		periodT:   1 / fps,
		dir:       1,
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
	max := cs.SpritesCount() - 1
	n += cs.dir
	switch cs.cycleType {
	case Cycle_OneTime:
		if n > max {
			n = max
		}
	case Cycle_Loop:
		if n > max {
			n = 0
		}
	case Cycle_PingPong:
		if n > max {
			n = max - 1
			cs.dir = -1
		}
		if n < 0 {
			n = 1
			cs.dir = 1
		}
	}
	if n < 0 {
		n = 0
	} else if n > max {
		n = max
	}
	cs.SetSpriteN(n)
}
