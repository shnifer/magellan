package commons

import (
	"github.com/Shnifer/magellan/graph/flow"
	"github.com/Shnifer/magellan/v2"
	"log"
	"math"
	"strconv"
	"strings"
)

const (
	velDev      = 25
	spawnPosDev = 10
	attrFDev    = 20
)

func sigFuncStrDecoder(data string) (fn string, params []float64) {
	a := strings.Split(data, " ")
	fn = a[0]
	params = make([]float64, len(a)-1)
	for i := 0; i < len(params); i++ {
		f, err := strconv.ParseFloat(a[i+1], 64)
		if err != nil {
			f = 0
			log.Println("can't parse ", data, "value", a[i+1])
		}
		params[i] = f
	}
	return fn, params
}

func SignatureVelSpawn(sig Signature) (VelocityF func(pos v2.V2) v2.V2, SpawnPos func() (pos v2.V2)) {
	fn, params := sigFuncStrDecoder(sig.Type().VelAndSpawnStr)
	k := params[0]
	devV := sig.DevV(SIG_VELSPAWN)
	devKVel := sig.DevK(SIG_VELSPAWN, velDev) * k
	devKPos := sig.DevK(SIG_VELSPAWN, spawnPosDev) * k

	switch fn {
	case "const":
		VelocityF = func(v2.V2) v2.V2 { return v2.ZV }
		SpawnPos = func() (pos v2.V2) {
			return flow.RandomInCirc(1*devKPos)().AddMul(devV, 0.5)
		}

	case "rotation":
		VelocityF = flow.ComposeRadial(flow.ConstC(k), flow.ConstC(params[1]*devKVel)).AddMul(devV, 0.5)
		SpawnPos = flow.RandomInCirc(1)
	case "sinFloat":
		sinx := func(x, y float64) float64 {
			return math.Sin(y*params[1]) * params[2]
		}
		VelocityF = flow.ComposeDecart(sinx, flow.ConstC(devV.Len())).Rot(devV.Dir())
		SpawnPos = flow.RandomOnSide(devV.Normed().Mul(-1), 0.2*devKPos)
	case "linearFloat":
		VelocityF = func(v2.V2) v2.V2 { return devV.Mul(k) }
		SpawnPos = flow.RandomOnSide(devV.Normed().Mul(-1), 0.2*devKPos)
	default:
		log.Panicln("unknown VelAndSpawnStr", sig.Type().VelAndSpawnStr)
	}
	return VelocityF, SpawnPos
}

func SignatureAttrF(sig Signature, fstr string, koefName string) (res flow.AttrF) {
	fn, params := sigFuncStrDecoder(fstr)
	k := params[0]
	devK := sig.DevK(koefName, attrFDev) * k

	switch fn {
	case "const":
		res = func(p flow.Point) float64 {
			return devK
		}
	case "upAndDown":
		res = flow.SinMaxTime(0, devK, 0.5)
	case "sinPulse":
		res = func(p flow.Point) float64 {
			return flow.SinLifeTime(0, 1, params[1]*devK)(p) * k + params[2]
		}
	case "linear":
		res = func(p flow.Point) float64 {
			return flow.LinearLifeTime(params[1], params[2])(p) * devK
		}
	default:
		log.Panicln("unknown SignatureAttrF", fn)
	}

	return res
}
