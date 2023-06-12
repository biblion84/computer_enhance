package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
)

func main() {

	file, err := os.Open(path.Join("..", "data", "haversine.bin"))
	if err != nil {
		panic(err)
	}

	var haversineSum float64

	for {
		var h float64
		if err := binary.Read(file, binary.LittleEndian, &h); err != nil {
			if !errors.Is(err, io.EOF) {
				panic(err)
			}
			break
		}
		haversineSum = h
	}

	fmt.Printf("haversine sum :%f\n", haversineSum)
}
