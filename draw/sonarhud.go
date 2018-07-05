package draw

import (
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/v2"
	"github.com/hajimehoshi/ebiten"
)

type SonarHUD struct {
	*SignaturePack

	inQ *graph.DrawQueue

	params graph.CamParams
	pos    v2.V2
	size   float64

	cut   *graph.Sprite
	img   *ebiten.Image
	layer int
}

func NewSonarHUD(pos v2.V2, size float64, params graph.CamParams, layer int) *SonarHUD {
	tempImgSize := size
	img, _ := ebiten.NewImage(int(tempImgSize), int(tempImgSize), ebiten.FilterDefault)
	innerCam := &graph.Camera{
		Scale:           size / 2,
		ClipW:           0,
		ClipH:           0,
		Center:          v2.V2{X: tempImgSize / 2, Y: tempImgSize / 2},
		DenyGlobalScale: true,
	}
	innerCam.Recalc()
	sigPack := NewSignaturePack(innerCam.Deny(), 0)

	cut := graph.NewSprite(graph.CircleTex(), graph.NoCam)
	cut.SetSize(tempImgSize, tempImgSize)
	cut.SetPivot(graph.TopLeft())
	return &SonarHUD{
		params:        params,
		pos:           pos,
		size:          size,
		SignaturePack: sigPack,
		cut:           cut,
		img:           img,
		layer:         layer,
		inQ:           graph.NewDrawQueue(),
	}
}

func (s *SonarHUD) Req(Q *graph.DrawQueue) {
	s.img.Clear()
	s.inQ.Clear()
	s.inQ.Append(s.SignaturePack)
	s.inQ.Run(s.img)
	im, op := s.cut.ImageOp()
	op.CompositeMode = ebiten.CompositeModeDestinationIn
	s.img.DrawImage(im, op)
	tex := graph.TexFromImage(s.img, ebiten.FilterDefault, 0, 0, 1, "sonarHUD")
	sprite := graph.NewSprite(tex, s.params)
	sprite.SetPos(s.pos)
	sprite.SetSize(s.size, s.size)
	Q.Add(sprite, s.layer)
}
