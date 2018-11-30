package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/shnifer/magellan/draw"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/v2"
)

var s *graph.SlidingSprite

func run(window *ebiten.Image) error {
	window.Clear()

	s.AddSlide(0.003)
	s.Draw(window)

	return nil
}

func main() {
	draw.InitTexAtlas()
	tex := draw.GetAtlasTex("terr1")
	tex = graph.SlidingTex(tex)
	sprite := graph.NewSprite(tex, graph.NoCam)
	sprite.SetPivot(graph.TopLeft())
	sprite.SetPivot(v2.ZV)
	s = graph.NewSlidingSprite(sprite)
	ebiten.Run(run, 200, 200, 1, "test")
}
