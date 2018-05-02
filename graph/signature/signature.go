package signature

import (
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/v2"
)

//Signature draws an animated pseudo-random effects with many parameters
type Signature struct {
}

type point struct {
	pos V2

	params map[string]float64
}

type Generator struct {
}

type Params struct {
	Tex       graph.Tex
	LinearTex bool

	SpawnRate float64

	SpawnF    SpawnF
	VelocityF VelocityF
	SizeF     SizeF
	RotF      RotF
	AlphaF    AlphaF
}

//generates distributed spawn point
//Drawing field is normed to [-1..1,-1..1]
type SpawnF func() (pos V2)

//V(x) for streaming like motion
type VelocityF func(pos V2) V2

type SizeF func(t float64) float64
type RotF func(t float64) float64
type AlphaF func(t float64) float64
