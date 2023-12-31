package timer

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

const MEGABYTE = 1024 * 1024
const GIGABYTE = MEGABYTE * 1024

func Rdtscp() int

func BusySleep(duration time.Duration) {
	cyclesToWait := int(duration.Seconds() * float64(t.cyclesPerSecond))
	cyclesElapsed := 0
	start := Rdtscp()
	for {
		end := Rdtscp()
		cyclesElapsed += end - start

		if cyclesElapsed >= cyclesToWait {
			return
		}
		start = end
	}
}

const MAX_LABELS = 128

type RdtscTimer struct {
	pageFaultHandler PageFaultHandler
	cyclesPerSecond  int

	lastLabel int
	labels    [MAX_LABELS]string

	// The profile currently being timed
	// this allows a 'child' profile to know that it have a parent, see Profile
	currentProfile int
	profiles       [MAX_LABELS]regionProfile

	called      int
	totalCycles int
}
type regionProfile struct {
	parentId     int
	runningTimer int
	timer        int

	runningTimerWithChildren int
	timerWithChildren        int
	bytes                    int

	runningPageFaults int
	pageFaults        int
}

const PROFILE = true

func Profile(timerName string, bytes ...int) {
	if !PROFILE {
		return
	}
	if timerName == "main loop" {
		fmt.Println("main loop")
	}
	timer := Rdtscp()
	pageFaults := t.pageFaultHandler.GetPageFaults()
	profileId := t.getLabelIndex(timerName)
	profile := t.profiles[profileId]

	if profile.runningTimer == 0 {
		profile.parentId = t.currentProfile
		t.currentProfile = profileId

		profile.runningPageFaults = pageFaults
		profile.runningTimer = timer
		profile.runningTimerWithChildren = timer
		if profile.parentId != 0 {
			parentProfile := t.profiles[profile.parentId]
			parentProfile.Pause(timer, pageFaults)
			t.profiles[profile.parentId] = parentProfile
		}
	} else {
		// mean we want to stop the profile
		t.currentProfile = profile.parentId
		profile.pageFaults += pageFaults - profile.runningPageFaults
		profile.timer += timer - profile.runningTimer
		profile.timerWithChildren += timer - profile.runningTimerWithChildren
		profile.runningTimer = 0
		if profile.parentId != 0 {
			parentProfile := t.profiles[profile.parentId]
			parentProfile.UnPause(timer, pageFaults)
			t.profiles[profile.parentId] = parentProfile
		}
	}

	processedBytes := 0
	for _, b := range bytes {
		processedBytes += b
	}

	profile.bytes += processedBytes
	t.profiles[profileId] = profile
	t.called++
}

func (p *regionProfile) Pause(rdtscp int, pageFaults int) {
	p.timer += rdtscp - p.runningTimer
	p.pageFaults += pageFaults - p.runningPageFaults
}

func (p *regionProfile) UnPause(rdtscp int, pageFaults int) {
	p.runningTimer = rdtscp
	p.runningPageFaults = pageFaults
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

	pageFaultHandler, err := InitPageFaultHandler()
	if err != nil {
		panic(err)
	}

	t = RdtscTimer{
		cyclesPerSecond:  (endRdtsc - startRdtsc) * secondDiviser,
		pageFaultHandler: pageFaultHandler,
	}

}

func Print() {

	w := &tabwriter.Writer{}
	w.Init(os.Stdout, 8, 8, 0, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "totalCycles time: \t %s µs \t totalCycles cycles : \t %s \t profiler called %s times\n",
		prettyPrint(t.cyclesToMicroSeconds(t.totalCycles)), prettyPrint(t.totalCycles), prettyPrint(t.called))
	for i := 1; i <= t.lastLabel; i++ {
		label := t.labels[i]
		profile := t.profiles[i]

		percentOfTotal := (float64(profile.timer) / float64(t.totalCycles)) * 100
		percentOfTotalWithChildren := (float64(profile.timerWithChildren) / float64(t.totalCycles)) * 100

		fmt.Fprintf(w, "%s: \t %s \t µs\t ", label, prettyPrint(t.cyclesToMicroSeconds(profile.timer)))
		fmt.Fprintf(w, "%s \t cycles \t %.2f %% \t with children: %.2f %%\t", prettyPrint(profile.timer),
			percentOfTotal, percentOfTotalWithChildren)

		if profile.bytes > 0 {
			seconds := t.cyclesToSeconds(profile.timer)
			bytesPerSecond := float64(profile.bytes) / seconds

			megabytes := float64(profile.bytes) / MEGABYTE
			gigabytesPerSecond := bytesPerSecond / GIGABYTE

			fmt.Fprintf(w, "%.3fmb at %.2fgb/s\t", megabytes, gigabytesPerSecond)

		}

		if profile.pageFaults > 0 {
			fmt.Fprintf(w, "%d page faults\t", profile.pageFaults)
		}

		fmt.Fprintf(w, "\n")
	}
}

func (t RdtscTimer) cyclesToMicroSeconds(cycles int) int {
	return int((float64(cycles) / float64(t.cyclesPerSecond)) * 1_000_000)
}

func (t RdtscTimer) cyclesToSeconds(cycles int) float64 {
	return float64(cycles) / float64(t.cyclesPerSecond)
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
	t.totalCycles = Rdtscp()
}

func End() {
	t.totalCycles = Rdtscp() - t.totalCycles
}
