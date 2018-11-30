package commons

import (
	. "github.com/shnifer/magellan/log"
	"github.com/shnifer/magellan/v2"
)

//Used by Pilot to carefully calculate gravity
//Other clients should use Galaxy.Update(SessionTime) and Ship RB predictor
//Calls Galaxy.Update(SessionTime) at the and, so no need to call it again

//it adds dt to data.PilotData.SessionTime
//it changes ship.Ang by ship.AngVel*sumT
//it applies PilotData.ThrustVector and Gravity each littleT
//rest of sumT is stored and data.Galaxy.fixedTimeRest

const minGravityToMove = 0.01

func UpdateGalaxyAndShip(data TData, sumT float64, dt float64) {
	if data.Galaxy == nil {
		Log(LVL_ERROR, "UpdateGalaxyAndShip called with nil Galaxy")
		return
	}
	if data.PilotData == nil {
		Log(LVL_ERROR, "UpdateGalaxyAndShip called with nil PilotData")
		return
	}

	galaxy := data.Galaxy
	ship := data.PilotData.Ship
	sessionTime := data.PilotData.SessionTime
	thrustF := data.PilotData.ThrustVector.Len()

	var thrust v2.V2
	var grav v2.V2

	moveList := calcMoveList(galaxy, ship.Pos, minGravityToMove)

	sumT += galaxy.fixedTimeRest
	for sumT >= dt {
		sessionTime += dt
		moveMasses(galaxy, sessionTime, moveList)
		sumT -= dt

		grav = SumGravityAcc(ship.Pos, galaxy)
		thrust = v2.InDir(ship.Ang).Mul(thrustF)
		ship.Vel.DoAddMul(v2.Add(grav, thrust), dt)
		ship.Pos.DoAddMul(ship.Vel, dt)
		ship.Ang += ship.AngVel * dt
	}

	//final update for all and every object, slow but once
	data.Galaxy.Update(sessionTime)

	data.PilotData.Ship = ship
	data.Galaxy.fixedTimeRest = sumT
	data.PilotData.SessionTime = sessionTime
}

func calcMoveList(galaxy *Galaxy, shipPos v2.V2, minLevel float64) map[string]struct{} {
	moveList := make(map[string]struct{})
	var l2, g float64
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		l2 = shipPos.Sub(obj.Pos).LenSqr()
		g = obj.Mass / l2
		if g >= minLevel {
			moveList[obj.ID] = struct{}{}
		}
	}
	return moveList
}

func moveMasses(galaxy *Galaxy, sessionTime float64, moveList map[string]struct{}) {

	//do not use parent position map here, cz moveList must not be large
	var parent v2.V2

	//skip lvl 0 objects, they do not move
	for id := range moveList {
		obj := galaxy.Points[id]
		if obj.ParentID == "" {
			continue
		}

		parent = galaxy.Points[obj.ParentID].Pos
		angle := (360/obj.Period)*sessionTime + obj.AngPhase
		obj.Pos = parent.AddMul(v2.InDir(angle), obj.Orbit)
	}
}
