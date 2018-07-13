package main

import (
	"testing"
	"math"
)

func Benchmark_pow(b *testing.B){
	for i:=0;i<b.N;i++{
		x:=math.Pow(123, 3)
		_=x
	}
}

func Benchmark_mul(b *testing.B){
	for i:=0;i<b.N;i++{
		x:=123*123*123
		_=x
	}
}
