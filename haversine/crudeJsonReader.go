package main

import (
	"errors"
	"fmt"
	"github.com/biblion84/computerEnhance/haversine/common"
	"github.com/biblion84/computerEnhance/haversine/models"
	"io"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

func main() {

	filename := path.Join("data", "haversine.json")

	file, err := os.Open(filename)
	p(err)

	fileinfo, err := os.Stat(filename)
	p(err)

	// we only want to parse the haversine.json for the computer_enhance lesson
	// this is not a valid json parser

	var elapsedParsing time.Duration

	// roughly ~100 bytes per pair
	pairs := make([]models.Pair, fileinfo.Size()/100)

	buffer := make([]byte, 4096*100)

	toRead := [4][]byte{} // [1000 iterations][4 digits]
	for i := range toRead {
		toRead[i] = make([]byte, 0, 21) // '-' + 3 digit + '.' + 16 digits of precision
	}

	start := time.Now()

	targetNesting := 3
	rightOfColumn := false
	// the value index that we're currently reading
	currentlyReading := 0
	total := 0
	wg := sync.WaitGroup{}
	index := 0
	for {
		_, err := file.Read(buffer)
		if err != nil {
			if errors.Is(io.EOF, err) {
				break
			} else {
				p(err)
			}
		}

		for _, c := range buffer {
			if c == '{' || c == '[' {
				targetNesting--
				rightOfColumn = false
				currentlyReading = 0
			} else if c == '}' || c == '}' {
				targetNesting++
				rightOfColumn = false
				if currentlyReading == len(toRead)-1 {
					wg.Add(1)
					total++
					go func(index int, toRead [4]string) {
						defer wg.Done()

						s := time.Now()
						read := make([]float64, len(toRead))
						for i, r := range toRead {
							value, err := strconv.ParseFloat(string(r), 64)
							p(err)
							read[i] = value
						}

						pair := models.Pair{
							X0: read[0],
							Y0: read[1],
							X1: read[2],
							Y1: read[3],
						}

						pairs[index] = pair
						elapsedParsing += time.Since(s)
					}(index, [4]string{string(toRead[0]), string(toRead[1]), string(toRead[2]), string(toRead[3])})

					index++
					currentlyReading = 0
					for i := range toRead {
						toRead[i] = toRead[i][:0]
					}
				}
			}

			if targetNesting != 0 {
				continue
			}
			switch c {
			case ':':
				rightOfColumn = true
			case ',':
				rightOfColumn = false
				currentlyReading++
			}

			if rightOfColumn && (c >= '0' && c <= '9' || c == '.' || c == '-') {
				toRead[currentlyReading] = append(toRead[currentlyReading], c)
			}
		}
	}
	wg.Wait()
	fmt.Printf("time parsing %d ms\n", elapsedParsing.Milliseconds())

	timeElapsed := time.Since(start)
	mbPerSecond := int(((float64(fileinfo.Size()) / 1_000_000.0) / float64(timeElapsed.Milliseconds())) * 1000)

	haversineSum := float64(0)
	for _, pair := range pairs {
		haversineSum += common.HaversineDistance(pair.X0, pair.Y0, pair.X1, pair.Y1, 6372.8)
	}

	fmt.Printf("took %d ms, %d MB/s\n", timeElapsed.Milliseconds(), mbPerSecond)
	fmt.Printf("haversine average : %f\n", haversineSum/float64(total))
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
