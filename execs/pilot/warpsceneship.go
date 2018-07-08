package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
)

func (s *warpScene) updateShipControl(dt float64) {
	s.procControlForward(dt)
	s.procControlTurn(dt)
}

func (s *warpScene) procControlForward(dt float64) {
	thrustInput := input.GetF("forward")

	switch {
	case thrustInput >= 0:
		s.thrustLevel = s.thrustLevel + Data.SP.Warp_engine.Distort_acc/100*thrustInput*dt
	case thrustInput < 0:
		s.thrustLevel = s.thrustLevel + Data.SP.Warp_engine.Distort_slow/100*thrustInput*dt
	}

	s.thrustLevel = Clamp(s.thrustLevel, 0, 1)

	Data.PilotData.Ship.Vel = v2.InDir(Data.PilotData.Ship.Ang).Mul(s.thrustLevel * Data.SP.Warp_engine.Distort_max)
}

func (s *warpScene) procControlTurn(dt float64) {
	turnInput := input.GetF("turn")
	Data.PilotData.Ship.AngVel = turnInput * Data.SP.Warp_engine.Turn_speed * s.thrustLevel *60
}

func (s *warpScene) procShipGravity(dt float64) {
	var sumV v2.V2
	for _, obj := range Data.Galaxy.Points {
		v := obj.Pos.Sub(Data.PilotData.Ship.Pos)
		len2 := v.LenSqr()
		vel := Data.PilotData.Ship.Vel.LenSqr()
		F := WarpGravity(obj.Mass, len2, vel, obj.Size/2)
		sumV.DoAddMul(v.Normed(), F)
	}
	//Data.PilotData.Ship.Vel.DoAddMul(sumV, dt)
}
