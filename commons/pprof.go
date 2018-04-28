package commons

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func StartProfile(cpufn string) {
	if cpufn != "" {
		f, err := os.Create(cpufn)
		if err != nil {
			log.Panicln("can't create cpu profile", cpufn, err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Panicln("can't start CPU profile ", err)
		}
	}
}

func StopProfile(cpufn string, memfn string) {
	if cpufn != "" {
		pprof.StopCPUProfile()
	}

	if memfn != "" {
		f, err := os.Create(memfn)
		if err != nil {
			log.Panicln("can't create mem profile", memfn)
		}
		runtime.GC()
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			log.Panicln("can't start mem profile", err)
		}
		f.Close()
	}
}
