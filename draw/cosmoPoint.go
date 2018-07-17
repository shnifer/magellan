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
	WarpInner *graph.Sprite
	WarpOuter *graph.Sprite
	WarpGreen *graph.WavedCircle
	WarpSize  float64

	ID   string
	Pos  v2.V2
	Size float64
	Type string

	caption     *graph.Text
	captionText string

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
		typeGlowScale = 0.3
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
		case GPT_WARP, GPT_ASTEROID, GPT_WORMHOLE, BUILDING_BEACON, BUILDING_BLACKBOX:
			//cycling sprite
			simpleSprite := NewAtlasSprite(spriteAN, params)
			simpleSprite.SetSize(pd.Size*2, pd.Size*2)
			simpleSprite.SetAng(ang)
			if pd.Color != zeroColor {
				simpleSprite.SetColor(pd.Color)
			}
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
		res.WarpSize = pd.WarpRedOutDist

		res.WarpInner = NewAtlasSprite(WarpInnerAN, params)
		res.WarpInner.SetSize(pd.WarpRedOutDist*2, pd.WarpRedOutDist*2)
		res.WarpInner.SetPos(pd.Pos)
		res.WarpInner.SetColor(pd.InnerColor)
		res.WarpOuter = NewAtlasSprite(WarpOuterAN, params)
		res.WarpOuter.SetSize(pd.WarpYellowOutDist*2, pd.WarpYellowOutDist*2)
		res.WarpOuter.SetPos(pd.Pos)
		res.WarpOuter.SetColor(pd.OuterColor)

		sprite := NewAtlasSprite(WarpGreenAN, params)
		sprite.SetColor(pd.GreenColor)

		opts := graph.WavedCircleOpts{
			Sprite:  sprite,
			PCount:  32,
			Layer:   graph.Z_UNDER_OBJECT + 50,
			Params:  params,
			RandGen: pd.ID,
		}
		res.WarpGreen = graph.NewWavedCircle(
			pd.Pos, pd.WarpGreenInDist, pd.WarpGreenOutDist, opts)
	}

	return &res
}

func (cp *CosmoPoint) Update(dt float64) {
	if cp.MarkSprite != nil {
		cp.MarkSprite.Update(dt)
	}
	if cp.SlidingSphere != nil {
		cp.SlidingSphere.Update(dt)
	}
	if cp.CycledSprite != nil {
		cp.CycledSprite.Update(dt)
	}
	if cp.WarpGreen != nil {
		cp.WarpGreen.Update(dt)
	}
}

func (cp *CosmoPoint) Req(Q *graph.DrawQueue) {

	clipped := false
	if cp.cam != nil {
		size := cp.Size
		if cp.WarpSize > size {
			size = cp.WarpSize
		}
		clipped = !cp.cam.CircleInSpace(cp.Pos, size)
	}
	if clipped {
		return
	}

	cp.recalcSprite()

	markAlpha, spriteAlpha := MarkAlpha(cp.Size*2/cp.markLevelScale, cp.cam)

	if markAlpha > 0 && cp.MarkSprite != nil {
		cp.MarkSprite.SetAlpha(markAlpha)
		cp.MarkGlowSprite.SetAlpha(markAlpha)

		Q.Add(cp.MarkSprite, graph.Z_GAME_OBJECT-cp.level)
		Q.Add(cp.MarkGlowSprite, graph.Z_GAME_OBJECT-10)
	}

	if spriteAlpha > 0 {
		if lowQ {
			if cp.SimpleSprite != nil {
				cp.SimpleSprite.SetAlpha(spriteAlpha)
				Q.Add(cp.SimpleSprite, graph.Z_GAME_OBJECT)
			}
		} else {
			if cp.SlidingSphere != nil {
				cp.SlidingSphere.SetAlpha(spriteAlpha)
				Q.Add(cp.SlidingSphere, graph.Z_GAME_OBJECT)
			}
			if cp.CycledSprite != nil {
				cp.CycledSprite.SetAlpha(spriteAlpha)
				Q.Add(cp.CycledSprite, graph.Z_GAME_OBJECT)
			}
		}
	}

	if cp.caption != nil {
		base := cp.cam.Apply(cp.Pos)
		off := v2.V2{X: 0, Y: 30}.Mul(graph.GS())
		cp.caption.SetPosPivot(base.Add(off), graph.Center())
		Q.Add(cp.caption, graph.Z_ABOVE_OBJECT)
	}

	if markAlpha == 0 {
		if cp.WarpOuter != nil {
			Q.Add(cp.WarpOuter, graph.Z_UNDER_OBJECT)
		}
		if cp.WarpInner != nil {
			Q.Add(cp.WarpInner, graph.Z_UNDER_OBJECT+1)
		}
		if cp.WarpGreen != nil {
			Q.Append(cp.WarpGreen)
		}
	}

	cp.glyphs.setPos(cp.cam.Apply(cp.Pos))
	cp.glyphs.setSize(cp.cam.Scale * cp.Size)

	Q.Append(cp.glyphs)
}

func (cp *CosmoPoint) recalcSprite() {
	if cp.SimpleSprite != nil {
		cp.SimpleSprite.SetPos(cp.Pos)
	}
	if cp.SlidingSphere != nil {
		cp.SlidingSphere.SetPos(cp.Pos)
	}
	if cp.CycledSprite != nil {
		cp.CycledSprite.SetPos(cp.Pos)
	}
	if cp.MarkSprite != nil {
		cp.MarkSprite.SetPos(cp.Pos)
	}
	if cp.MarkGlowSprite != nil {
		cp.MarkGlowSprite.SetPos(cp.Pos)
	}
	if cp.WarpGreen != nil {
		cp.WarpGreen.SetPos(cp.Pos)
	}
	if cp.WarpOuter != nil {
		cp.WarpOuter.SetPos(cp.Pos)
	}
	if cp.WarpInner != nil {
		cp.WarpInner.SetPos(cp.Pos)
	}
}

//must be set once before creating CosmoPoints
func LowQualityCosmoPoint(v bool) {
	lowQ = v
}

func (cp *CosmoPoint) SetCaption(caption string, clr color.Color) {
	cp.caption = graph.NewText(caption, Fonts[Face_cap], clr)
	cp.captionText = caption
}

func (cp *CosmoPoint) GetCaption() string {
	return cp.captionText
}
