package commons

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

var mutex *os.File

func StartProfile(prefix string) {
	defer LogFunc("StartProfile " + prefix)()

	f, err := os.Create(prefix + "cpu.prof")
	if err != nil {
		log.Panicln("can't create cpu profile", prefix, err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Panicln("can't start CPU profile ", err)
	}

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	return

	mutex, err = os.Create(prefix + "mutex.prof")
	if err != nil {
		log.Panicln("can't create profile mutex")
	}
	go func() {
		for {
			time.Sleep(time.Second)
			err := pprof.Lookup("mutex").WriteTo(mutex, 1)
			if err != nil {
				log.Println(err)
				break
			}
		}
	}()
}

func heap(fn string) {
	f, err := os.Create(fn)
	if err != nil {
		log.Panicln("can't create mem profile", fn)
	}
	defer f.Close()
	runtime.GC()
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		log.Panicln("can't start mem profile", err)
	}
}

func StopProfile(prefix string) {
	defer LogFunc("StopProfile " + prefix)()

	//mutex.Close()
	pprof.StopCPUProfile()

	heap(prefix + "mem.prof")
}
