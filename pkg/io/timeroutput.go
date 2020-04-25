package io

import (
	"fmt"
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
		fmt.Printf("Timer %s running, value: %f\n", t.name, t.runTimeTodayInHours+t.now().Sub(t.lastStart).Hours())
		return t.runTimeTodayInHours + t.now().Sub(t.lastStart).Hours()
	}
	fmt.Printf("Timer %s not running, value: %f\n", t.name, t.runTimeTodayInHours)
	return t.runTimeTodayInHours
}

func (t *TimerOutput) resetTotalIfDayChanged() {
	if !t.lastStart.Round(24 * time.Hour).Equal(t.now().Round(24 * time.Hour)) {
		t.runTimeTodayInHours = 0
		fmt.Printf("Timer %s, date cycle\n", t.name)
	} else {
		fmt.Printf("Timer %s, no date cycle\n", t.name)
	}
}

func (t *TimerOutput) Switch(value bool) {
	t.resetTotalIfDayChanged()
	now := t.now()
	if !t.lastValue && value {
		fmt.Printf("Dummy timer %s starting\n", t.name)
		t.lastStart = now
	}
	if t.lastValue && !value {
		t.runTimeTodayInHours += now.Sub(t.lastStart).Hours()
		fmt.Printf("Dummy timer %s stopped. total run today: %f\n", t.name, t.runTimeTodayInHours)
	}
	t.lastValue = value
	t.output.Switch(value)
}
