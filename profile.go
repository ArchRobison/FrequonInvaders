// Routines related to profiling

package main

import (
	"flag"
	"github.com/ArchRobison/Gophetica/nimble"
	"os"
	"runtime/pprof"
)

var benchmarkMode = false

var fourierFrameCount = 0 // Number of fourier frames rendered

func tallyFourierFrame() {
	if debugMode && benchmarkMode {
		if fourierFrameCount++; fourierFrameCount >= 1000 {
			nimble.Quit()
		}
	}
}

// profileStart starts any profilig requested on the command line.
// It returns a slice of functions to be executed when main exits.
func profileStart() []func() {
	fun := make([]func(), 0, 2)
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	memProfile := flag.String("memprofile", "", "write mem profile to file")
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		fun = append(fun, pprof.StopCPUProfile)
		benchmarkMode = true
	}
	if *memProfile != "" {
		heapProfileFile, err := os.Create(*memProfile)
		if err != nil {
			panic(err)
		}
		fun = append(fun, func() {
			pprof.WriteHeapProfile(heapProfileFile)
			heapProfileFile.Close()
		})
		benchmarkMode = true
	}
	return fun
}
