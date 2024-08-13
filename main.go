package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
)

type TempData struct {
	city        string
	temperature float32
}

type Result struct {
	min     float32
	max     float32
	mean    float32
	visited int
}

func main() {
	path := "../1bcfile/measurements.txt"

	// profile files
	cpuProfile := "cpu.prof"
	memProfile := "mem.prof"

	// start CPU profiling
	startCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	// read the file and process data
	readFileLineByLine(path)

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

// readFileLineByLine reads a file and computes statistics for each city
func readFileLineByLine(logfile string) {
	file, err := os.OpenFile(logfile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// map to store statistics for each city
	stats := make(map[string]Result)

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		temp := decodeLine(sc.Text())
		if cityStat, ok := stats[temp.city]; ok {
			cityStat.min = float32(math.Min(float64(cityStat.min), float64(temp.temperature)))
			cityStat.max = float32(math.Max(float64(cityStat.max), float64(temp.temperature)))
			cityStat.visited++
			cityStat.mean = cityStat.mean + (temp.temperature-cityStat.mean)/float32(cityStat.visited)
			stats[temp.city] = cityStat
		} else {
			stats[temp.city] = Result{
				min:     temp.temperature,
				max:     temp.temperature,
				mean:    temp.temperature,
				visited: 1,
			}
		}
	}

	log.Printf("Found cities - %d\n", len(stats))
	if err := sc.Err(); err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	// log stats for each city
	for city, stat := range stats {
		log.Printf("City: %s, Min: %.2f, Max: %.2f, Mean: %.2f, Visited: %d\n",
			city, stat.min, stat.max, stat.mean, stat.visited)
	}
}

// decodeLine decodes a line of the log file into TempData.
func decodeLine(line string) TempData {
	parts := strings.Split(line, ";")
	temperature, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 32)
	return TempData{
		city:        parts[0],
		temperature: float32(temperature),
	}
}
