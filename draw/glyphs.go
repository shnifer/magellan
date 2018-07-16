package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"strconv"
	"image/color"
)

const glyphSize = 32
const maxGlyphsInRow = 3

type glyphs struct {
	//hardcoded 2 rows
	Glyphs0   []*graph.Sprite
	Glyphs1   []*graph.Sprite
	CountText []*graph.Text

	pos  v2.V2
	size float64
}

//todo: beacon and blackbox glyphs?
//todo: warp mines for control center
func newGlyphs(pd *GalaxyPoint) glyphs {
	glyphs0 := make([]*graph.Sprite, 0)
	glyphs1 := make([]*graph.Sprite, 0)
	countText:= make([]*graph.Text, 0)

	for _, owner :=range CorpNames{
		if _,ok:=pd.Mines[owner]; !ok{
			continue
		}
		text:=graph.NewText(strconv.Itoa(len(pd.Mines[owner])), Fonts[Face_list], color.White)
		countText = append(countText,text)
		if len(glyphs0) < maxGlyphsInRow {
			glyphs0 = append(glyphs0, newGlyph(BUILDING_MINE, owner))
		} else {
			glyphs1 = append(glyphs0, newGlyph(BUILDING_MINE, owner))
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

func (g glyphs) Req(Q *graph.DrawQueue) {
	var pos v2.V2
	basePos := g.pos.AddMul(v2.V2{X: -1, Y: -1}, g.size)
	n:=0
	for i, sprite := range g.Glyphs0 {
		pos = basePos.AddMul(v2.V2{X: glyphSize, Y: 0}, float64(i))
		sprite.SetPos(pos)
		g.CountText[n].SetPosPivot(pos, graph.Center())
		n++
		Q.Add(sprite, graph.Z_ABOVE_OBJECT)
		Q.Add(sprite, graph.Z_ABOVE_OBJECT+1)
	}
	basePos.DoAddMul(v2.V2{X: 0, Y: glyphSize}, 1)
	for i, sprite := range g.Glyphs1 {
		pos = basePos.AddMul(v2.V2{X: glyphSize, Y: 0}, float64(i))
		sprite.SetPos(pos)
		g.CountText[n].SetPosPivot(pos, graph.Center())
		n++
		Q.Add(sprite, graph.Z_ABOVE_OBJECT)
		Q.Add(sprite, graph.Z_ABOVE_OBJECT+1)
	}
}

func (g *glyphs) setPos(pos v2.V2) {
	g.pos = pos
}

func (g *glyphs) setSize(size float64) {
	g.size = size/2 + glyphSize/2
	if g.size < Mark_size/2 {
		g.size = Mark_size / 2
	}
}
