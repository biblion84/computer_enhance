package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/biblion84/computerEnhance/haversine/common"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
	same as buffered_goroutine_generator.go but this one use a streaming algorithm, allowing us to have no upper limit

	limit none
	at     10_000_000:
		using batch size of 1000000
		iteration: 0 generating the pairs took 44 ms
		iteration: 0 formatting the floats to strings took 256 ms
		iteration: 0 writing to output took 173 ms
		iteration: 1 generating the pairs took 31 ms
		iteration: 1 formatting the floats to strings took 297 ms
		iteration: 1 writing to output took 146 ms
		iteration: 2 generating the pairs took 31 ms
		iteration: 2 formatting the floats to strings took 244 ms
		iteration: 2 writing to output took 147 ms
		iteration: 3 generating the pairs took 24 ms
		iteration: 3 formatting the floats to strings took 252 ms
		iteration: 3 writing to output took 163 ms
		iteration: 4 generating the pairs took 27 ms
		iteration: 4 formatting the floats to strings took 253 ms
		iteration: 4 writing to output took 147 ms
		iteration: 5 generating the pairs took 25 ms
		iteration: 5 formatting the floats to strings took 259 ms
		iteration: 5 writing to output took 149 ms
		iteration: 6 generating the pairs took 25 ms
		iteration: 6 formatting the floats to strings took 251 ms
		iteration: 6 writing to output took 159 ms
		iteration: 7 generating the pairs took 25 ms
		iteration: 7 formatting the floats to strings took 280 ms
		iteration: 7 writing to output took 155 ms
		iteration: 8 generating the pairs took 33 ms
		iteration: 8 formatting the floats to strings took 255 ms
		iteration: 8 writing to output took 148 ms
		iteration: 9 generating the pairs took 26 ms
		iteration: 9 formatting the floats to strings took 289 ms
		iteration: 9 writing to output took 154 ms
		iteration: 10 generating the pairs took 0 ms
		iteration: 10 formatting the floats to strings took 0 ms
		iteration: 10 writing to output took 0 ms
		total done in 4487 ms
*/

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

// How many pairs will be generated at a time before they're written to disk
// the larger the size the bigger the memory footprint
const MAX_BATCH_SIZE = 1_000_000
const MIN_BATCH_SIZE = 10000

func main() {
	start := time.Now()
	s := time.Now()

	if len(os.Args) != 4 {
		fmt.Printf("[uniform/cluster] [seed] [number of coordinate pairs to generate]\n")
		os.Exit(0)
	}

	clusterGeneration := false
	if strings.EqualFold(os.Args[1], "cluster") {
		clusterGeneration = true
	}

	seed, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("could not convert the [seed] '%s' to a number\n", os.Args[2])
		os.Exit(0)
	}

	// seeded rand
	srand := rand.New(rand.NewSource(int64(seed)))

	pairsQuantity, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("could not convert the [number of coordinate pairs to generate] '%s' to a number\n", os.Args[3])
		os.Exit(0)
	}

	binaryFile, err := os.OpenFile(path.Join("..", "data", "haversine.bin"), os.O_CREATE|os.O_TRUNC, 0o664)
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.OpenFile(path.Join("..", "data", "haversine.json"), os.O_CREATE|os.O_TRUNC, 0o664)
	if err != nil {
		panic(err)
	}

	outputBinary := bufio.NewWriterSize(binaryFile, 4096*10)
	outputFile := bufio.NewWriterSize(jsonFile, 4096*10)
	if _, err := outputFile.WriteString("{\"pairs\":[\n"); err != nil {
		panic(err)
	}

	haversineSum := float64(0)

	outerBatchSize := pairsQuantity / 10
	if outerBatchSize > MAX_BATCH_SIZE {
		outerBatchSize = MAX_BATCH_SIZE
	} else if outerBatchSize < MIN_BATCH_SIZE {
		outerBatchSize = MIN_BATCH_SIZE
	}

	fmt.Printf("using batch size of %d\n", outerBatchSize)

	for outerBatch := 0; outerBatch <= pairsQuantity/outerBatchSize; outerBatch++ {
		pairsToGenerate := outerBatchSize
		// last iteration, pick up the scraps
		if outerBatch == pairsQuantity/outerBatchSize {
			pairsToGenerate = pairsQuantity % outerBatchSize
		}
		pairs := make([]Pair, pairsToGenerate)
		if clusterGeneration {
			clusterNumber := srand.Intn(40) + 2
			clusterReferences := make([]Pair, 0, clusterNumber)
			for i := 0; i < clusterNumber; i++ {
				clusterReferences = append(clusterReferences, Pair{
					X0: float64((i * 72) % 180), Y0: float64((i * 72) % 180),
					X1: float64((i * 123) % 360), Y1: float64((i * 123) % 360),
				})
			}
			fmt.Printf("using cluster size of %d\n", clusterNumber)
			for i := 0; i < pairsToGenerate; i++ {
				ref := clusterReferences[i%clusterNumber]
				pairs[i] = Pair{
					X0: ref.X0 + ((srand.Float64() - 1.0) / 10), Y0: ref.Y0 + ((srand.Float64() - 1.0) / 10),
					X1: ref.X1 + ((srand.Float64() - 1.0) / 10), Y1: ref.Y1 + ((srand.Float64() - 1.0) / 10),
				}
			}
		} else {
			for i := 0; i < pairsToGenerate; i++ {
				pairs[i] = Pair{
					X0: (srand.Float64() * 360) - 180, Y0: (srand.Float64() * 180) - 90,
					X1: (srand.Float64() * 360) - 180, Y1: (srand.Float64() * 180) - 90,
				}
			}
		}

		timing(&s, fmt.Sprintf("iteration: %d generating the pairs", outerBatch))

		pairsJsoned := make([]string, pairsToGenerate)

		// when tinkering with this value I saw little difference between 10 and 10_000
		batchSize := 128
		wg := sync.WaitGroup{}
		for i := 0; i <= pairsToGenerate/batchSize; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < batchSize; j++ {
					// this combined with <= on the outer loop allow us to pick up the scraps
					if i == pairsToGenerate/batchSize && j >= pairsToGenerate%batchSize {
						break
					}
					index := j + (i * batchSize)
					sb := strings.Builder{}
					sb.WriteString("\t{\"X0\":")
					sb.WriteString(strconv.FormatFloat(pairs[index].X0, 'f', 16, 64))
					sb.WriteString(`,"Y0":`)
					sb.WriteString(strconv.FormatFloat(pairs[index].Y0, 'f', 16, 64))
					sb.WriteString(`,"X1":`)
					sb.WriteString(strconv.FormatFloat(pairs[index].X1, 'f', 16, 64))
					sb.WriteString(`,"Y1":`)
					sb.WriteString(strconv.FormatFloat(pairs[index].Y1, 'f', 16, 64))
					sb.WriteString(`}`)
					pairsJsoned[index] = sb.String()
				}
			}(i)
		}
		wg.Wait()

		timing(&s, fmt.Sprintf("iteration: %d formatting the floats to strings", outerBatch))

		for i, pair := range pairsJsoned {
			hd := common.HaversineDistance(pairs[i].X0, pairs[i].Y0, pairs[i].X1, pairs[i].Y1, 6372.8)
			p(binary.Write(outputBinary, binary.LittleEndian, hd))
			haversineSum += hd
			// do not put , on the first iteration
			if outerBatch == 0 && i == 0 {
				pw(outputFile.WriteString(pair))
			} else {
				pw(outputFile.WriteString(",\n" + pair))
			}
		}
		timing(&s, fmt.Sprintf("iteration: %d writing to output", outerBatch))
	}

	p(binary.Write(outputBinary, binary.LittleEndian, haversineSum/float64(pairsQuantity)))
	p(outputBinary.Flush())
	pw(outputFile.WriteString("\n]}"))
	p(outputFile.Flush())

	fmt.Printf("total done in %d ms\n", time.Since(start).Milliseconds())
	fmt.Printf("haversine average: %f\n", haversineSum/float64(pairsQuantity))

}

// panic write
func pw(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

func p(err error) {
	if err != nil {
		panic(err)
	}
}
