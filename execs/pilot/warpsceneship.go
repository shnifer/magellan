package main

import (
	"github.com/Shnifer/magellan/input"
	"github.com/Shnifer/magellan/commons"
)

func (s *warpScene) updateShipControl(dt float64) {
	s.procControlForward(dt)
	s.procControlTurn(dt)
}

func (s *warpScene) procControlForward(dt float64) {
	w := input.WarpLevel("warpspeed")

	s.thrustLevel = commons.Clamp(w,
		s.thrustLevel-Data.SP.Warp_engine.Distort_slow/100*dt,
		s.thrustLevel+Data.SP.Warp_engine.Distort_acc/100*dt)

	if s.thrustLevel == 0{
		Data.PilotData.Distortion = 0
	} else {
		Data.PilotData.Distortion = DEFVAL.MinDistortion +
			s.thrustLevel*Data.SP.Warp_engine.Distort_max
	}
}

func (s *warpScene) procControlTurn(dt float64) {
	turnInput := input.GetF("turn")
	s.maneurLevel = turnInput
	Data.PilotData.Ship.AngVel = turnInput * Data.SP.Warp_engine.Turn_speed
}
