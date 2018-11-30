package draw

import (
	. "github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/v2"
	"golang.org/x/image/colornames"
	"strconv"
)

const glyphW = 37
const glyphH = 50
const maxGlyphsInRow = 3

type glyphs struct {
	//hardcoded 2 rows
	Glyphs0   []*graph.Sprite
	Glyphs1   []*graph.Sprite
	CountText []*graph.Text

	pos  v2.V2
	size float64
}

func newGlyphs(pd *GalaxyPoint) glyphs {
	glyphs0 := make([]*graph.Sprite, 0)
	glyphs1 := make([]*graph.Sprite, 0)
	countText := make([]*graph.Text, 0)

	for _, owner := range CorpNames {
		if _, ok := pd.Mines[owner]; !ok {
			continue
		}
		ct := strconv.Itoa(len(pd.Mines[owner]))
		if ct == "1" {
			ct = ""
		}
		text := graph.NewText(ct,
			Fonts[Face_list], colornames.Darkviolet)
		countText = append(countText, text)

		if len(glyphs0) < maxGlyphsInRow {
			glyphs0 = append(glyphs0, newGlyph(BUILDING_MINE+"_"+owner, owner))
		} else {
			glyphs1 = append(glyphs0, newGlyph(BUILDING_MINE+"_"+owner, owner))
		}
	}

	//we removed fishhouses
	/*
		for owner := range pd.FishHouses {
			ng := newGlyph(BUILDING_FISHHOUSE, owner)
			if num <= maxGlyphsInRow {
				glyphs0 = append(glyphs0, ng)
			} else {
				glyphs1 = append(glyphs1, ng)
			}

		}
	*/
	return glyphs{
		Glyphs0:   glyphs0,
		Glyphs1:   glyphs1,
		CountText: countText,
	}
}

func newGlyph(t string, owner string) *graph.Sprite {
	res := NewAtlasSprite("MAGIC_GLYPH_"+t, graph.NoCam)
	res.SetSize(glyphW, glyphH)
	//clr := ColorByOwner(owner)
	//res.SetColor(clr)
	return res
}

func (g glyphs) Req(Q *graph.DrawQueue) {
	var pos v2.V2
	basePos := g.pos.AddMul(v2.V2{X: -1, Y: -1}, g.size/2).Sub(v2.V2{X: glyphW / 2, Y: glyphH / 2})
	n := 0
	for i, sprite := range g.Glyphs0 {
		pos = basePos.AddMul(v2.V2{X: glyphW, Y: 0}, float64(i))
		sprite.SetPos(pos)
		g.CountText[n].SetPosPivot(pos, graph.Center())
		Q.Add(sprite, graph.Z_ABOVE_OBJECT)
		Q.Add(g.CountText[n], graph.Z_ABOVE_OBJECT+1)
		n++
	}
	basePos.DoAddMul(v2.V2{X: 0, Y: glyphH}, 1)
	for i, sprite := range g.Glyphs1 {
		pos = basePos.AddMul(v2.V2{X: glyphW, Y: 0}, float64(i))
		sprite.SetPos(pos)
		g.CountText[n].SetPosPivot(pos, graph.Center())
		Q.Add(sprite, graph.Z_ABOVE_OBJECT)
		Q.Add(g.CountText[n], graph.Z_ABOVE_OBJECT+1)
		n++
	}
}

func (g *glyphs) setPos(pos v2.V2) {
	g.pos = pos
}

func (g *glyphs) setSize(size float64) {
	g.size = size / 2
	if g.size < Mark_size/2 {
		g.size = Mark_size / 2
	}
}
