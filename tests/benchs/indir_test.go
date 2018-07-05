package main

import (
	"github.com/Shnifer/magellan/v2"
	"math"
	"testing"
)

func BenchmarkInDir(b *testing.B) {
	var ang float64
	var sum v2.V2
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sum = InDir1(ang)
		ang += 0.1
	}
	b.StopTimer()
	_ = sum
}

func InDir1(angle float64) v2.V2 {
	a := angle * v2.Deg2Rad
	return v2.V2{-math.Sin(a), math.Cos(a)}
}

func InDir2(angle float64) v2.V2 {
	a := angle * v2.Deg2Rad
	s, c := math.Sincos(a)
	return v2.V2{-s, c}
}
