package main

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/input"
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

	Data.PilotData.Distortion = s.thrustLevel * Data.SP.Warp_engine.Distort_max
}

func (s *warpScene) procControlTurn(dt float64) {
	turnInput := input.GetF("turn")

	s.maneurLevel = turnInput
}
