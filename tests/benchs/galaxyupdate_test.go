package main

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"io/ioutil"
	"testing"
)

var galaxy *Galaxy

func BenchmarkUpdate1(b *testing.B) {
	loadgalaxy()
	var st float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st += 0.1
		Update1(st)
	}
}

func BenchmarkUpdate2(b *testing.B) {
	loadgalaxy()
	var st float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st += 0.1
		Update2(st)
	}
}

func BenchmarkUpdate3(b *testing.B) {
	loadgalaxy()
	Data := TData{}
	Data.Galaxy = galaxy
	Data.PilotData = &PilotData{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Update3(Data, 0.016, 0.001)
	}
}

func Update1(st float64) {
	for _, obj := range galaxy.Ordered {
		if obj.ParentID == "" {
			continue
		}
		parent := galaxy.Points[obj.ParentID].Pos
		angle := (360 / obj.Period) * st
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}

func Update2(st float64) {
	posMap := make(map[string]v2.V2)

	for _, obj := range galaxy.Ordered {
		if obj.ParentID == "" {
			continue
		}
		parent, ok := posMap[obj.ParentID]
		if !ok {
			parent = galaxy.Points[obj.ParentID].Pos
			posMap[obj.ParentID] = parent
		}
		angle := (360 / obj.Period) * st
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}

var fixedTimeRest float64

func Update3(data TData, sumT float64, dt float64) {

	galaxy := data.Galaxy
	ship := data.PilotData.Ship
	sessionTime := data.PilotData.SessionTime
	thrustF := data.PilotData.ThrustVector.Len()

	var thrust v2.V2
	var grav v2.V2

	moveList := make(map[string]struct{})
	var l2, g float64
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		l2 = ship.Pos.Sub(obj.Pos).LenSqr()
		g = obj.Mass / l2
		if g > 0.001 {
			moveList[obj.ID] = struct{}{}
		}
	}

	sumT += fixedTimeRest
	for sumT >= dt {
		sessionTime += dt
		UpdateGal(galaxy, sessionTime, moveList)
		sumT -= dt

		grav = SumGravityAcc(ship.Pos, galaxy)
		thrust = v2.InDir(ship.Ang).Mul(thrustF)
		ship.Vel.DoAddMul(v2.Add(grav, thrust), dt)
		ship.Pos.DoAddMul(ship.Vel, dt)
		ship.Ang += ship.AngVel * dt
	}

	data.PilotData.Ship = ship
	fixedTimeRest = sumT
	data.PilotData.SessionTime = sessionTime
}

func UpdateGal(galaxy *Galaxy, sessionTime float64, moveList map[string]struct{}) {
	if galaxy == nil {
		return
	}

	//bench tells that this way is faster
	var parent v2.V2

	//skip lvl 0 objects, they do not move
	for id := range moveList {
		obj := galaxy.Points[id]
		if obj.ParentID == "" {
			continue
		}
		parent = galaxy.Points[obj.ParentID].Pos

		angle := (360 / obj.Period) * sessionTime
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}

var gravGalaxy gravGalaxyT

func BenchmarkUpdate4(b *testing.B) {
	loadgalaxy()
	Data := TData{}
	Data.Galaxy = galaxy
	Data.PilotData = &PilotData{}
	gravGalaxy = calcGravGalaxy(galaxy)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Update4(Data, 0.016, 0.001)
	}
}

func Update4(data TData, sumT float64, dt float64) {

	galaxy := data.Galaxy
	ship := data.PilotData.Ship
	sessionTime := data.PilotData.SessionTime
	thrust := data.PilotData.ThrustVector

	var grav v2.V2

	loadGravGalaxy(gravGalaxy, galaxy)

	sumT += fixedTimeRest
	for sumT >= dt {
		sessionTime += dt
		UpdateGravGal(gravGalaxy, sessionTime)
		sumT -= dt

		grav = sumGrav(ship.Pos, gravGalaxy)
		ship.Vel.DoAddMul(v2.Add(grav, thrust), dt)
		ship.Pos.DoAddMul(ship.Vel, dt)
		ship.Ang += ship.AngVel * dt
	}

	data.PilotData.Ship = ship
	fixedTimeRest = sumT
	data.PilotData.SessionTime = sessionTime
}

type gravP = struct {
	id        string
	parentInd int
	pos       v2.V2
	orbit     float64
	period    float64
	mass      float64
	gDepth    float64
}
type gravGalaxyT = []gravP

func UpdateGravGal(galaxy gravGalaxyT, sessionTime float64) {
	if galaxy == nil {
		return
	}

	var parent v2.V2

	//skip lvl 0 objects, they do not move
	for i, gravP := range galaxy {
		if gravP.parentInd == -1 {
			continue
		}

		parent = galaxy[gravP.parentInd].pos

		angle := (360 / gravP.period) * sessionTime
		galaxy[i].pos = parent.AddMul(v2.InDir(angle), gravP.orbit)
	}
}

func calcGravGalaxy(galaxy *Galaxy) gravGalaxyT {
	res := make(gravGalaxyT, 0, len(galaxy.Ordered))

	ord := make(map[string]int)
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		p := gravP{
			id:     obj.ID,
			pos:    obj.Pos,
			orbit:  obj.Orbit,
			period: obj.Period,
			mass:   obj.Mass,
			gDepth: obj.GDepth,
		}
		if obj.ParentID == "" {
			p.parentInd = -1
		} else {
			p.parentInd = ord[obj.ParentID]
		}
		res = append(res, p)
		ord[obj.ID] = len(res) - 1
	}

	return res
}

func loadGravGalaxy(gg gravGalaxyT, galaxy *Galaxy) {
	for i, p := range gg {
		gg[i].pos = galaxy.Points[p.id].Pos
	}
}

func sumGrav(pos v2.V2, gg gravGalaxyT) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range gg {
		v = obj.pos.Sub(pos)
		len2 = v.LenSqr()
		G = Gravity(obj.mass, len2, obj.gDepth)
		sumF.DoAddMul(v.Normed(), G)
	}
	return sumF
}

func loadgalaxy() {
	buf, err := ioutil.ReadFile("testgalaxy.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &galaxy)
	if err != nil {
		panic(err)
	}

	//First restore ID's
	for id, v := range galaxy.Points {
		if v.ID == "" {
			v.ID = id
			galaxy.Points[id] = v
		}
	}
	//Second - recalc lvls!
	galaxy.RecalcLvls()
}
