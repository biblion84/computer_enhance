package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
	Limit : 100_000_000 (not enough memory)
	at       10_000_000:
				make took 4 ms
				populating took 337 ms
				pair struct took 0 ms
				marshall took 9368 ms
				write file took 987 ms
				done in 10697 ms
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

	pairs := make([]Pair, pairsQuantity)

	for i := 0; i < len(pairs); i++ {
		pairs[i] = Pair{
			X0: (srand.Float64() * 360) - 180, Y0: (srand.Float64() * 180) - 90,
			X1: (srand.Float64() * 360) - 180, Y1: (srand.Float64() * 180) - 90,
		}
	}

	type PairsJson struct {
		Pairs []Pair `json:"pairs"`
	}

	output := PairsJson{
		pairs,
	}

	outputJson, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("haversine.json", outputJson, 0o664); err != nil {
		panic(err)
	}

	fmt.Printf("done in %d ms\n", time.Since(start).Milliseconds())
}
