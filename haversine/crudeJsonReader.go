package main

import (
	"errors"
	"fmt"
	"github.com/biblion84/computerEnhance/haversine/common"
	"github.com/biblion84/computerEnhance/haversine/models"
	"github.com/biblion84/computerEnhance/haversine/timer"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

func main() {
	timer.Begin()

	filename := path.Join("data", "haversine.json")

	timer.Profile("os.Open")
	file, err := os.Open(filename)
	timer.Profile("os.Open")

	p(err)

	timer.Profile("os.Stat")
	fileinfo, err := os.Stat(filename)
	timer.Profile("os.Stat")
	p(err)

	// we only want to parse the haversine.json for the computer_enhance lesson
	// this is not a valid json parser

	timer.Profile("setup")
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
	// the value index that we're currently reading, because we're reading 4 floats (x0, y0, x1, y1) it can be 0 through 3
	currentlyReading := 0
	total := 0
	index := 0

	timer.Profile("setup")

	timer.Profile("main loop")
	for {
		timer.Profile("file.Read")
		_, err := file.Read(buffer)
		if err != nil {
			if errors.Is(io.EOF, err) {
				break
			} else {
				p(err)
			}
		}
		timer.Profile("file.Read")

		timer.Profile("ifs")
		for _, c := range buffer {
			if c == '{' || c == '[' {
				targetNesting--
				rightOfColumn = false
				currentlyReading = 0
			} else if c == '}' || c == ']' {
				targetNesting++
				rightOfColumn = false
				if currentlyReading == len(toRead)-1 {
					total++

					s := time.Now()
					read := make([]float64, len(toRead))
					for i, r := range toRead {
						timer.Profile("parseFloat")
						value, err := strconv.ParseFloat(string(r), 64)
						timer.Profile("parseFloat")
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
		timer.Profile("ifs")

	}
	timer.Profile("main loop")

	fmt.Printf("time parsing %d ms\n", elapsedParsing.Milliseconds())

	timeElapsed := time.Since(start)
	mbPerSecond := int(((float64(fileinfo.Size()) / 1_000_000.0) / float64(timeElapsed.Milliseconds())) * 1000)

	// We over-allocated pairs, now bringing it back to the size it should have been
	// but we didn't know how much elements were in the json, thus the estimation of the len above
	pairs = pairs[:total]
	haversineSum := float64(0)
	timer.Profile("haversineDistance")
	for _, pair := range pairs {
		haversineSum += common.HaversineDistance(pair.X0, pair.Y0, pair.X1, pair.Y1, 6372.8)
	}
	timer.Profile("haversineDistance")

	fmt.Printf("took %d ms, %d MB/s\n", timeElapsed.Milliseconds(), mbPerSecond)
	fmt.Printf("haversine average : %f\n", haversineSum/float64(len(pairs)))

	timer.End()
	timer.Print()
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
