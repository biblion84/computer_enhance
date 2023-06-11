package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

/*
	Slower than the naive one.
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

	file, err := os.OpenFile("haversine.json", os.O_CREATE|os.O_TRUNC, 0o664)
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteString(`{"pairs":[`); err != nil {
		panic(err)
	}

	pairsQuantity, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("could not convert the [number of coordinate pairs to generate] '%s' to a number\n", os.Args[2])
		os.Exit(0)
	}

	pairs := make([]Pair, pairsQuantity)

	for i := 0; i < len(pairs); i++ {
		if _, err := file.WriteString(fmt.Sprintf(`{"X0":%f,"Y0":%f,"X1":%f,"Y1":%f}`,
			(srand.Float64()*360)-180, (srand.Float64()*180)-90,
			(srand.Float64()*360)-180, (srand.Float64()*180)-90)); err != nil {
			panic(err)
		}
		if i != len(pairs)-1 {
			if _, err := file.WriteString(","); err != nil {
				panic(err)
			}
		}
	}

	if _, err := file.WriteString("]}"); err != nil {
		panic(err)
	}

	fmt.Printf("done in %d ms\n", time.Since(start).Milliseconds())
}
