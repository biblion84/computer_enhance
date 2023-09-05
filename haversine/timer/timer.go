package timer

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

func Rdtscp() int

type RdtscTimer struct {
	cyclesPerSecond int
	runningTimers   map[string]int
	timers          map[string]int
	total           int

	durationTimers        map[string]time.Duration
	durationRunningTimers map[string]time.Time

	called int
}

var t RdtscTimer

func init() {
	startTimer := time.Now()
	startRdtsc := Rdtscp()
	secondDiviser := 100
	for time.Since(startTimer) < (time.Second / time.Duration(secondDiviser)) {

	}
	endRdtsc := Rdtscp()

	t = RdtscTimer{
		cyclesPerSecond:       (endRdtsc - startRdtsc) * secondDiviser,
		runningTimers:         make(map[string]int, 100),
		timers:                make(map[string]int, 100),
		durationTimers:        make(map[string]time.Duration, 100),
		durationRunningTimers: make(map[string]time.Time, 100),
	}
}

const MEASURE_CYCLES = true

func Profile(timerName string) {
	t.called++
	if MEASURE_CYCLES {
		if t.runningTimers[timerName] != 0 {
			t.timers[timerName] += Rdtscp() - t.runningTimers[timerName]
			t.runningTimers[timerName] = 0
		} else {
			t.runningTimers[timerName] = Rdtscp()
		}
	} else {
		if !t.durationRunningTimers[timerName].IsZero() {
			t.durationTimers[timerName] += time.Since(t.durationRunningTimers[timerName])
			t.durationRunningTimers[timerName] = time.Time{}
		} else {
			t.durationRunningTimers[timerName] = time.Now()
		}
	}

}

func Print() {

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()

	if MEASURE_CYCLES {
		fmt.Fprintf(w, "total time: \t %s µs \t total cycles : \t %s \t profiler called %s times\n",
			prettyPrint(t.cyclesToMicroSeconds(t.total)), prettyPrint(t.total), prettyPrint(t.called))
		for k, v := range t.timers {
			if k == "total" {
				continue
			}
			percentOfTotal := (float64(v) / float64(t.total)) * 100
			fmt.Fprintf(w, "%s: \t %s \t µs\t %s \t cycles \t %.2f %% \n",
				k, prettyPrint(t.cyclesToMicroSeconds(v)), prettyPrint(v), percentOfTotal)
		}
	} else {
		fmt.Fprintf(w, "total time \t %s µs \t profiler called %s times\n",
			prettyPrint(int(t.durationTimers["total"].Microseconds())), prettyPrint(t.called))
		for k, v := range t.durationTimers {
			if k == "total" {
				continue
			}
			percentOfTotal := (float64(v.Microseconds()) / float64(t.durationTimers["total"].Microseconds())) * 100
			fmt.Fprintf(w, "%s: \t %s \t µs \t %.2f %% \n",
				k, prettyPrint(int(v.Microseconds())), percentOfTotal)
		}
	}
}

func (t RdtscTimer) cyclesToMicroSeconds(cycles int) int {
	return int((float64(cycles) / float64(t.cyclesPerSecond)) * 1_000_000)
}

func TimeFunction(callerName string) func() {
	Profile(callerName)
	return func() {
		Profile(callerName)
	}
}

func prettyPrint(x int) string {
	printed := []rune(strconv.Itoa(x))

	prettyPrinted := ""

	for i := 0; i < len(printed); i++ {
		if i%3 == 0 && i != 0 {
			prettyPrinted = "_" + prettyPrinted
		}

		prettyPrinted = string(printed[len(printed)-1-i]) + prettyPrinted
	}

	return prettyPrinted
}

func Begin() {
	Profile("total")
	t.total = Rdtscp()
}

func End() {
	Profile("total")
	t.total = Rdtscp() - t.total
}
