package versions

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/Saur4ig/1brc_go/types"
)

const WORKER_COUNT = 16
const MAX_LINE_LEN = 100 // in bytes

type chunk struct {
	start int64
	size  int64
}

func ProcessParallelV1(path string) {
	chunks := parallelFile(path)

	resultsChan := make(chan map[string]*types.Result, WORKER_COUNT)
	log.Println(chunks)

	for _, chunk := range chunks {
		go part(path, chunk.start, chunk.size, resultsChan)
	}

	results := make(map[string]*types.Result)

	for i := 0; i < len(chunks); i++ {
		result := <-resultsChan
		for station, s := range result {
			resSt, ok := results[station]
			if ok {
				resSt.Min = math.Min(resSt.Min, s.Min)
				resSt.Max = math.Max(resSt.Max, s.Max)
				resSt.Visited += s.Visited
				resSt.Sum += s.Sum
			} else {
				results[station] = &types.Result{
					Min:     s.Min,
					Max:     s.Max,
					Sum:     s.Sum,
					Visited: s.Visited,
				}
			}
		}
	}

	for station, data := range results {
		mean := data.Sum / float64(data.Visited)
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", station, data.Min, mean, data.Max)
	}
}

func part(path string, offset, size int64, res chan map[string]*types.Result) {
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

	stats := make(map[string]*types.Result)

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

func processLine(line []byte, stats map[string]*types.Result) {
	sepIndex := bytes.IndexByte(line, ';')
	if sepIndex == -1 {
		return
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
