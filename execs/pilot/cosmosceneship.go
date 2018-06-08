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
	massK:=1000/Data.BSP.Mass
	var min, max float64
	switch {
	case s.maneurLevel >= 0:
		max = s.maneurLevel + Data.SP.Turn_acc*massK/100*dt
		min = s.maneurLevel - Data.SP.Turn_slow*massK/100*dt
	case s.maneurLevel < 0:
		max = s.maneurLevel + Data.SP.Turn_slow*massK/100*dt
		min = s.maneurLevel - Data.SP.Turn_acc*massK/100*dt
	}
	s.maneurLevel = Clamp(turnInput, min, max)
	Data.PilotData.Ship.AngVel = s.maneurLevel * Data.SP.Turn_max
}

func (s *cosmoScene) procControlForward(dt float64) {
	thrustInput := input.GetF("forward")
	massK:=1000/Data.BSP.Mass
	var min, max float64
	switch {
	case s.thrustLevel >= 0:
		max = s.thrustLevel + Data.SP.Thrust_acc*massK/100*dt
		min = s.thrustLevel - Data.SP.Thrust_slow*massK/100*dt
	case s.thrustLevel < 0:
		max = s.thrustLevel + Data.SP.Thrust_rev_slow*massK/100*dt
		min = s.thrustLevel - Data.SP.Thrust_rev_acc*massK/100*dt
	}
	if Data.SP.Thrust_rev == 0 && min < 0 {
		min = 0
	}
	s.thrustLevel = Clamp(thrustInput, min, max)

	var accel float64
	switch {
	case s.thrustLevel >= 0:
		accel = s.thrustLevel * Data.SP.Thrust / Data.BSP.Mass
	case s.thrustLevel < 0:
		accel = s.thrustLevel * Data.SP.Thrust_rev / Data.BSP.Mass
	}

	Data.PilotData.ThrustVector = v2.InDir(Data.PilotData.Ship.Ang).Mul(accel)
	Data.PilotData.Ship.Vel.DoAddMul(v2.InDir(Data.PilotData.Ship.Ang), accel*dt)
}

func (s *cosmoScene) procShipGravity(dt float64) {
	sumF := SumGravity(Data.PilotData.Ship.Pos, Data.StateData.Galaxy)
	Data.PilotData.Ship.Vel.DoAddMul(sumF, dt/Data.BSP.Mass)
}

func (s *cosmoScene) procEmissions(dt float64) {
	emissions := CalculateEmissions(Data.Galaxy, Data.PilotData.Ship.Pos)
	for emiType, emiVal := range emissions {
		switch emiType {
		case EMISSION_SLOW:
			Data.PilotData.Ship.Vel =
				Data.PilotData.Ship.Vel.Mul(1 - emiVal*dt/100)
		}
	}
}
