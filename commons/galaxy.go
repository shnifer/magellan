package commons

import (
	"github.com/Shnifer/magellan/v2"
)

func (galaxy *Galaxy) RecalcLvls() {
	defer LogFunc("galaxy.RecalcLvls")()

	if galaxy == nil {
		return
	}

	maxLvl := 0
	lvls := make(map[string]int)

	var Lvl func(*GalaxyPoint) int
	Lvl = func(p *GalaxyPoint) int {
		parent := p.ParentID
		if parent == "" {
			lvls[p.ID] = 0
			return 0
		}
		l, ok := lvls[parent]
		if ok {
			lvls[p.ID] = l + 1
			return l + 1
		}

		l = Lvl(galaxy.Points[parent])
		lvls[p.ID] = l + 1
		return l + 1
	}

	for _, point := range galaxy.Points {
		lvl := Lvl(point)
		if lvl > maxLvl {
			maxLvl = lvl
		}
	}

	galaxy.maxLvl = maxLvl
	for lvl := 0; lvl <= maxLvl; lvl++ {
		for id, p := range galaxy.Points {
			if lvls[id] == lvl {
				galaxy.Ordered = append(galaxy.Ordered, p)
			}
		}
	}
}

func (galaxy *Galaxy) Update(sessionTime float64) {
	defer LogFunc("galaxy.Update")()

	if galaxy == nil {
		return
	}
	//skip lvl 0 objects, they do not move
	for _, obj := range galaxy.Ordered {
		if obj.ParentID == "" {
			continue
		}
		parent := galaxy.Points[obj.ParentID].Pos
		angle := (360 / obj.Period) * sessionTime
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}
