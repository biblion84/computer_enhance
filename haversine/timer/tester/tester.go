package main

import "github.com/biblion84/computerEnhance/haversine/timer"

const TESTER_RUN = 10

func main() {

	for i := 0; i < 10; i++ {
		
	}
}

type TesterState int

const (
	TesterRunning TesterState = iota
	TesterDone
)

type RepetitionTester struct {
	State        TesterState
	Cycles       []int
	currentCycle int
}

func (t *RepetitionTester) Start() {
	t.State = TesterRunning
	t.currentCycle = timer.Rdtscp()
}

func (t *RepetitionTester) Stop() {
	stopCycle := timer.Rdtscp()
	t.Cycles = append(t.Cycles, stopCycle-t.currentCycle)

	t.currentCycle = 0

	if len(t.Cycles) >= TESTER_RUN {
		t.State = TesterDone
	}
}
