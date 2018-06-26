package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
)

const glyphSize = 32
const maxGlyphsInRow = 5

type glyphs struct {
	//hardcoded 2 rows
	Glyphs0 []*graph.Sprite
	Glyphs1 []*graph.Sprite

	pos  v2.V2
	size float64
}

func newGlyphs(pd *GalaxyPoint) glyphs {
	glyphs0 := make([]*graph.Sprite, 0)
	glyphs1 := make([]*graph.Sprite, 0)

	num := len(pd.Mines) + len(pd.FishHouses)
	for owner := range pd.Mines {
		glyphs0 = append(glyphs0, newGlyph(BUILDING_MINE, owner))
	}
	for owner := range pd.FishHouses {
		ng := newGlyph(BUILDING_FISHHOUSE, owner)
		if num <= maxGlyphsInRow {
			glyphs0 = append(glyphs0, ng)
		} else {
			glyphs1 = append(glyphs1, ng)
		}

	}
	return glyphs{
		Glyphs0: glyphs0,
		Glyphs1: glyphs1,
	}
}

func newGlyph(t string, owner string) *graph.Sprite {
	res := NewAtlasSprite("MAGIC_GLYPH_"+t, graph.NoCam)
	res.SetSize(glyphSize, glyphSize)
	clr := ColorByOwner(owner)
	res.SetColor(clr)
	return res
}

func (g glyphs) Req() *graph.DrawQueue {
	res := graph.NewDrawQueue()
	var pos v2.V2
	basePos := g.pos.AddMul(v2.V2{X: -1, Y: -1}, g.size)
	for i, sprite := range g.Glyphs0 {
		pos = basePos.AddMul(v2.V2{X: glyphSize, Y: 0}, float64(i))
		sprite.SetPos(pos)
		res.Add(sprite, graph.Z_ABOVE_OBJECT)
	}
	basePos.DoAddMul(v2.V2{X: 0, Y: glyphSize}, 1)
	for i, sprite := range g.Glyphs1 {
		pos = basePos.AddMul(v2.V2{X: glyphSize, Y: 0}, float64(i))
		sprite.SetPos(pos)
		res.Add(sprite, graph.Z_ABOVE_OBJECT)
	}

	return res
}

func (g *glyphs) setPos(pos v2.V2) {
	g.pos = pos
}

func (g *glyphs) setSize(size float64) {
	g.size = size/2 + glyphSize/2
	if g.size < mark_size/2 {
		g.size = mark_size / 2
	}
}
