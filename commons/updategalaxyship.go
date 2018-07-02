package commons

import (
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
)

//Used by Pilot to carefully calculate gravity
//Other clients should use Galaxy.Update(SessionTime) and Ship RB predictor

//it adds dt to data.PilotData.SessionTime
//it changes ship.Ang by ship.AngVel*sumT
//it applies PilotData.ThrustVector and Gravity each littleT
//rest of sumT is stored and data.Galaxy.fixedTimeRest
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

	sumT += galaxy.fixedTimeRest
	for sumT >= dt {
		sessionTime += dt
		//todo: DO THIS FASTER, Bench ready.
		galaxy.Update(sessionTime)
		sumT -= dt

		grav = SumGravityAcc(ship.Pos, galaxy)
		thrust = v2.InDir(ship.Ang).Mul(thrustF)
		ship.Vel.DoAddMul(v2.Add(grav, thrust), dt)
		ship.Pos.DoAddMul(ship.Vel, dt)
		ship.Ang += ship.AngVel * dt
	}

	data.PilotData.Ship = ship
	data.Galaxy.fixedTimeRest = sumT
	data.PilotData.SessionTime = sessionTime
}
