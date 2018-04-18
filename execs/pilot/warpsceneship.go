package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
)

func (s *warpScene) updateShipControl(dt float64) {
	s.procControlTurn(dt)
	s.procControlForward(dt)
}

func (s *warpScene) procControlTurn(dt float64) {
	turnInput := input.GetF("turn")
	var min, max float64
	switch {
	case s.maneurLevel >= 0:
		max = s.maneurLevel + Data.BSP.Turn_acc/100*dt
		min = s.maneurLevel - Data.BSP.Turn_slow/100*dt
	case s.maneurLevel < 0:
		max = s.maneurLevel + Data.BSP.Turn_slow/100*dt
		min = s.maneurLevel - Data.BSP.Turn_acc/100*dt
	}
	s.maneurLevel = Clamp(turnInput, min, max)
	Data.PilotData.Ship.AngVel = s.maneurLevel * Data.BSP.Turn_max
}

func (s *warpScene) procControlForward(dt float64) {
	thrustInput := input.GetF("forward")
	var min, max float64
	switch {
	case s.thrustLevel >= 0:
		max = s.thrustLevel + Data.BSP.Thrust_acc/100*dt
		min = s.thrustLevel - Data.BSP.Thrust_slow/100*dt
	case s.thrustLevel < 0:
		max = s.thrustLevel + Data.BSP.Thrust_rev_slow/100*dt
		min = s.thrustLevel - Data.BSP.Thrust_rev_acc/100*dt
	}
	if Data.BSP.Thrust_rev == 0 && min < 0 {
		min = 0
	}
	s.thrustLevel = Clamp(thrustInput, min, max)

	var accel float64
	switch {
	case s.thrustLevel >= 0:
		accel = s.thrustLevel * Data.BSP.Thrust
	case s.thrustLevel < 0:
		accel = s.thrustLevel * Data.BSP.Thrust_rev
	}
	Data.PilotData.Ship.Vel.DoAddMul(v2.InDir(Data.PilotData.Ship.Ang), accel*dt)
}

func (s *warpScene) procShipGravity(dt float64) {
	var F v2.V2
	for _, obj := range s.objects {
		V := obj.Pos.Sub(Data.PilotData.Ship.Pos)
		D2 := V.LenSqr() + DEFVAL.GravityZ2
		F = F.Add(V.Normed().Mul(obj.Mass * DEFVAL.GravityConst / D2))
	}
	Data.PilotData.Ship.Vel.DoAddMul(F, dt)
}
