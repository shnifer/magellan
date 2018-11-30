package draw

import (
	. "github.com/shnifer/magellan/commons"
	"github.com/shnifer/magellan/graph"
	"github.com/shnifer/magellan/graph/flow"
)

const (
	spawnPeriodDev = 25
	lifeTimeDev    = 50
)

type SignaturePack struct {
	camParams  graph.CamParams
	baseLayer  int
	flows      map[Signature]*flow.Flow
	deltaLayer map[Signature]int
	deltas     []bool
}

func NewSignaturePack(camParams graph.CamParams, baseLayer int) *SignaturePack {
	return &SignaturePack{
		camParams:  camParams,
		baseLayer:  baseLayer,
		flows:      make(map[Signature]*flow.Flow),
		deltaLayer: make(map[Signature]int),
		deltas:     make([]bool, 0),
	}
}

func (sp *SignaturePack) ActiveSignatures(sigs []Signature) {
	for _, sig := range sigs {
		if _, ok := sp.flows[sig]; !ok {
			//found delta
			delta := -1
			for i, occupied := range sp.deltas {
				if !occupied {
					delta = i
					break
				}
			}
			if delta == -1 {
				delta = len(sp.deltas)
				sp.deltas = append(sp.deltas, false)
			}
			sp.deltas[delta] = true
			sp.deltaLayer[sig] = delta
			sp.flows[sig] = createSignatureFlow(sig, sp.camParams, sp.baseLayer+delta)
		}
	}

	for sig, f := range sp.flows {
		found := false
		for _, val := range sigs {
			if val == sig {
				found = true
				break
			}
		}
		f.SetActive(found)
	}
}

func (sp *SignaturePack) Update(dt float64) {
	for sig, f := range sp.flows {
		f.Update(dt)
		if f.IsEmpty() {
			sp.deltas[sp.deltaLayer[sig]] = false
			delete(sp.deltaLayer, sig)
			delete(sp.flows, sig)
		}
	}
}

func (sp *SignaturePack) Req(Q *graph.DrawQueue) {
	for _, f := range sp.flows {
		Q.Append(f)
	}
}

func createSignatureFlow(signature Signature, camParams graph.CamParams, layer int) *flow.Flow {
	particle := signature.Particle()
	param := signature.Type()

	sprite := NewAtlasSprite(signature.Particle().SpriteName, camParams)
	sprite.SetColor(param.Color)
	drawer := flow.SpriteDrawerParams{
		Sprite:       sprite,
		DoRandomLine: particle.DoRandomLine,
		FPS:          particle.FPS,
		CycleType:    particle.CycleType,
		Layer:        layer,
	}.New

	AttrFs := flow.NewAttrFs()

	spawnPeriod := param.SpawnPeriod * signature.DevK(SIG_SPAWNPERIOD, spawnPeriodDev)
	spawnLife := flow.NormRand(param.LifeTime, signature.DevK(SIG_LIFETIME, lifeTimeDev))

	velocityF, spawnPos := SignatureVelSpawn(signature)

	AttrFs["Ang"] = SignatureAttrF(signature, param.AngStr, SIG_ANGF)
	AttrFs["Size"] = SignatureAttrF(signature, param.SizeStr, SIG_SIZEF)
	AttrFs["Alpha"] = SignatureAttrF(signature, param.AlphaStr, SIG_ALPHAF)

	return flow.Params{
		SpawnPeriod:    spawnPeriod,
		SpawnPos:       spawnPos,
		VelocityF:      velocityF,
		SpawnLife:      spawnLife,
		SpawnUpdDrawer: drawer,
		AttrFs:         AttrFs,
	}.New()
}
