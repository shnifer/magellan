package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
)

const glyphSize = 32

type CosmoPoint struct {
	Sprite        *graph.CycledSprite
	SlidingSphere *SlidingSphere
	EmissionRange *graph.Sprite

	ID   string
	Pos  v2.V2
	Size float64
	Type string

	Glyphs []*graph.Sprite
	cam    *graph.Camera
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

	Glyphs := make([]*graph.Sprite, 0)
	if pd.HasMine {
		Glyphs = append(Glyphs, newGlyph(BUILDING_MINE, pd.MineOwner))
	}
	if pd.HasFishHouse {
		Glyphs = append(Glyphs, newGlyph(BUILDING_FISHHOUSE, pd.FishHouseOwner))
	}

	var slidingSphere *SlidingSphere
	if pd.ID == "earth" {
		period := 2 + rand.Float64()*2
		slidingSphere = NewAtlasSlidingSphere("terr1", params, graph.Z_GAME_OBJECT, period)
		slidingSphere.SetSize(pd.Size*2, pd.Size*2)

		if !isMark {
			slidingSphere.SetAng(rand.Float64() * 360)
		}
	}

	res := CosmoPoint{
		Sprite:        cycledSprite,
		SlidingSphere: slidingSphere,
		EmissionRange: emissionRange,
		Pos:           pd.Pos,
		ID:            pd.ID,
		Size:          pd.Size,
		Type:          pd.Type,
		Glyphs:        Glyphs,
		cam:           params.Cam,
	}
	res.recalcSprite()
	return &res
}

//CosmoPoint update takes Absolute session time to calculate cosmic clocks position
func (co *CosmoPoint) Update(dt float64) {
	co.Sprite.Update(dt)
	if co.SlidingSphere != nil {
		co.SlidingSphere.Update(dt)
	}
	co.recalcSprite()
}

/*
func (co *CosmoPoint) Draw(dest *ebiten.Image) {
	if co.EmissionRange != nil {
		co.EmissionRange.Draw(dest)
	}
	co.Sprite.Draw(dest)
}
*/

func (co *CosmoPoint) Req() (res *graph.DrawQueue) {
	res = graph.NewDrawQueue()
	if co.EmissionRange != nil {
		res.Add(co.EmissionRange, graph.Z_ABOVE_OBJECT)
	}
	if co.SlidingSphere == nil {
		res.Add(co.Sprite, graph.Z_GAME_OBJECT)
	} else {
		res.Append(co.SlidingSphere)
	}

	for i, sprite := range co.Glyphs {
		pos := co.cam.Apply(co.Pos)
		size := co.cam.Scale*co.Size/2 + glyphSize/2
		pos.DoAddMul(v2.V2{X: -1, Y: -1}, size)
		pos.DoAddMul(v2.V2{X: glyphSize, Y: 0}, float64(i))
		sprite.SetPos(pos)
		res.Add(sprite, graph.Z_ABOVE_OBJECT)
	}
	return res
}

func (co *CosmoPoint) recalcSprite() {
	co.Sprite.SetPos(co.Pos)
	if co.SlidingSphere != nil {
		co.SlidingSphere.SetPos(co.Pos)
	}
	if co.EmissionRange != nil {
		co.EmissionRange.SetPos(co.Pos)
	}
}

func newGlyph(t string, owner string) *graph.Sprite {
	res := NewAtlasSprite(t, graph.NoCam)
	res.SetSize(glyphSize, glyphSize)
	clr := ColorByOwner(owner)
	res.SetColor(clr)
	return res
}
