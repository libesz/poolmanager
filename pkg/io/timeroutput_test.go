package io

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	dummyOutput := DummyOutput{Name_: "dummy"}
	eightpm := func() time.Time {
		return time.Date(2020, 04, 15, 20, 0, 0, 0, time.Local)
	}
	ninepm := func() time.Time {
		return time.Date(2020, 04, 15, 21, 0, 0, 0, time.Local)
	}
	tenpm := func() time.Time {
		return time.Date(2020, 04, 15, 22, 0, 0, 0, time.Local)
	}
	tomorrow := func() time.Time {
		return time.Date(2020, 04, 16, 01, 0, 0, 0, time.Local)
	}
	timer := NewTimerOutput("timer", &dummyOutput, eightpm)
	if timer.Get() {
		t.Fatal("Timer output shall be initially off")
	}
	setResult := timer.Set(false)
	if setResult {
		t.Fatal("Set() from false to false shall not result in true")
	}
	if timer.Value() != 0.0 {
		t.Fatal("Timer shall be zero initially")
	}
	if dummyOutput.Value {
		t.Fatal("Dummy output shall not set to true")
	}
	setResult = timer.Set(true)
	if !setResult {
		t.Fatal("Set() from false to true shall not result in false")
	}
	if !dummyOutput.Value {
		t.Fatal("Dummy output shall set to true")
	}

	timer.now = ninepm
	if timer.Value() != 1.0 {
		t.Fatal("Timer shall show 1 hours time")
	}

	timer.now = tenpm
	setResult = timer.Set(false)
	if !setResult {
		t.Fatal("Set() from true to false shall result in true")
	}
	if timer.Value() != 2.0 {
		t.Fatal("Timer shall show 2 hours time")
	}
	if dummyOutput.Value {
		t.Fatal("Dummy output shall not set to true")
	}

	timer.now = tomorrow
	if timer.Value() != 0.0 {
		t.Fatalf("Timer shall show 0, but shows: %f", timer.Value())
	}
}
