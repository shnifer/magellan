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
	//map[id]level
	lvls := make(map[string]int)
	glvls := make(map[string]int)

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
	var GLvl func(*GalaxyPoint) int
	var L int
	GLvl = func(p *GalaxyPoint) int {
		parent := p.ParentID
		if parent == "" {
			glvls[p.ID] = 0
			return 0
		}
		l, ok := glvls[parent]
		if ok {
			L = l + 1
			if galaxy.Points[parent].IsVirtual {
				L = l
			}
			glvls[p.ID] = L
			return L
		}

		l = GLvl(galaxy.Points[parent])
		L = l + 1
		if galaxy.Points[parent].IsVirtual {
			L = l
		}
		glvls[p.ID] = L
		return L
	}

	var lvl, glvl int
	for id, point := range galaxy.Points {
		lvl = Lvl(point)
		glvl = GLvl(point)
		galaxy.Points[id].Level = lvl
		galaxy.Points[id].GLevel = glvl
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
	for id, v := range galaxy.Ordered {
		if v.Mines == nil {
			v.Mines = make(map[string][]string)
		}
		if v.FishHouses == nil {
			v.FishHouses = make(map[string]string)
		}
		if v.Beacons == nil {
			v.Beacons = make(map[string]string)
		}
		if v.BlackBoxes == nil {
			v.BlackBoxes = make(map[string]string)
		}
		if v.Emissions == nil {
			v.Emissions = make([]Emission, 0)
		}
		if v.Signatures == nil {
			v.Signatures = make([]Signature, 0)
		}
		galaxy.Ordered[id] = v
	}
}

//must be used once per frame, to recalc all positions
func (galaxy *Galaxy) Update(sessionTime float64) {
	defer LogFunc("galaxy.Update")()

	if galaxy == nil {
		return
	}

	//bench tells that this way is faster
	posMap := make(map[string]v2.V2)
	//skip lvl 0 objects, they do not move
	for _, obj := range galaxy.Ordered {
		if obj.ParentID == "" {
			continue
		}
		parent, ok := posMap[obj.ParentID]
		if !ok {
			parent = galaxy.Points[obj.ParentID].Pos
			posMap[obj.ParentID] = parent
		}
		angle := (360 / obj.Period) * sessionTime
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}

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
		if _, exist := gp.Mines[b.OwnerID]; exist {
			for _,v:=range gp.Mines[b.OwnerID]{
				if v==fullKey{
					Log(LVL_ERROR, "trying to add mine on planet ", b.PlanetID, " but already has mine with fk ",fullKey)
					return
				}
			}
		}
		gp.Mines[b.OwnerID] = append(gp.Mines[b.OwnerID],fullKey)
	case BUILDING_FISHHOUSE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add fishhouse on non existant planet with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.FishHouses[b.OwnerID]; exist {
			Log(LVL_ERROR, "trying to add fishhouse on planet ", b.PlanetID, " but already has fishhouse")
			return
		}
		gp.FishHouses[b.OwnerID] = fullKey
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

//works with already calced and ordered Galaxy
func (galaxy *Galaxy) AddWarpBuilding(b Building) {
	fullKey := b.FullKey
	if _, exist := galaxy.Points[fullKey]; exist {
		//already exist
		Log(LVL_WARN, "trying to add building with already exist Fullkey:", fullKey)
		return
	}

	switch b.Type {
	case BUILDING_BEACON:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add beacon on non existant star with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.Beacons[fullKey]; exist {
			Log(LVL_ERROR, "trying to add beacon on star ", b.PlanetID, " but already has this beacon")
			return
		}
		gp.Beacons[fullKey] = b.Message
	case BUILDING_BLACKBOX:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add blackbox on non existant star with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.BlackBoxes[fullKey]; exist {
			Log(LVL_ERROR, "trying to add blackbox on star ", b.PlanetID, " but already has this blackbox")
			return
		}
		gp.BlackBoxes[fullKey] = b.Message
	default:
		Log(LVL_ERROR, "warpscene addBuilding, unknown building type", b.Type)
	}
}

//works with already calced and ordered Galaxy
func (galaxy *Galaxy) DelBuilding(b Building) {
	switch b.Type {
	case BUILDING_MINE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add del on non existant planet with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.Mines[b.OwnerID]; !exist {
			Log(LVL_ERROR, "trying to del mine on planet", b.PlanetID, "but do not has mine")
			return
		}
		fks:=gp.Mines[b.OwnerID]
		found:=false
		for i,v:=range fks{
			if v==b.FullKey{
				fks = append(fks[:i], fks[i+1:]...)
				found = true
				break
			}
		}
		if !found{
			Log(LVL_ERROR, "trying to del mine on planet", b.PlanetID,
				"but it do not has mine with fk ", b.FullKey)
			return
		}
		if len(fks)==0 {
			delete(gp.Mines, b.OwnerID)
		} else {
			gp.Mines[b.OwnerID] = fks
		}
	case BUILDING_FISHHOUSE:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to add fishhouse on non existant planet with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.FishHouses[b.OwnerID]; !exist {
			Log(LVL_ERROR, "trying to del fishhouse on planet", b.PlanetID, "but do not has mine")
			return
		}
		delete(gp.FishHouses, b.OwnerID)
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

//works with already calced and ordered Galaxy
func (galaxy *Galaxy) DelWarpBuilding(b Building) {
	fullKey := b.FullKey
	switch b.Type {
	case BUILDING_BEACON:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to del beacon on non existant planet with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.Beacons[fullKey]; !exist {
			Log(LVL_ERROR, "trying to del beacon on planet", b.PlanetID, "but do not has mine")
			return
		}
		delete(gp.Beacons, fullKey)
	case BUILDING_BLACKBOX:
		gp, ok := galaxy.Points[b.PlanetID]
		if !ok {
			Log(LVL_ERROR, "trying to del blackbox on non existant planet with ID:", b.PlanetID)
			return
		}
		if _, exist := gp.BlackBoxes[fullKey]; !exist {
			Log(LVL_ERROR, "trying to del blackbox on planet", b.PlanetID, "but do not has mine")
			return
		}
		delete(gp.BlackBoxes, fullKey)
	default:
		Log(LVL_ERROR, "warp delBuilding, unknown building type", b.Type)
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
