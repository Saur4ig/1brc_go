package versions

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/Saur4ig/1brc_go/types"
)

// ReadFileLineByLine reads a file and computes statistics for each city
func ReadFileLineByLine(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	// map to store statistics for each city
	stats := make(map[string]types.Result)

	sc := bufio.NewScanner(file)
	sc.Buffer(make([]byte, 1024*1024), 1024*1024)
	for sc.Scan() {
		temp := decodeLine(sc.Bytes())
		if cityStat, ok := stats[temp.City]; ok {
			cityStat.Min = math.Min(cityStat.Min, temp.Temperature)
			cityStat.Max = math.Max(cityStat.Max, temp.Temperature)
			cityStat.Visited++
			cityStat.Sum += temp.Temperature
			stats[temp.City] = cityStat
		} else {
			stats[temp.City] = types.Result{
				Min:     temp.Temperature,
				Max:     temp.Temperature,
				Sum:     temp.Temperature,
				Visited: 1,
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
			city, stat.Min, stat.Max, stat.Sum/float64(stat.Visited), stat.Visited)
	}
}

// decodeLine decodes a line of the log file into TempData.
func decodeLine(line []byte) types.TempData {
	parts := strings.Split(string(line), ";")
	temperature, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 32)
	return types.TempData{
		City:        parts[0],
		Temperature: temperature,
	}
}
