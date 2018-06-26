package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	. "github.com/Shnifer/magellan/log"
	"github.com/Shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
)

const glyphSize = 32

var lowQ bool

type CosmoPoint struct {
	level int
	//mark for far scale
	markLevelScale float64
	MarkGlowSprite *graph.Sprite
	MarkSprite     *graph.CycledSprite

	//static sprite for lowQ
	SimpleSprite *graph.Sprite

	//for not-so-lowQ
	//we use sliding or cycling depends on Type
	SlidingSphere *SlidingSphere
	CycledSprite  *graph.CycledSprite

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

	markNameAN := "MAGIC_MARK_" + pd.Type
	if pd.ID == "earth" {
		markNameAN = MARKtheEarthAN
	} else if pd.ID == "magellan" {
		markNameAN = MARKtheMagellanAN
	}

	markSprite = NewAtlasSprite(markNameAN, markParam)
	markScale := MarkScaleLevel(pd.Level)
	markSprite.SetScale(markScale, markScale)
	markCS := graph.NewCycledSprite(markSprite, graph.Cycle_Loop, 10)

	var typeGlowScale float64
	switch pd.Type {
	case GPT_ASTEROID:
		typeGlowScale = 0.1
	case GPT_STAR:
		typeGlowScale = 2
	default:
		typeGlowScale = 1.5
	}
	if pd.ID == "magellan" {
		typeGlowScale = 2.5
	}
	var markGlow *graph.Sprite
	markGlow = NewAtlasSprite(MARKGLOWAN, markParam)
	markGlow.SetScale(markScale*typeGlowScale, markScale*typeGlowScale)
	markGlow.SetColor(colornames.Deepskyblue)

	spriteAN := pd.SpriteAN
	if spriteAN == "" {
		spriteAN = "MAGIC_DEFAULT_" + pd.Type
	}

	var simpleSprite *graph.Sprite
	var slidingSphere *SlidingSphere
	var cycledSprite *graph.CycledSprite

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
		switch pd.Type {
		case GPT_STAR, GPT_GASPLANET, GPT_HARDPLANET:
			//sliding sphere
			period := 2 + rand.Float64()*2
			slidingSphere = NewAtlasSlidingSphere(spriteAN, params, period)
			if pd.Color != zeroColor {
				slidingSphere.SetColor(pd.Color)
			}
			slidingSphere.SetSize(pd.Size*2, pd.Size*2)
			slidingSphere.SetAng(ang)
		case GPT_WARP, GPT_ASTEROID, BUILDING_BEACON, BUILDING_BLACKBOX:
			//cycling sprite
			simpleSprite := NewAtlasSprite(spriteAN, params)
			simpleSprite.SetSize(pd.Size*2, pd.Size*2)
			simpleSprite.SetAng(ang)
			cycledSprite = graph.NewCycledSprite(simpleSprite, graph.Cycle_Loop, 20)
		default:
			Log(LVL_ERROR, "Unknown galaxy point type ", pd.Type)
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
		level:          pd.Level,
		MarkGlowSprite: markGlow,
		MarkSprite:     markCS,
		markLevelScale: markScale,
		SimpleSprite:   simpleSprite,
		SlidingSphere:  slidingSphere,
		CycledSprite:   cycledSprite,
		EmissionRange:  emissionRange,
		Pos:            pd.Pos,
		ID:             pd.ID,
		Size:           pd.Size,
		Type:           pd.Type,
		Glyphs:         Glyphs,
		cam:            params.Cam,
	}
	res.recalcSprite()
	return &res
}

func (co *CosmoPoint) Update(dt float64) {
	if co.MarkSprite != nil {
		co.MarkSprite.Update(dt)
	}
	if co.SlidingSphere != nil {
		co.SlidingSphere.Update(dt)
	}
	if co.CycledSprite != nil {
		co.CycledSprite.Update(dt)
	}

	co.recalcSprite()
}

func (co *CosmoPoint) Req() (res *graph.DrawQueue) {
	res = graph.NewDrawQueue()
	if co.EmissionRange != nil {
		res.Add(co.EmissionRange, graph.Z_UNDER_OBJECT)
	}

	markAlpha, spriteAlpha := markAlpha(co.Size*2/co.markLevelScale, co.cam)

	if markAlpha > 0 && co.MarkSprite != nil {
		co.MarkSprite.SetAlpha(markAlpha)
		co.MarkGlowSprite.SetAlpha(markAlpha)

		res.Add(co.MarkSprite, graph.Z_GAME_OBJECT-co.level)
		res.Add(co.MarkGlowSprite, graph.Z_GAME_OBJECT-10)
	}

	if spriteAlpha > 0 {
		if lowQ {
			if co.SimpleSprite != nil {
				co.SimpleSprite.SetAlpha(spriteAlpha)
				res.Add(co.SimpleSprite, graph.Z_GAME_OBJECT)
			}
		} else {
			if co.SlidingSphere != nil {
				co.SlidingSphere.SetAlpha(spriteAlpha)
				res.Add(co.SlidingSphere, graph.Z_GAME_OBJECT)
			}
			if co.CycledSprite != nil {
				co.CycledSprite.SetAlpha(spriteAlpha)
				res.Add(co.CycledSprite, graph.Z_GAME_OBJECT)
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
	if co.CycledSprite != nil {
		co.CycledSprite.SetPos(co.Pos)
	}
	if co.MarkSprite != nil {
		co.MarkSprite.SetPos(co.Pos)
	}
	if co.MarkGlowSprite != nil {
		co.MarkGlowSprite.SetPos(co.Pos)
	}
	if co.EmissionRange != nil {
		co.EmissionRange.SetPos(co.Pos)
	}
}

func newGlyph(t string, owner string) *graph.Sprite {
	res := NewAtlasSprite("MAGIC_GLYPH_"+t, graph.NoCam)
	res.SetSize(glyphSize, glyphSize)
	clr := ColorByOwner(owner)
	res.SetColor(clr)
	return res
}

//must be set once before creating CosmoPoints
func LowQualityCosmoPoint(v bool) {
	lowQ = v
}
