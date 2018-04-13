package commons

import (
	"github.com/Shnifer/magellan/v2"
)

//just Rigid Body data
type RBData struct {
	Pos    v2.V2
	Ang    float64
	Vel    v2.V2
	AngVel float64
}

func (rb RBData) Extrapolate(dt float64) RBData {
	rb.Pos.DoAddMul(rb.Vel, dt)
	rb.Ang += rb.AngVel * dt
	return rb
}
