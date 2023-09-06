//go:build amd64

package main

import (
	"fmt"
	"math"
	"time"
)

func Rdtscp() int

const ITERATIONS = 10_000

func main() {

	timer := NewRdtscTimer()
	times := [ITERATIONS]int64{}
	tscs := [ITERATIONS]int{}

	timer.Time("iterations")

	overhead := 1<<63 - 1
	for iteration := 0; iteration < ITERATIONS; iteration++ {
		startTime := time.Now()
		start := Rdtscp()

		for i := 0; i < 3; i++ {
			e := Rdtscp()
			s := Rdtscp()

			if e-s < overhead {
				overhead = s - e
			}
		}

		end := Rdtscp()
		endTime := time.Now()

		times[iteration] = endTime.Sub(startTime).Nanoseconds()
		tscs[iteration] = end - start
	}

	timer.Time("iterations")

	fmt.Printf("smallest rdtscp overhead : %d\n", overhead)

	timesMax := int64(0)
	timesMin := 1 << 62
	timesTotal := int64(0)
	for _, t := range times {
		if t > timesMax {
			timesMax = t
		}
		if int(t) < timesMin {
			timesMin = int(t)
		}
		timesTotal += t
	}
	timesMean := timesTotal / ITERATIONS

	tscMax := 0
	tscMin := 1 << 62
	tscTotal := 0
	for _, t := range tscs {
		if t > tscMax {
			tscMax = t
		}
		if t < tscMin {
			tscMin = t
		}
		tscTotal += t
	}
	tscMean := tscTotal / ITERATIONS

	timesDeviation := int64(0)
	for _, t := range times {
		d := t - timesMean
		timesDeviation += d * d
	}

	tscDeviation := 0
	for _, t := range tscs {
		d := t - tscMean
		tscDeviation += d * d
	}
	timesStandardDeviation := math.Sqrt(float64(timesDeviation / ITERATIONS))
	tscStandardDeviation := math.Sqrt(float64(tscDeviation / ITERATIONS))

	timesDeviationPercent := timesStandardDeviation / float64(timesMean)
	tscDeviationPercent := tscStandardDeviation / float64(tscMean)

	fmt.Printf("times \tvariance :%f, \t %%deviation : %f, \tmean : %d, \ttotal : %d, \tmin : %d, \t max : %d\n",
		timesStandardDeviation, timesDeviationPercent, timesMean, timesTotal, timesMin, timesMax)
	fmt.Printf("tsc \tvariance :%f, \t %%deviation : %f, \tmean : %d, \ttotal : %d, \tmin : %d, \t max : %d\n",
		tscStandardDeviation, tscDeviationPercent, tscMean, tscTotal, tscMin, tscMax)

	timer.Time("second")
	st := time.Now()
	for time.Since(st) < time.Second {
	}
	timer.Time("second")

	timer.Print()

}

type RdtscTimer struct {
	cyclesPerSecond int
	runningTimers   map[string]int
	Timers          map[string]int
}

func NewRdtscTimer() RdtscTimer {
	startTimer := time.Now()
	startRdtsc := Rdtscp()
	secondDiviser := 100
	for time.Since(startTimer) < (time.Second / time.Duration(secondDiviser)) {

	}
	endRdtsc := Rdtscp()

	return RdtscTimer{
		cyclesPerSecond: (endRdtsc - startRdtsc) * secondDiviser,
		runningTimers:   map[string]int{},
		Timers:          map[string]int{},
	}
}

func (t *RdtscTimer) Time(timerName string) {
	if t.runningTimers[timerName] != 0 {
		t.Timers[timerName] += Rdtscp() - t.runningTimers[timerName]
		t.runningTimers[timerName] = 0
	} else {
		t.runningTimers[timerName] = Rdtscp()
	}
}

func (t *RdtscTimer) Print() {
	for k, v := range t.Timers {
		fmt.Printf("cycles : %d, cyclesPerSecond : %d\t\t", v, t.cyclesPerSecond)
		fmt.Printf("%s : %dµs\t %d\n", k, t.cyclesToMicroSeconds(v), v)
	}
}

func (t *RdtscTimer) cyclesToMicroSeconds(cycles int) int {
	return int((float64(cycles) / float64(t.cyclesPerSecond)) * 1_000_000)
}
