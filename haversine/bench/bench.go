package main

import (
	"fmt"
	"github.com/biblion84/computerEnhance/haversine/timer"
)

func main() {

	timer.Begin()

	smallestCycles := 1 << 62
	for i := 0; i < 100_000; i++ {
		start := timer.Rdtscp()
		timer.Rdtscp()
		timer.Rdtscp()
		end := timer.Rdtscp()
		if end-start < smallestCycles {
			smallestCycles = end - start
		}
	}
	timer.End()

	timer.Print()

	fmt.Printf("smallest cycles taken by gotsc : %d\n", smallestCycles)
}
