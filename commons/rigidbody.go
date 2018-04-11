package commons

import (
	"github.com/Shnifer/magellan/graph"
)

//just Rigid Body data
type RBData struct {
	Pos    graph.Point
	Ang    float64
	Vel    graph.Point
	AngVel float64
}

func (rb RBData) Extrapolate(dt float64) RBData {
	rb.Pos.X += rb.Vel.X * dt
	rb.Pos.Y += rb.Vel.Y * dt
	rb.Ang += rb.AngVel * dt
	return rb
}
