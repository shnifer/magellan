package commons

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
)

//fixed
const BuildingSize = 20

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

	galaxy.Ordered = galaxy.Ordered[:0]
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

//works with already calced and ordered Galaxy
func (galaxy *Galaxy) addBuilding(b Building) {
	fullKey := b.FullKey
	if _, exist := galaxy.Points[fullKey]; exist {
		//already exist
		Log(LVL_WARN, "trying to add building with already exist Fullkey:", fullKey)
		return
	}

	switch b.Type {
	case BUILDING_MINE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add mine on non existant planet with ID:", b.PlanetID)
			return
		}
		if gp.HasMine {
			Log(LVL_ERROR, "trying to add mine on planet", b.PlanetID, " but already has mine")
			return
		}
		gp.HasMine = true
		gp.MineOwner = b.OwnerID
	case BUILDING_BEACON, BUILDING_BLACKBOX:
		parentID := ""
		if len(galaxy.Ordered) > 0 {
			parentID = galaxy.Ordered[0].ID
		}
		gp := GalaxyPoint{}.outerFromBuilding(b, parentID, galaxy.SpawnDistance)
		galaxy.Points[fullKey] = gp
		galaxy.Ordered = append(galaxy.Ordered, gp)
	}
}
func (GalaxyPoint) outerFromBuilding(b Building, parentID string, dist float64) (gp *GalaxyPoint) {
	gp = &GalaxyPoint{
		ID:         b.FullKey,
		ParentID:   parentID,
		Type:       b.Type,
		Orbit:      dist,
		Period:     b.Period,
		ScanData:   b.Message,
		Size:       BuildingSize,
		Emissions:  []Emission{},
		Signatures: []Signature{},
	}

	if b.Message != "" {
		gp.Signatures = append(gp.Signatures, Signature{SigString: b.Message})
	}

	return gp
}
