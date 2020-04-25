package io

import (
	"fmt"
	"time"
)

type TimedGPIO struct {
	Name                string
	runTimeTodayInHours float64
	lastStart           time.Time
	lastValue           bool
	Now                 func() time.Time
}

func (t *TimedGPIO) Type() string {
	return "time"
}

func (t *TimedGPIO) Degree() string {
	return "H"
}

func (t *TimedGPIO) Value() float64 {
	t.resetTotalIfDayChanged()
	if t.lastValue {
		fmt.Printf("Timer %s running, value: %f\n", t.Name, t.runTimeTodayInHours+t.Now().Sub(t.lastStart).Hours())
		return t.runTimeTodayInHours + t.Now().Sub(t.lastStart).Hours()
	}
	fmt.Printf("Timer %s not running, value: %f\n", t.Name, t.runTimeTodayInHours)
	return t.runTimeTodayInHours
}

func (t *TimedGPIO) resetTotalIfDayChanged() {
	if !t.lastStart.Round(24 * time.Hour).Equal(t.Now().Round(24 * time.Hour)) {
		t.runTimeTodayInHours = 0
		fmt.Printf("Timer %s, date cycle\n", t.Name)
	} else {
		fmt.Printf("Timer %s, no date cycle\n", t.Name)
	}
}

func (t *TimedGPIO) Switch(value bool) {
	t.resetTotalIfDayChanged()
	now := t.Now()
	if !t.lastValue && value {
		fmt.Printf("Dummy timer %s starting\n", t.Name)
		t.lastStart = now
	}
	if t.lastValue && !value {
		t.runTimeTodayInHours += now.Sub(t.lastStart).Hours()
		fmt.Printf("Dummy timer %s stopped. total run today: %f\n", t.Name, t.runTimeTodayInHours)
	}
	t.lastValue = value
}
