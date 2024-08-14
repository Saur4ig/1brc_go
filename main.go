package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/Saur4ig/1brc_go/versions"
)

func main() {
	path := "../1bcfile/measurements.txt"
	//path := "data/test.txt"

	// profile files
	cpuProfile := "cpu.prof"
	memProfile := "mem.prof"

	// start CPU profiling
	startCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	// read the file and process data
	versions.ReadFileLineByLine(path)

	// capture memory profile.
	writeMemoryProfile(memProfile)
}

// startCPUProfile starts CPU profiling and saves it to the specified file
func startCPUProfile(cpuProfile string) {
	f, err := os.Create(cpuProfile)
	if err != nil {
		log.Fatalf("could not create CPU profile: %s", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatalf("could not start CPU profile: %s", err)
	}
}

// writeMemoryProfile writes the current memory profile to the specified file
func writeMemoryProfile(memProfile string) {
	mf, err := os.Create(memProfile)
	if err != nil {
		log.Fatalf("could not create memory profile: %s", err)
	}
	defer mf.Close()
	runtime.GC() // run garbage collection to update memory statistics
	if err := pprof.WriteHeapProfile(mf); err != nil {
		log.Fatalf("could not write memory profile: %s", err)
	}
}
