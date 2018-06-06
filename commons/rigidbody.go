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

type RBFollower struct{
	rb RBData
	elastic float64

	inited bool
	delta v2.V2
	elasticT float64
}

func NewRBFollower (elastic float64) *RBFollower{
	return &RBFollower{
		elastic: elastic,
	}
}

func (f *RBFollower) JumpTo(rb RBData) {
	f.delta = v2.ZV
	f.elasticT = 0
	f.inited = true
	f.rb = rb
}

func (f *RBFollower) MoveTo(rb RBData) {
	if f.inited{
		f.delta = f.RB().Pos.Sub(rb.Pos)
		f.elasticT = 0
	} else {
		f.inited = true
	}
	f.rb = rb
}

func (f *RBFollower) Update(dt float64) {
	if !f.inited{
		return
	}
	f.rb = f.rb.Extrapolate(dt)
	f.elasticT += dt
}

func (f *RBFollower) RB() RBData{
	if !f.inited || f.elastic==0{
		return f.rb
	}

	res := f.rb
	if f.elasticT<f.elastic{
		k:=(f.elastic-f.elasticT) / f.elastic
		res.Pos.DoAddMul(f.delta, k)
	}

	return res
}