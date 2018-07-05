package draw

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"golang.org/x/image/font"
	"image/color"
	"log"
)

type TextInput struct {
	sprite *graph.Sprite
	gText  *graph.Text
	face   font.Face
	clr    color.Color
	text   string
	layer  int
	//enter onDone(text, true), esc onDone(text, false)
	onDone func(string, bool)
}

func NewTextInput(sprite *graph.Sprite, face font.Face, clr color.Color, layer int, onDone func(string, bool)) *TextInput {
	return &TextInput{
		sprite: sprite,
		face:   face,
		clr:    clr,
		layer:  layer,
		onDone: onDone,
	}
}

func (ti *TextInput) Update(dt float64) {
	input := ebiten.InputChars()
	back := inpututil.IsKeyJustPressed(ebiten.KeyBackspace)
	enter := inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeyKPEnter)

	esc := inpututil.IsKeyJustPressed(ebiten.KeyEscape)
	if len(input) == 0 && ti.gText != nil && !back && !enter && !esc {
		return
	}
	ti.text += string(input)
	if back && len(ti.text) > 0 {
		runes := []rune(ti.text)
		runes = runes[0 : len(runes)-1]
		ti.text = string(runes)
	}
	r := ti.sprite.GetRect()
	h := r.Max.Y - r.Min.Y
	p := v2.V2{X: float64(r.Min.X), Y: float64(r.Min.Y)}.AddMul(v2.V2{X: 1, Y: 1}, float64(h/2))
	log.Println(r)
	log.Println(h)
	log.Println(p)
	log.Println("==")
	ti.gText = graph.NewText(ti.text, ti.face, ti.clr)
	ti.gText.SetPosPivot(p, graph.MidLeft())
	if enter {
		ti.onDone(ti.text, true)
	}
	if esc {
		ti.onDone(ti.text, false)
	}
}

func (ti *TextInput) Req() *graph.DrawQueue {
	R := graph.NewDrawQueue()
	R.Add(ti.sprite, ti.layer)
	R.Add(ti.gText, ti.layer+1)
	return R
}
