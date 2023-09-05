package timer

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

func Rdtscp() int

const MAX_LABELS = 128

type RdtscTimer struct {
	cyclesPerSecond int

	lastLabel int
	labels    [MAX_LABELS]string

	total int

	currentProfile int
	profiles       [MAX_LABELS]regionProfile

	durationTimers        map[string]time.Duration
	durationRunningTimers map[string]time.Time

	called int
}
type regionProfile struct {
	parentId     int
	runningTimer int
	timer        int
}

func Profile(timerName string) {
	timer := Rdtscp()
	profileId := t.getLabelIndex(timerName)
	profile := t.profiles[profileId]

	if profile.runningTimer == 0 {
		profile.parentId = t.currentProfile
		t.currentProfile = profileId

		profile.runningTimer = timer
		if profile.parentId != 0 {
			parentProfile := t.profiles[profile.parentId]
			parentProfile.Pause(timer)
			t.profiles[profile.parentId] = parentProfile
		}
	} else {
		t.currentProfile = profile.parentId
		// mean we want to stop the profile
		profile.timer += timer - profile.runningTimer
		profile.runningTimer = 0
		if profile.parentId != 0 {
			parentProfile := t.profiles[profile.parentId]
			parentProfile.UnPause(timer)
			t.profiles[profile.parentId] = parentProfile
		}
	}

	t.profiles[profileId] = profile
	t.called++
}

func (p *regionProfile) Pause(rdtscp int) {
	p.timer += rdtscp - p.runningTimer
	p.runningTimer = 0
}

func (p *regionProfile) UnPause(rdtscp int) {
	p.runningTimer = rdtscp
}

func (t *RdtscTimer) getLabelIndex(label string) int {
	for i := 1; i <= t.lastLabel; i++ {
		if t.labels[i] == label {
			return i
		}
	}
	t.lastLabel++
	if t.lastLabel >= MAX_LABELS {
		panic("profiler: reached the max label number")
	}
	t.labels[t.lastLabel] = label
	return t.lastLabel
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
		durationTimers:        make(map[string]time.Duration, 100),
		durationRunningTimers: make(map[string]time.Time, 100),
	}
}

const MEASURE_CYCLES = true
const PROFILE = true

func Print() {

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()

	if MEASURE_CYCLES {
		fmt.Fprintf(w, "total time: \t %s µs \t total cycles : \t %s \t profiler called %s times\n",
			prettyPrint(t.cyclesToMicroSeconds(t.total)), prettyPrint(t.total), prettyPrint(t.called))
		for i := 0; i < t.lastLabel; i++ {
			label := t.labels[i]
			profile := t.profiles[i]

			percentOfTotal := (float64(profile.timer) / float64(t.total)) * 100
			fmt.Fprintf(w, "%s: \t %s \t µs\t %s \t cycles \t %.2f %% \n",
				label, prettyPrint(t.cyclesToMicroSeconds(profile.timer)), prettyPrint(profile.timer), percentOfTotal)
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
