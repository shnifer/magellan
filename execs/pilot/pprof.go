package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func startProfile() {
	cpufn := DEFVAL.CpuProfFileName
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

func stopProfile() {
	if DEFVAL.CpuProfFileName != "" {
		pprof.StopCPUProfile()
	}

	memfn := DEFVAL.MemProfFileName
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
