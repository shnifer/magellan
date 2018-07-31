package queue

import "testing"

const N = 10000

func BenchmarkAppend(b *testing.B) {
	x:=make([]int, N)
	for i:=0;i<b.N;i++{
		x=append(x, i)
		x=x[1:]
	}
}
func BenchmarkCopy(b *testing.B) {
	x:=make([]int, N)
	for i:=0;i<b.N;i++{
		x=append(x[1:], i)
	}
}
