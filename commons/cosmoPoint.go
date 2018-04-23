package commons

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
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
	tex, col := GetAtlasTexColor(pd.Type)
	sprite := graph.NewSprite(tex, cam, false, false)
	sprite.SetColor(col)
	sprite.SetSize(pd.Size*2, pd.Size*2)
	cycledSprite := graph.NewCycledSprite(sprite, graph.Cycle_Loop, 20)
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

	co.Sprite.Update(dt)

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
