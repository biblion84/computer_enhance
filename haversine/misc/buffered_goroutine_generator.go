package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	faster than the naive one

	limit 100_000_000 (oom)
	at     10_000_000:
		generating the pairs took 377 ms
		formatting the floats to strings took 2453 ms
		writing to output took 1390 ms
		total done in 4221 ms
*/

func radianFromDegrees(degrees float64) float64 {
	return degrees * 0.01745329251994329577
}

func square(a float64) float64 {
	return a * a
}

func haversineDistance(x0, y0, x1, y1, radius float64) float64 {

	dLat := radianFromDegrees(y1 - y0)
	dLon := radianFromDegrees(x1 - x0)
	lat1 := radianFromDegrees(y0)
	lat2 := radianFromDegrees(y1)

	a := square(math.Sin(dLat/2)) + math.Cos(lat1)*math.Cos(lat2)*square(math.Sin(dLon/2))
	c := 2 * math.Asin(math.Sqrt(a))

	return radius * c
}

type Pair struct {
	X0 float64
	Y0 float64
	X1 float64
	Y1 float64
}

func timing(timer *time.Time, message string) {
	fmt.Printf("%s took %d ms\n", message, time.Since(*timer).Milliseconds())
	*timer = time.Now()
}

func main() {
	start := time.Now()
	s := time.Now()

	if len(os.Args) != 3 {
		fmt.Printf("[seed] [number of coordinate pairs to generate]\n")
		os.Exit(0)
	}

	seed, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("could not convert the [seed] '%s' to a number\n", os.Args[1])
		os.Exit(0)
	}

	// seeded rand
	srand := rand.New(rand.NewSource(int64(seed)))

	file, err := os.OpenFile(path.Join("data", "haversine.json"), os.O_CREATE|os.O_TRUNC, 0o664)
	if err != nil {
		panic(err)
	}

	outputFile := bufio.NewWriterSize(file, 4096*10)
	if _, err := outputFile.WriteString("{\"pairs\":[\n"); err != nil {
		panic(err)
	}

	pairsQuantity, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("could not convert the [number of coordinate pairs to generate] '%s' to a number\n", os.Args[2])
		os.Exit(0)
	}

	pairs := make([]Pair, pairsQuantity)
	for i := 0; i < pairsQuantity; i++ {
		pairs[i] = Pair{
			X0: (srand.Float64() * 360) - 180, Y0: (srand.Float64() * 180) - 90,
			X1: (srand.Float64() * 360) - 180, Y1: (srand.Float64() * 180) - 90,
		}
	}

	timing(&s, "generating the pairs")

	pairsJsoned := make([]string, pairsQuantity)

	batchSize := 32
	wg := sync.WaitGroup{}
	for i := 0; i <= pairsQuantity/batchSize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < batchSize; j++ {
				// this combined with <= on the outer loop allow us to pick up the scraps
				if i == pairsQuantity/batchSize && j >= pairsQuantity%batchSize {
					break
				}
				index := j + (i * batchSize)
				sb := strings.Builder{}
				sb.WriteString("\t{\"X0\":")
				sb.WriteString(strconv.FormatFloat(pairs[index].X0, 'f', 10, 64))
				sb.WriteString(`,"Y0":`)
				sb.WriteString(strconv.FormatFloat(pairs[index].Y0, 'f', 10, 64))
				sb.WriteString(`,"X1":`)
				sb.WriteString(strconv.FormatFloat(pairs[index].X1, 'f', 10, 64))
				sb.WriteString(`,"Y1":`)
				sb.WriteString(strconv.FormatFloat(pairs[index].Y1, 'f', 10, 64))
				sb.WriteString(`}`)
				pairsJsoned[index] = sb.String()
			}
		}(i)
	}
	wg.Wait()

	timing(&s, "formatting the floats to strings")

	for i, p := range pairsJsoned {
		if i != len(pairsJsoned)-1 {
			pw(outputFile.WriteString(p + ",\n"))
		} else {
			pw(outputFile.WriteString(p + "\n]}"))
			break
		}
	}

	if err := outputFile.Flush(); err != nil {
		panic(err)
	}
	timing(&s, "writing to output")

	fmt.Printf("total done in %d ms\n", time.Since(start).Milliseconds())
}

// panic write
func pw(_ int, err error) {
	if err != nil {
		panic(err)
	}
}
