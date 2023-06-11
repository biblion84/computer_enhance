package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

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

const BatchSize = 100_000

func main() {
	start := time.Now()

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

	pairsQuantity, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("could not convert the [number of coordinate pairs to generate] '%s' to a number\n", os.Args[2])
		os.Exit(0)
	}

	s := time.Now()
	pairs := make([]Pair, pairsQuantity)

	fmt.Printf("make took %d ms\n", time.Since(s).Milliseconds())
	s = time.Now()

	outputFile, err := os.OpenFile("haversine.json", os.O_CREATE|os.O_TRUNC, 0o664)
	if err != nil {
		panic(err)
	}

	for b := 0; b < pairsQuantity/BatchSize; b++ {
		for i := 0; i < len(pairs); i++ {
			index := (b * BatchSize) + i
			pairs[index] = Pair{
				X0: (srand.Float64() * 360) - 180, Y0: (srand.Float64() * 180) - 90,
				X1: (srand.Float64() * 360) - 180, Y1: (srand.Float64() * 180) - 90,
			}
		}

		outputJson, err := json.Marshal(pairs)
		if err != nil {
			panic(err)
		}

		if err := outputFile.Write("haversine.json", outputJson, 0o664); err != nil {
			panic(err)
		}
	}

	outputFile.
		fmt.Printf("write outputFile took %d ms\n", time.Since(s).Milliseconds())
	s = time.Now()

	fmt.Printf("done in %d ms\n", time.Since(start).Milliseconds())
}
