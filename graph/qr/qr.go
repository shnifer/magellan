package qr

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/magellan/graph"
	. "github.com/shnifer/magellan/log"
	"github.com/skip2/go-qrcode"
)

func NewQRSpriteHUD(text string, size int) *graph.Sprite {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		Log(LVL_PANIC, err)
	}
	image, err := ebiten.NewImageFromImage(qr.Image(size), ebiten.FilterDefault)
	if err != nil {
		Log(LVL_PANIC, err)
	}
	tex := graph.TexFromImage(image, ebiten.FilterDefault, 0, 0, 0, "~qr")
	return graph.NewSpriteHUD(tex)
}
