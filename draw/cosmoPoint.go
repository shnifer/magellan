package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
)

//TODO: add mines and fishhouses
type CosmoPoint struct {
	Sprite        *graph.CycledSprite
	EmissionRange *graph.Sprite

	ID   string
	Pos  v2.V2
	Size float64
	Type string

	lastT float64
}

func NewCosmoPoint(pd *GalaxyPoint, params graph.CamParams) *CosmoPoint {

	isMark := pd.Type == BUILDING_BLACKBOX || pd.Type == BUILDING_BEACON
	if isMark {
		params.DenyScale = true
		params.DenyAngle = true
	}

	sprite := NewAtlasSprite(pd.Type, params)
	zeroColor := color.RGBA{}
	if pd.Color != zeroColor {
		sprite.SetColor(pd.Color)
	}

	sprite.SetSize(pd.Size*2, pd.Size*2)

	if !isMark {
		sprite.SetAng(rand.Float64() * 360)
	}

	//Random spin speed
	fps := 20 * (0.5 + rand.Float64())
	cycledSprite := graph.NewCycledSprite(sprite, graph.Cycle_Loop, fps)

	emiR := 0.0
	for _, emi := range pd.Emissions {
		if emi.MainRange > emiR {
			emiR = emi.MainRange
		}
	}
	var emissionRange *graph.Sprite
	if emiR > 0 {
		emissionRange = graph.NewSprite(graph.CircleTex(), params)
		emissionRange.SetSize(emiR*2, emiR*2)
		emissionRange.SetColor(colornames.Orchid)
		emissionRange.SetAlpha(0.3)
		emissionRange.SetPos(pd.Pos)
	}

	res := CosmoPoint{
		Sprite:        cycledSprite,
		EmissionRange: emissionRange,
		Pos:           pd.Pos,
		ID:            pd.ID,
		Size:          pd.Size,
		Type:          pd.Type,
	}
	res.recalcSprite()
	return &res
}

//CosmoPoint update takes Absolute session time to calculate cosmic clocks position
func (co *CosmoPoint) Update(dt float64) {
	co.Sprite.Update(dt)
	co.recalcSprite()
}

func (co *CosmoPoint) Draw(dest *ebiten.Image) {
	if co.EmissionRange != nil {
		co.EmissionRange.Draw(dest)
	}
	co.Sprite.Draw(dest)
}

func (co *CosmoPoint) Req() (res *graph.DrawQueue) {
	res = graph.NewDrawQueue()
	if co.EmissionRange != nil {
		res.Add(co.EmissionRange, graph.Z_ABOVE_OBJECT)
	}
	res.Add(co.Sprite, graph.Z_GAME_OBJECT)
	return res
}

func (co *CosmoPoint) recalcSprite() {
	co.Sprite.SetPos(co.Pos)
	if co.EmissionRange != nil {
		co.EmissionRange.SetPos(co.Pos)
	}
}
