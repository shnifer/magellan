package draw

import (
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/graph"
	"github.com/Shnifer/magellan/graph/flow"
	"github.com/Shnifer/magellan/v2"
)

var signatureFlowCache map[Signature]*flow.Flow

const (
	spawnPeriodDev = 25
)

func init (){
	signatureFlowCache = make(map[Signature]*flow.Flow)
}

func GetSignatureFlow(signature Signature, camParams graph.CamParams, layer int) *flow.Flow {
	particle:=signature.Particle()
	param:=signature.Type()

	sprite := NewAtlasSprite(signature.Particle().SpriteName, camParams)
	drawer := flow.SpriteDrawerParams{
		Sprite:       sprite,
		DoRandomLine: particle.DoRandomLine,
		FPS:          particle.FPS,
		CycleType:    particle.CycleType,
		Layer:        layer,
	}.New

	AttrFs := flow.NewAttrFs()

	spawnPeriod:=param.SpawnPeriod*signature.DevK(SIG_SPAWNPERIOD,spawnPeriodDev)

	velocityF,spawnPos:=SignatureVelSpawn(signature)

	return flow.Params{
		SpawnPeriod:    spawnPeriod,
		SpawnPos:       flow.RandomInCirc(1),
		VelocityF:      velFs["rotation"],
		SpawnLife:      flow.NormRand(medLife, devLife),
		SpawnUpdDrawer: drawer,
		AttrFs:         AttrFs,
	}.New()
}
