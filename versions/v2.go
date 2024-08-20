package versions

import (
	"bytes"
	"io"
	"log"
	"math"
	"os"

	"github.com/Saur4ig/1brc_go/types"
)

const BUFFER_SIZE = 1024 * 1024

func ReadFileLineByLineV2(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buffer := make([]byte, BUFFER_SIZE)
	stats := make(map[string]*types.Result)
	var line []byte
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}

		start := 0
		for i := 0; i < n; i++ {
			if buffer[i] == '\n' {
				line = append(line, buffer[start:i]...)

				sepIndex := bytes.IndexByte(line, ';')
				if sepIndex == -1 {
					continue
				}
				station := string(line[:sepIndex])
				secondPartFloat := parseTemp(line[sepIndex+1:])
				if cityStat, ok := stats[station]; ok {
					cityStat.Min = math.Min(cityStat.Min, secondPartFloat)
					cityStat.Max = math.Max(cityStat.Max, secondPartFloat)
					cityStat.Visited++
					cityStat.Sum += secondPartFloat
				} else {
					stats[station] = &types.Result{
						Min:     secondPartFloat,
						Max:     secondPartFloat,
						Sum:     secondPartFloat,
						Visited: 1,
					}
				}

				// Reset for next line
				line = nil
				start = i + 1
			}
		}
	}

	// log stats for each city
	for city, stat := range stats {
		log.Printf("City: %s, Min: %.2f, Max: %.2f, Mean: %.2f, Visited: %d\n",
			city, stat.Min, stat.Max, stat.Sum/float64(stat.Visited), stat.Visited)
	}
}

func parseTemp(bytes []byte) float64 {
	if len(bytes) == 0 {
		return 0
	}

	index := 0
	negative := false

	if bytes[index] == '-' {
		negative = true
		index++
	}

	if index >= len(bytes) {
		return 0
	}

	temp := float64(bytes[index] - '0')
	index++
	if index < len(bytes) && bytes[index] != '.' {
		temp = temp*10 + float64(bytes[index]-'0')
		index++
	}

	if index+1 < len(bytes) && bytes[index] == '.' {
		index++ // skip the '.'
		temp += float64(bytes[index]-'0') * 0.1
	}

	if negative {
		temp = -temp
	}

	return temp
}
