package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/v2"
)

func (s *cosmoScene) updateShipControl(dt float64) {
	s.procControlTurn(dt)
	s.procControlForward(dt)
}

func (s *cosmoScene) procControlTurn(dt float64) {
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

func (s *cosmoScene) procControlForward(dt float64) {
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

func (s *cosmoScene) procShipGravity(dt float64) {
	var sumF v2.V2
	for _, obj := range s.objects {
		v := obj.Pos.Sub(Data.PilotData.Ship.Pos)
		len2 := v.LenSqr()
		F := Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), F)
	}
	Data.PilotData.Ship.Vel.DoAddMul(sumF, dt)
}
