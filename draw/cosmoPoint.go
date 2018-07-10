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

	//warpRanges
	WarpCircles []*graph.CircleLine

	ID   string
	Pos  v2.V2
	Size float64
	Type string

	caption *graph.Text

	//hardcoded 2 rows
	glyphs glyphs

	cam *graph.Camera
}

func NewCosmoPoint(pd *GalaxyPoint, params graph.CamParams) *CosmoPoint {

	if pd.IsVirtual {
		Log(LVL_ERROR, "NewCosmoPoint called for virtual galaxyPoint")
		return nil
	}

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
	markScale := MarkScaleLevel(pd.GLevel)
	markSize := markScale * Mark_size
	markSprite.SetSize(markSize*2, markSize*2)
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
		simpleSprite = NewAtlasRoundSprite(spriteAN, params)
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

	glyphs := newGlyphs(pd)

	emiR := 0.0
	for _, emi := range pd.Emissions {
		if emi.MainRange > emiR {
			emiR = emi.MainRange
		}
	}

	var captionText *graph.Text
	res := CosmoPoint{
		level:          pd.GLevel,
		MarkGlowSprite: markGlow,
		MarkSprite:     markCS,
		markLevelScale: markScale,
		SimpleSprite:   simpleSprite,
		SlidingSphere:  slidingSphere,
		CycledSprite:   cycledSprite,
		Pos:            pd.Pos,
		ID:             pd.ID,
		Size:           pd.Size,
		Type:           pd.Type,
		glyphs:         glyphs,
		caption:        captionText,
		cam:            params.Cam,
	}
	switch pd.Type {
	case BUILDING_BEACON, BUILDING_BLACKBOX:
		res.SetCaption(pd.ScanData, colornames.Red)
	default:
	}

	if pd.Type == GPT_WARP {
		rads := make([]float64, 4)
		rads[0] = pd.WarpRedOutDist
		rads[1] = pd.WarpGreenInDist
		rads[2] = pd.WarpGreenOutDist
		rads[3] = pd.WarpYellowOutDist

		circles := make([]*graph.CircleLine, 4)
		opts := graph.CircleLineOpts{
			PCount: 64,
			Layer:  graph.Z_UNDER_OBJECT,
			Params: params,
		}
		for i := 0; i < len(circles); i++ {
			switch i {
			case 0:
				opts.Clr = colornames.Red
			case 1, 2:
				opts.Clr = colornames.Green
			case 3:
				opts.Clr = colornames.Yellow
			}
			circles[i] = graph.NewCircleLine(v2.ZV, rads[i], opts)
		}
		res.WarpCircles = circles
	}

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

}

func (co *CosmoPoint) Req(Q *graph.DrawQueue) {

	co.recalcSprite()

	markAlpha, spriteAlpha := MarkAlpha(co.Size*2/co.markLevelScale, co.cam)

	if markAlpha > 0 && co.MarkSprite != nil {
		co.MarkSprite.SetAlpha(markAlpha)
		co.MarkGlowSprite.SetAlpha(markAlpha)

		Q.Add(co.MarkSprite, graph.Z_GAME_OBJECT-co.level)
		Q.Add(co.MarkGlowSprite, graph.Z_GAME_OBJECT-10)
	}

	if spriteAlpha > 0 {
		if lowQ {
			if co.SimpleSprite != nil {
				co.SimpleSprite.SetAlpha(spriteAlpha)
				Q.Add(co.SimpleSprite, graph.Z_GAME_OBJECT)
			}
		} else {
			if co.SlidingSphere != nil {
				co.SlidingSphere.SetAlpha(spriteAlpha)
				Q.Add(co.SlidingSphere, graph.Z_GAME_OBJECT)
			}
			if co.CycledSprite != nil {
				co.CycledSprite.SetAlpha(spriteAlpha)
				Q.Add(co.CycledSprite, graph.Z_GAME_OBJECT)
			}
		}
	}

	if co.caption != nil {
		base := co.cam.Apply(co.Pos)
		off := v2.V2{X: 0, Y: 30}.Mul(graph.GS())
		co.caption.SetPosPivot(base.Add(off), graph.Center())
		Q.Add(co.caption, graph.Z_ABOVE_OBJECT)
	}

	if markAlpha == 0 && co.WarpCircles != nil {
		for _, wcl := range co.WarpCircles {
			Q.Append(wcl)
		}
	}

	co.glyphs.setPos(co.cam.Apply(co.Pos))
	co.glyphs.setSize(co.cam.Scale * co.Size)

	Q.Append(co.glyphs)
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
	if co.WarpCircles != nil {
		for _, c := range co.WarpCircles {
			c.SetPos(co.Pos)
		}
	}
}

//must be set once before creating CosmoPoints
func LowQualityCosmoPoint(v bool) {
	lowQ = v
}

func (s *CosmoPoint) SetCaption(caption string, clr color.Color) {
	s.caption = graph.NewText(caption, Fonts[Face_cap], clr)
}
