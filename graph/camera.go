package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
)

//exec Recalc() after changes
type Camera struct {
	//center in world
	Pos v2.V2
	//center on screen
	Center v2.V2
	//angle in deg, PLUS camera counterclock (image clockwise)
	AngleDeg float64
	//scale
	Scale float64
	g, r  ebiten.GeoM
}

func NewCamera() *Camera {
	res := &Camera{
		Scale: 1,
	}
	res.Recalc()
	return res
}

func (c *Camera) Geom() ebiten.GeoM {
	return c.g
}
func (c *Camera) Apply(p v2.V2) v2.V2 {
	x, y := c.g.Apply(p.X, p.Y)
	return v2.V2{x, y}
}

func (c *Camera) UnApply(p v2.V2) v2.V2 {
	x, y := c.r.Apply(p.X, p.Y)
	return v2.V2{x, y}
}

func (c *Camera) Recalc() {
	G := ebiten.GeoM{}
	//Translate relevante to camera
	G.Translate(-c.Pos.X, -c.Pos.Y)
	//Rotate and scale
	G.Rotate(-c.AngleDeg * Deg2Rad)
	G.Scale(c.Scale, -c.Scale)
	//Translate to screen center
	G.Translate(c.Center.X, c.Center.Y)
	c.g = G
	G.Invert()
	c.r = G
}
