package graph

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/skip2/go-qrcode"
)

func NewQRSpriteHUD(text string, size int) *Sprite {
	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		panic(err)
	}
	image, err := ebiten.NewImageFromImage(qr.Image(size), ebiten.FilterDefault)
	if err != nil {
		panic(err)
	}
	tex := TexFromImage(image, ebiten.FilterDefault, 0, 0, 0, "~qr")
	return NewSpriteHUD(tex)
}
