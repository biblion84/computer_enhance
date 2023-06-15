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
	"time"
)

func main() {

	filename := path.Join("data", "haversine.json")

	start := time.Now()

	file, err := os.Open(filename)
	p(err)

	fileinfo, err := os.Stat(filename)
	p(err)

	// we only want to parse the haversine.json for the computer_enhance lesson
	// this is not a valid json parser

	var elapsedParsing time.Duration

	sum := float64(0)

	pairs := []models.Pair{}

	// first we try to find the first value
	buffer := make([]byte, 4096)

	toRead := [4]string{}

	nesting := 3
	inString := false
	rightOfColumn := false
	// the value index that we're currently reading
	currentlyReading := 0
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
			switch c {
			case '{', '[':
				nesting--
				rightOfColumn = false
				inString = false
				currentlyReading = 0
			case '}', ']':
				nesting++
				rightOfColumn = false
				inString = false
				if currentlyReading == len(toRead)-1 {
					s := time.Now()
					read := make([]float64, len(toRead))

					for i, r := range toRead {
						value, err := strconv.ParseFloat(r, 64)
						p(err)
						read[i] = value
					}

					pair := models.Pair{
						X0: read[0],
						Y0: read[1],
						X1: read[2],
						Y1: read[3],
					}

					pairs = append(pairs, pair)
					sum += common.HaversineDistance(pair.X0, pair.Y0, pair.X1, pair.Y1, 6372.8)
					toRead = [4]string{}
					elapsedParsing += time.Since(s)
				}
				currentlyReading = 0
			}
			if nesting != 0 {
				continue
			}
			if c == '"' {
				inString = !inString
			} else if c == ':' {
				rightOfColumn = true
			} else if c == ',' {
				rightOfColumn = false
				currentlyReading++
			}

			if rightOfColumn && (c >= '0' && c <= '9' || c == '.' || c == '-') {
				toRead[currentlyReading] += string(c)
			}

		}

	}
	fmt.Printf("time parsing %d ms\n", elapsedParsing.Milliseconds())

	timeElapsed := time.Since(start)
	mbPerSecond := int(((float64(fileinfo.Size()) / 1_000_000.0) / float64(timeElapsed.Milliseconds())) * 1000)

	fmt.Printf("took %d ms, %d MB/s\n", timeElapsed.Milliseconds(), mbPerSecond)
	fmt.Printf("haversine average : %f\n", sum/float64(len(pairs)))
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
