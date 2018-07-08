package commons

import (
	"github.com/Shnifer/magellan/v2"
	."github.com/Shnifer/magellan/log"
)

//Vel = Distortion*VelDistK
const VelDistK = 0.5
//Acc = gravAcc*Distortion^3*AccDistK
const AccDistK = 0.1

func UpdateWarpAndShip(data TData, sumT float64, dt float64) {
	if data.Galaxy == nil {
		Log(LVL_ERROR, "UpdateWarpAndShip called with nil Galaxy")
		return
	}
	if data.PilotData == nil {
		Log(LVL_ERROR, "UpdateWarpAndShip called with nil PilotData")
		return
	}

	galaxy := data.Galaxy
	ship := data.PilotData.Ship
	sessionTime := data.PilotData.SessionTime
	distortion := data.PilotData.Distortion
	gravK := distortion * distortion * distortion * AccDistK
	vel := VelDistK * distortion

	var grav v2.V2

	sumT += galaxy.fixedTimeRest
	for sumT >= dt {
		sessionTime += dt
		sumT -= dt

		grav = SumWarpGravityAcc(ship.Pos, galaxy).Mul(gravK)
		ship.Vel.DoAddMul(grav,dt)
		ship.Vel = ship.Vel.Normed().Mul(vel)
		ship.Pos.DoAddMul(ship.Vel, dt)
	}
	ship.Ang = ship.Vel.Dir()

	//final update for all and every object, slow but once
	data.Galaxy.Update(sessionTime)

	data.PilotData.Ship = ship
	data.Galaxy.fixedTimeRest = sumT
	data.PilotData.SessionTime = sessionTime
}

