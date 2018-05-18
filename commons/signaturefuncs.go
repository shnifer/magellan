package commons

import (
	"github.com/Shnifer/magellan/graph/flow"
	"github.com/Shnifer/magellan/v2"
	"log"
	"math"
)

const (
	velDev      = 25
	spawnPosDev = 10
	attrFDev    = 20
)

func SignatureVelSpawn(sig Signature) (VelocityF func(pos v2.V2) v2.V2, SpawnPos func() (pos v2.V2)) {
	k := sig.Coef(SIG_VELSPAWN)
	devV := sig.DevV(SIG_VELSPAWN)
	devKVel := sig.DevK(SIG_VELSPAWN, velDev) * k
	devKPos := sig.DevK(SIG_VELSPAWN, spawnPosDev) * k

	switch sig.Type().VelAndSpawnF {
	case "const":
		VelocityF = func(v2.V2) v2.V2 { return v2.ZV }
		SpawnPos = func() (pos v2.V2) {
			return flow.RandomInCirc(1*devKPos)().AddMul(devV, 0.5)
		}

	case "rotationOut":
		VelocityF = flow.ComposeRadial(flow.ConstC(k), flow.ConstC(0.1*devKVel)).AddMul(devV, 0.5)
		SpawnPos = flow.RandomInCirc(1)
	case "sinFloat":
		sinx := func(x, y float64) float64 {
			return math.Sin(y*5/k) * k / 5
		}
		VelocityF = flow.ComposeDecart(sinx, flow.ConstC(devV.Len())).Rot(devV.Dir())
		SpawnPos = flow.RandomOnSide(devV.Normed().Mul(-1), 0.2*devKPos)
	case "linearFloat":
		VelocityF = func(v2.V2) v2.V2 { return devV.Mul(k) }
		SpawnPos = flow.RandomOnSide(devV.Normed().Mul(-1), 0.2*devKPos)
	default:
		log.Panicln("unknown VelAndSpawnF", sig.Type().VelAndSpawnF)
	}
	return VelocityF, SpawnPos
}

func SignatureAttrF(sig Signature, fname string, koefName string) (res flow.AttrF) {
	k := sig.Coef(koefName)
	devK := sig.DevK(koefName, attrFDev) * k

	switch fname {
	case "const":
		res = func(p flow.Point) float64 {
			return devK
		}
	case "upAndDown":
		res = flow.SinMaxTime(0, devK, 0.5)
	case "sinPulse":
		res = func(p flow.Point) float64 {
			return flow.SinLifeTime(0, 1, devK)(p) * k
		}
	case "linearUp":
		res = func(p flow.Point) float64 {
			return flow.LinearLifeTime(0, 1)(p) * devK
		}
	case "linearDown":
		res = func(p flow.Point) float64 {
			return flow.LinearLifeTime(1, 0)(p) * devK
		}
	default:
		log.Panicln("unknown SignatureAttrF", fname)
	}

	return res
}
