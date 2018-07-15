package storage

import (
	"github.com/peterbourgon/diskv"
	"testing"
	"strconv"
)

var d *disk
const storagePath = "teststorage"

func init() {
	diskOpts := diskv.Options{
		BasePath:     storagePath,
		CacheSizeMax: 1024 * 1024,
	}
	d = newDisk(diskOpts, 0)
}

func createNRecords(n int) {
	d.EraseAll()
	var s string
	for i:=0;i<n;i++{
		s=strconv.Itoa(i)
		d.Write(s, []byte(s))
	}
}

const n=1000

func Benchmark_LoadKeys(b *testing.B) {
	createNRecords(n)
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		for s:=range d.Keys(nil) {
			_=s
		}
	}
}