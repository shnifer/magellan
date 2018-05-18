package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/graph/flow"
)

func CreateSignatureFlow(signature Signature, params graph.CamParams) *flow.Flow {
	sprite := NewAtlasSprite(signature.Particle().SpriteName, params)

	drawer := flow.SpriteDrawerParams{
		Sprite:       sprite,
		DoRandomLine: true,
		FPS:          10,
		CycleType:    graph.Cycle_PingPong,
		Layer:        graph.Z_GAME_OBJECT,
	}.New
	_ = drawer()
	return nil
}
