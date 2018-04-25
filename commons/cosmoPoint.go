package commons

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"image/color"
	"math/rand"
)

type CosmoPoint struct {
	Sprite *graph.CycledSprite

	ID string

	Pos  v2.V2
	Size float64

	Parent   *CosmoPoint
	Orbit    float64
	AngVel   float64
	AngPhase float64

	lastT float64

	Mass float64

	ScanData string
}

func NewCosmoPoint(pd GalaxyPoint, cam *graph.Camera) *CosmoPoint {
	sprite := NewAtlasSprite(pd.Type, cam, false, false)
	zeroColor := color.RGBA{}
	if pd.Color != zeroColor {
		sprite.SetColor(pd.Color)
	}

	sprite.SetSize(pd.Size*2, pd.Size*2)

	//Random spin speed
	fps := 20 * (0.5 + rand.Float64())
	cycledSprite := graph.NewCycledSprite(sprite, graph.Cycle_Loop, fps)
	res := CosmoPoint{
		Sprite:   cycledSprite,
		ID:       pd.ID,
		Pos:      pd.Pos,
		Size:     pd.Size,
		Orbit:    pd.Orbit,
		AngVel:   360 / pd.Period,
		Mass:     pd.Mass,
		ScanData: pd.ScanData,
	}
	res.recalcSprite()
	return &res
}

//CosmoPoint update takes Absolute session time to calculate cosmic clocks position
func (co *CosmoPoint) Update(sessionTime float64) {
	if co.lastT == 0 {
		co.lastT = sessionTime
	}
	dt := sessionTime - co.lastT
	co.lastT = sessionTime

	if dt>0 {
		co.Sprite.Update(dt)
	}

	if co.Parent != nil {
		angle := co.AngPhase + co.AngVel*sessionTime
		co.Pos = co.Parent.Pos.AddMul(v2.InDir(angle), co.Orbit)
		co.recalcSprite()
	}
}

func (co *CosmoPoint) Draw(dest *ebiten.Image) {
	co.Sprite.Draw(dest)
}

func (co *CosmoPoint) recalcSprite() {
	co.Sprite.SetPosAng(co.Pos, 0)
}

//for server calculation
func CalculateCosmoPos(name string, points []GalaxyPoint, sessionTime float64) v2.V2 {
	ind := make(map[string]int, len(points))
	for i, v := range points {
		ind[v.ID] = i
	}
	var f func(id string) v2.V2
	f = func(id string) v2.V2 {
		i, ok := ind[id]
		if !ok {
			return v2.ZV
		}
		point := points[i]
		parent := points[i].ParentID
		if parent == "" {
			return point.Pos
		}

		parentPos := f(parent)
		angle := (360 / point.Period) * sessionTime
		return parentPos.AddMul(v2.InDir(angle), point.Orbit)
	}
	return f(name)
}
