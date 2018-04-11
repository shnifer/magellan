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

func (rb RBData) Extrapolate(dt float64) (res RBData) {
	res = rb
	res.Pos.X += rb.Vel.X * dt
	res.Pos.Y += rb.Vel.Y * dt
	res.Ang += rb.AngVel * dt
	return res
}
