package commons

import (
	"github.com/shnifer/magellan/v2"
	"math"
)

var gravityConst float64
var warpGravityConst float64
var warpGravThreshold float64

var VelDistWarpK float64

func SetGravityConsts(G, W float64) {
	gravityConst = G
	warpGravityConst = W
}

func SetVelDistWarpK(k float64) {
	VelDistWarpK = k
}

func SetWarpGravThreshold(v float64) {
	warpGravThreshold = v
}

//gravity acceleration (g) from planet with given mass at given range
func Gravity(mass, lenSqr, zDist float64) float64 {
	d2 := lenSqr + zDist*zDist

	if d2 == 0 {
		return 0
	}

	return gravityConst * mass / d2

	//d2 = d2 * d2
	//return gravityConst * mass * lenSqr / d2
}

func UnGravity(mass, zDist, grav float64) float64 {
	if grav <= 0 {
		return 0
	}
	d2 := gravityConst * mass / grav
	lenSqr := d2 - zDist*zDist
	if lenSqr <= 0 {
		return 0
	}
	return math.Sqrt(lenSqr)
}

func SumGravityAcc(pos v2.V2, galaxy *Galaxy) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		G = Gravity(obj.Mass, len2, obj.GDepth)
		sumF.DoAddMul(v.Normed(), G)
	}
	return sumF
}

func SumGravityAccWithReport(pos v2.V2, galaxy *Galaxy, reportLevel float64) (sumF v2.V2, report []v2.V2) {
	var v v2.V2
	var len2, G float64
	report = make([]v2.V2, 0, len(galaxy.Ordered))
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		G = Gravity(obj.Mass, len2, obj.GDepth)
		sumF.DoAddMul(v.Normed(), G)
		if G > reportLevel {
			report = append(report, v.Normed().Mul(G))
		}
	}
	return sumF, report
}

func WarpGravity(mass, lenSqr, zDist float64) float64 {

	d2 := lenSqr + zDist*zDist
	if d2 == 0 {
		return 0
	}
	if warpGravThreshold > 0 && mass < d2*warpGravThreshold {
		return 0
	}
	return warpGravityConst * mass / d2
}

func SumWarpGravityAcc(pos v2.V2, galaxy *Galaxy) (sumF v2.V2) {
	var v v2.V2
	var len2, G float64
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		G = WarpGravity(obj.Mass, len2, obj.GDepth)
		if G > 0 {
			sumF.DoAddMul(v.Normed(), G)
		}
	}
	return sumF
}
