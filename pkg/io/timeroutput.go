package io

import (
	"log"
	"time"
)

type TimerOutput struct {
	name                string
	output              Output
	runTimeTodayInHours float64
	lastStart           time.Time
	lastValue           bool
	now                 func() time.Time
}

func NewTimerOutput(name string, output Output, now func() time.Time) TimerOutput {
	return TimerOutput{
		name:   name,
		output: output,
		now:    now,
	}
}

func (t *TimerOutput) Type() string {
	return "time"
}

func (t *TimerOutput) Degree() string {
	return "H"
}

func (t *TimerOutput) Value() float64 {
	t.resetTotalIfDayChanged()
	if t.lastValue {
		//log.Printf("Timer %s running, value: %f\n", t.name, t.runTimeTodayInHours+t.now().Sub(t.lastStart).Hours())
		return t.runTimeTodayInHours + t.now().Sub(t.lastStart).Hours()
	}
	//log.Printf("Timer %s not running, value: %f\n", t.name, t.runTimeTodayInHours)
	return t.runTimeTodayInHours
}

func (t *TimerOutput) resetTotalIfDayChanged() {
	if !t.lastStart.Round(24 * time.Hour).Equal(t.now().Round(24 * time.Hour)) {
		log.Printf("Timer %s, date change. Total runtime yesterday: %f\n", t.name, t.runTimeTodayInHours)
		t.runTimeTodayInHours = 0
		t.lastStart = t.now()
	} else {
		//log.Printf("Timer %s, no date cycle\n", t.name)
	}
}

func (t *TimerOutput) Set(value bool) bool {
	t.resetTotalIfDayChanged()
	now := t.now()
	if !t.lastValue && value {
		log.Printf("Dummy timer %s starting\n", t.name)
		t.lastStart = now
	}
	if t.lastValue && !value {
		t.runTimeTodayInHours += now.Sub(t.lastStart).Hours()
		log.Printf("Dummy timer %s stopped. total run today: %f\n", t.name, t.runTimeTodayInHours)
	}
	t.lastValue = value
	return t.output.Set(value)
}

func (t *TimerOutput) Get() bool {
	return t.output.Get()
}
