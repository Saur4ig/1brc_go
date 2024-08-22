package versions

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Saur4ig/1brc_go/types"
)

const WORKER_COUNT = 16
const MAX_LINE_LEN = 100 // in bytes
const APPROX_STATIONS_AMOUNT = 4300

type chunk struct {
	start int64
	size  int64
}

func ProcessParallelV1(path string) {
	chunks := parallelFile(path)

	resultsChan := make(chan map[uint64]*types.Res, WORKER_COUNT)
	log.Println(chunks)

	for _, chunk := range chunks {
		go part(path, chunk.start, chunk.size, resultsChan)
	}

	results := make(map[uint64]*types.Res, APPROX_STATIONS_AMOUNT)

	for i := 0; i < len(chunks); i++ {
		result := <-resultsChan
		for station, s := range result {
			resSt, ok := results[station]
			if ok {
				if resSt.Min > s.Min {
					resSt.Min = s.Min
				}
				if s.Max > resSt.Max {
					resSt.Max = s.Max
				}
				resSt.Visited += s.Visited
				resSt.Sum += s.Sum
			} else {
				results[station] = &types.Res{
					Min:     s.Min,
					Max:     s.Max,
					Sum:     s.Sum,
					Visited: s.Visited,
					Station: s.Station,
				}
			}
		}
	}

	for _, data := range results {
		mean := data.Sum / int64(data.Visited)
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", data.Station, float64(data.Min)/10, float64(mean)/10, float64(data.Max)/10)
	}
}

func part(path string, offset, size int64, res chan map[uint64]*types.Res) {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		panic(err)
	}
	f := io.LimitedReader{R: file, N: size}

	stats := make(map[uint64]*types.Res, APPROX_STATIONS_AMOUNT)

	buffer := make([]byte, BUFFER_SIZE)
	line := make([]byte, 0, MAX_LINE_LEN)

	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}

		start := 0
		for i := 0; i < n; i++ {
			if buffer[i] == '\n' {
				line = append(line, buffer[start:i]...)
				processLine(line, stats)
				line = line[:0]
				start = i + 1
			}
		}
		if start < n {
			line = append(line, buffer[start:n]...)
		}

		if err == io.EOF {
			if len(line) > 0 {
				processLine(line, stats)
			}
			break
		}
	}

	res <- stats
}

func processLine(line []byte, stats map[uint64]*types.Res) {
	sepIndex := getSemiColIndex(line)
	if sepIndex == -1 {
		return
	}
	station := hash(line[:sepIndex])
	temperature := parseToInt(line[sepIndex+1:])

	if cityStat, ok := stats[station]; ok {
		if cityStat.Min > temperature {
			cityStat.Min = temperature
		}
		if temperature > cityStat.Max {
			cityStat.Max = temperature
		}
		cityStat.Visited++
		cityStat.Sum += temperature
		return
	}

	stats[station] = &types.Res{
		Min:     temperature,
		Max:     temperature,
		Sum:     temperature,
		Station: string(line[:sepIndex]),
		Visited: 1,
	}
}

func parallelFile(path string) []chunk {
	file, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		panic(err)
	}
	size := st.Size()
	splitSize := size / WORKER_COUNT

	chunks := make([]chunk, 0, WORKER_COUNT)
	start := int64(0)
	for i := 0; i < WORKER_COUNT; i++ {
		var endOffset int64
		if i == WORKER_COUNT-1 {
			endOffset = size
		} else {
			endOffset = adjustOffset(file, start+splitSize)
		}
		chunks = append(chunks, chunk{start, endOffset - start})
		start = endOffset
	}

	return chunks
}

// move limit to the next new line
func adjustOffset(file *os.File, limit int64) int64 {
	_, err := file.Seek(limit, 0)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	offset := limit
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		offset += int64(len(line))
		break
	}

	return offset
}

// semicolon is more close to the right side, it is faster to search from the right side
func getSemiColIndex(line []byte) int {
	for i := len(line) - 4; i >= 0; i-- {
		if line[i] == ';' {
			return i
		}
	}
	return -1
}

// create a hash from bytes
func hash(b []byte) uint64 {
	var h uint64
	for i := 0; i < 8 && i < len(b); i++ {
		h = (h << 8) | uint64(b[i])
	}
	return h
}

func parseToInt(bytes []byte) int64 {
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

	temp := int64(bytes[index] - '0')
	index++
	if index < len(bytes) && bytes[index] != '.' {
		temp = temp*10 + int64(bytes[index]-'0')
		index++
	}

	if index+1 < len(bytes) && bytes[index] == '.' {
		index++ // skip the '.'
		temp = temp*10 + int64(bytes[index]-'0')
	}

	if negative {
		temp = -temp
	}

	return temp
}
