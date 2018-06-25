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

	for id, point := range galaxy.Points {
		lvl := Lvl(point)
		galaxy.Points[id].Level = lvl
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

//TODO: rework for multiple mines and houses
//works with already calced and ordered Galaxy
func (galaxy *Galaxy) AddBuilding(b Building) {
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
			Log(LVL_ERROR, "trying to add mine on planet ", b.PlanetID, " but already has mine")
			return
		}
		gp.HasMine = true
		gp.MineOwner = b.OwnerID
		gp.MineFullKey = b.FullKey
	case BUILDING_FISHHOUSE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add fishhouse on non existant planet with ID:", b.PlanetID)
			return
		}
		if gp.HasFishHouse {
			Log(LVL_ERROR, "trying to add fishhouse on planet", b.PlanetID, " but already has mine")
			return
		}
		gp.HasFishHouse = true
		gp.FishHouseOwner = b.OwnerID
		gp.FishHouseFullKey = b.FullKey
	case BUILDING_BEACON, BUILDING_BLACKBOX:
		parentID := ""
		if len(galaxy.Ordered) > 0 {
			parentID = galaxy.Ordered[0].ID
		}
		gp := GalaxyPoint{}.outerFromBuilding(b, parentID, galaxy.SpawnDistance)
		//append to the end without resorting cz buildings has no child
		galaxy.Points[fullKey] = gp
		galaxy.Ordered = append(galaxy.Ordered, gp)
	default:
		Log(LVL_ERROR, "cosmoscene addBuilding, unknown building type", b.Type)
	}
}

//TODO: rework for multiple mines and houses
//works with already calced and ordered Galaxy
func (galaxy *Galaxy) DelBuilding(b Building) {
	switch b.Type {
	case BUILDING_MINE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add del on non existant planet with ID:", b.PlanetID)
			return
		}
		if !gp.HasMine {
			Log(LVL_ERROR, "trying to del mine on planet", b.PlanetID, "but do not has mine")
			return
		}
		gp.HasMine = false
		gp.MineOwner = ""
		gp.MineFullKey = ""
	case BUILDING_FISHHOUSE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add fishhouse on non existant planet with ID:", b.PlanetID)
			return
		}
		if !gp.HasFishHouse {
			Log(LVL_ERROR, "trying to del fishhouse on planet", b.PlanetID, "but do not has mine")
			return
		}
		gp.HasFishHouse = false
		gp.FishHouseOwner = ""
		gp.FishHouseFullKey = ""
	case BUILDING_BEACON, BUILDING_BLACKBOX:
		fullKey := b.FullKey
		pointer, exist := galaxy.Points[fullKey]
		if !exist {
			Log(LVL_WARN, "trying to del building but can't find full key:", fullKey)
			return
		}
		for i, v := range galaxy.Ordered {
			if v == pointer {
				galaxy.Ordered = append(galaxy.Ordered[:i], galaxy.Ordered[i+1:]...)
				break
			}
		}
		delete(galaxy.Points, fullKey)
	default:
		Log(LVL_ERROR, "galaxy delBuilding, unknown building type", b.Type)
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
		Level:      1,
	}

	if b.Message != "" {
		gp.Signatures = append(gp.Signatures, Signature{SigString: b.Message})
	}

	return gp
}
