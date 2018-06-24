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

var lowQ bool

const (
	mark_size   = 40
	sprite_size = 50
)

type CosmoPoint struct {
	MarkSprite    *graph.Sprite
	SimpleSprite  *graph.Sprite
	SlidingSphere *SlidingSphere

	onlyMark bool

	EmissionRange *graph.Sprite

	ID   string
	Pos  v2.V2
	Size float64
	Type string

	Glyphs []*graph.Sprite
	cam    *graph.Camera
}

func NewCosmoPoint(pd *GalaxyPoint, params graph.CamParams) *CosmoPoint {

	markParam := params
	markParam.DenyScale = true
	markParam.DenyAngle = true

	var markSprite *graph.Sprite

	onlyMark := pd.Type == BUILDING_BLACKBOX || pd.Type == BUILDING_BEACON
	if onlyMark {
		markSprite = NewAtlasSprite("MAGIC_"+pd.Type, markParam)
	} else {
		markSprite = NewAtlasSprite("MAGIC_MARK_"+pd.Type, markParam)
	}
	spriteAN := pd.SpriteAN
	if spriteAN == "" && !onlyMark {
		spriteAN = "MAGIC_DEFAULT_" + pd.Type
	}

	var simpleSprite *graph.Sprite
	var slidingSphere *SlidingSphere

	if !onlyMark {
		zeroColor := color.RGBA{}
		ang := rand.Float64() * 360
		if lowQ {
			simpleSprite = NewAtlasSprite(spriteAN, params)
			if pd.Color != zeroColor {
				simpleSprite.SetColor(pd.Color)
			}
			simpleSprite.SetSize(pd.Size*2, pd.Size*2)
			simpleSprite.SetAng(ang)
		} else {
			period := 2 + rand.Float64()*2

			slidingSphere = NewAtlasSlidingSphere(spriteAN, params, period)
			if pd.Color != zeroColor {
				slidingSphere.SetColor(pd.Color)
			}
			slidingSphere.SetSize(pd.Size*2, pd.Size*2)
			slidingSphere.SetAng(ang)
		}
	}

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

	res := CosmoPoint{
		MarkSprite:    markSprite,
		SimpleSprite:  simpleSprite,
		SlidingSphere: slidingSphere,
		EmissionRange: emissionRange,
		onlyMark:      onlyMark,
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
	if co.SlidingSphere != nil {
		co.SlidingSphere.Update(dt)
	}
	co.recalcSprite()
}

func (co *CosmoPoint) Req() (res *graph.DrawQueue) {
	res = graph.NewDrawQueue()
	if co.EmissionRange != nil {
		res.Add(co.EmissionRange, graph.Z_UNDER_OBJECT)
	}

	if co.MarkSprite != nil {
		res.Add(co.MarkSprite, graph.Z_GAME_OBJECT)
	}

	if !co.onlyMark {
		size := co.Size
		if co.cam != nil {
			size *= co.cam.Scale
		}
		var markAlpha float64
		var spriteAlpha float64
		if size <= mark_size {
			markAlpha = 1
			spriteAlpha = 0
		} else if size >= sprite_size {
			markAlpha = 0
			spriteAlpha = 1
		} else {
			k := (size - mark_size) / (sprite_size - mark_size)
			markAlpha = 1 - k
			spriteAlpha = k
		}

		if lowQ {
			co.MarkSprite.SetAlpha(markAlpha)
			co.SimpleSprite.SetAlpha(spriteAlpha)
			if spriteAlpha > 0 {
				res.Add(co.SimpleSprite, graph.Z_GAME_OBJECT)
			}
		} else {
			co.SlidingSphere.SetAlpha(spriteAlpha)
			if spriteAlpha > 0 {
				res.Add(co.SlidingSphere, graph.Z_GAME_OBJECT)
			}
		}
	}

	for i, sprite := range co.Glyphs {
		pos := co.cam.Apply(co.Pos)
		size := co.cam.Scale*co.Size/2 + glyphSize/2
		if size < mark_size/2 {
			size = mark_size / 2
		}
		pos.DoAddMul(v2.V2{X: -1, Y: -1}, size)
		pos.DoAddMul(v2.V2{X: glyphSize, Y: 0}, float64(i))
		sprite.SetPos(pos)
		res.Add(sprite, graph.Z_ABOVE_OBJECT)
	}
	return res
}

func (co *CosmoPoint) recalcSprite() {
	if co.SimpleSprite != nil {
		co.SimpleSprite.SetPos(co.Pos)
	}
	if co.SlidingSphere != nil {
		co.SlidingSphere.SetPos(co.Pos)
	}
	if co.MarkSprite != nil {
		co.MarkSprite.SetPos(co.Pos)
	}
	if co.EmissionRange != nil {
		co.EmissionRange.SetPos(co.Pos)
	}
}

func newGlyph(t string, owner string) *graph.Sprite {
	res := NewAtlasSprite("MAGIC_"+t, graph.NoCam)
	res.SetSize(glyphSize, glyphSize)
	clr := ColorByOwner(owner)
	res.SetColor(clr)
	return res
}

//must be set once before creating CosmoPoints
func LowQualityCosmoPoint(v bool) {
	lowQ = v
}
