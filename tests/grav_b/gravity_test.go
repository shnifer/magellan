package main

import (
	"encoding/json"
	. "github.com/Shnifer/magellan/commons"
	"github.com/Shnifer/magellan/v2"
	"io/ioutil"
	"testing"
)

func SumGravity1(pos v2.V2, points map[string]*GalaxyPoint) (sumF v2.V2) {
	for _, obj := range points {
		v := obj.Pos.Sub(pos)
		len2 := v.LenSqr()
		F := Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), F)
	}
	return sumF
}

func SumGravity2(pos v2.V2, points map[string]*GalaxyPoint) (sumF v2.V2) {
	var v v2.V2
	var len2, F float64
	for _, obj := range points {
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		F = Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), F)
	}
	return sumF
}

func SumGravity3(pos v2.V2, points map[string]*GalaxyPoint) (sumF v2.V2) {
	var v v2.V2
	var len2, F float64
	for _, obj := range points {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		F = Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), F)
	}
	return sumF
}

func SumGravity4(pos v2.V2, points []*GalaxyPoint) (sumF v2.V2) {
	var v v2.V2
	var len2, F float64
	for _, obj := range points {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		F = Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), F)
	}
	return sumF
}

func SumGravity6(pos v2.V2, galaxy *Galaxy) (sumF v2.V2) {
	var v v2.V2
	var len2, F float64
	for _, obj := range galaxy.Ordered {
		if obj.Mass == 0 {
			continue
		}
		v = obj.Pos.Sub(pos)
		len2 = v.LenSqr()
		F = Gravity(obj.Mass, len2, obj.Size/2)
		sumF.DoAddMul(v.Normed(), F)
	}
	return sumF
}

func BenchmarkGravity1(b *testing.B) {
	galaxy := LoadGalaxy(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := v2.RandomInCircle(100)
		SumGravity1(pos, galaxy.Points)
	}
}
func BenchmarkGravity2(b *testing.B) {
	galaxy := LoadGalaxy(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := v2.RandomInCircle(100)
		SumGravity2(pos, galaxy.Points)
	}
}

func BenchmarkGravity3(b *testing.B) {
	galaxy := LoadGalaxy(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := v2.RandomInCircle(100)
		SumGravity3(pos, galaxy.Points)
	}
}

func BenchmarkGravity4(b *testing.B) {
	galaxy := LoadGalaxy(b)

	sliceOfPoint := make([]*GalaxyPoint, 0)
	for _, v := range galaxy.Points {
		//		if v.Mass != 0 {
		sliceOfPoint = append(sliceOfPoint, v)
		//		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := v2.RandomInCircle(100)
		SumGravity4(pos, sliceOfPoint)
	}
}

func BenchmarkGravity6(b *testing.B) {
	galaxy := LoadGalaxy(b)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := v2.RandomInCircle(100)
		SumGravity6(pos, galaxy)
	}
}

func LoadGalaxy(b *testing.B) (galaxy *Galaxy) {
	galaxy = new(Galaxy)
	GalaxyID := "solar"

	buf, err := ioutil.ReadFile("Galaxy_" + GalaxyID + ".json")
	if err != nil {
		b.Error("Can't open file for galaxyID ", GalaxyID)
		return
	}

	err = json.Unmarshal(buf, &galaxy)
	if err != nil {
		b.Error("can't unmarshal file for galaxy", GalaxyID)
		return
	}

	//First restore ID's
	for id, v := range galaxy.Points {
		if v.ID == "" {
			v.ID = id
			galaxy.Points[id] = v
		}
	}
	//Second - recalc lvls!
	galaxy.RecalcLvls()

	return galaxy
}
