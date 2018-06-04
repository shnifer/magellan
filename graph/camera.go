package graph

import (
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"math"
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

	//in screen pixels, zero means no clipping
	ClipW, ClipH float64
}

var NoCam CamParams

type CamParams struct {
	Cam       *Camera
	DenyScale bool
	DenyAngle bool
}

func NewCamParams(cam *Camera, denyScale, denyAngle bool) CamParams {
	return CamParams{
		Cam:       cam,
		DenyAngle: denyAngle,
		DenyScale: denyScale,
	}
}
func (c *Camera) Params(denyScale, denyAngle bool) CamParams {
	return NewCamParams(c, denyScale, denyAngle)
}
func (c *Camera) Phys() CamParams {
	return c.Params(false, false)
}
func (c *Camera) Deny() CamParams {
	return c.Params(true, true)
}
func (c *Camera) FixS() CamParams {
	return c.Params(true, false)
}

func NewCamera() *Camera {
	res := &Camera{
		Scale: 1,
		ClipW: winW,
		ClipH: winH,
	}
	res.Recalc()
	return res
}

func (c *Camera) Geom() ebiten.GeoM {
	return c.g
}
func (c *Camera) Apply(p v2.V2) v2.V2 {
	x, y := c.g.Apply(p.X, p.Y)
	return v2.V2{X: x, Y: y}
}

func (c *Camera) UnApply(p v2.V2) v2.V2 {
	x, y := c.r.Apply(p.X, p.Y)
	return v2.V2{X: x, Y: y}
}

func (c *Camera) Recalc() {
	G := ebiten.GeoM{}
	//Translate relevant to camera
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

//if object is in camera space always return true
//may return true if object is out of space
func (c *Camera) PointInSpace(p v2.V2) bool{
	if c.ClipW*c.ClipH==0 {
		return true
	}
	delta:=c.Apply(p).Sub(c.Center)
	maxX:=c.ClipW/2
	maxY:=c.ClipH/2
	if delta.X>maxX|| delta.X<(-maxX) ||
	delta.Y>maxY || delta.Y<(-maxY) {
		return false
	}
	return true
}
func (c *Camera) CircleInSpace(center v2.V2, radius float64) bool{
	if c.ClipW*c.ClipH==0 {
		return true
	}
	delta:=c.Apply(center).Sub(c.Center)
	scrRadius:=radius*c.Scale
	maxX:=c.ClipW/2+scrRadius
	maxY:=c.ClipH/2+scrRadius

	if delta.X>maxX|| delta.X<(-maxX) ||
		delta.Y>maxY || delta.Y<(-maxY) {
		return false
	}
	return true
}

//do not count angle, just use outer circle
func (c *Camera) RectInSpace(center v2.V2, w,h float64) bool{
	if c.ClipW*c.ClipH==0 {
		return true
	}
	radius:=math.Sqrt(w*w+h*h)/2
	return c.CircleInSpace(center, radius)
}

func (c *Camera) SetClip(w,h int) {
	c.ClipW=float64(w)
	c.ClipH=float64(h)
}