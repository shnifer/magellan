package main

import (
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var mu sync.Mutex
var items = make(map[int]struct{})

func main() {
	runtime.SetMutexProfileFraction(5)
	for i := 0; i < 100000; i++ {
		go func(i int) {
			mu.Lock()
			defer mu.Unlock()

			items[i] = struct{}{}
		}(i)
	}
	f, err := os.Create("mutex.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = pprof.Lookup("mutex").WriteTo(f, 1)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 3)
}
