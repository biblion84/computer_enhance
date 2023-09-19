package main

import (
	"github.com/biblion84/computerEnhance/haversine/timer"
	"github.com/dterei/gotsc"
	"testing"
)

func BenchmarkRdtscp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Rdtscp()
		timer.Rdtscp()
	}
}

func BenchmarkGoTsc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gotsc.BenchStart()
		gotsc.BenchEnd()
	}
}
