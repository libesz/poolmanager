package controller

import (
	"testing"
)

type MockTimer struct {
	nextValue     float64
	setCalledWith bool
	setReturns    bool
	getReturns    bool
}

func (t *MockTimer) Type() string {
	return "time"
}

func (t *MockTimer) Degree() string {
	return "H"
}

func (t *MockTimer) Value() float64 {
	return t.nextValue
}

func (t *MockTimer) Set(value bool) bool {
	t.setCalledWith = value
	return t.setReturns
}

func (t *MockTimer) Get() bool {
	return t.getReturns
}

func TestNormal(t *testing.T) {
	timer := &MockTimer{}
	c := NewPoolPumpController(timer, timer)
	config := Config{Ranges: map[string]float64{"desired runtime per day": 1}}
	c.Act(config)
	if timer.setCalledWith != true {
		t.Error("Timer output shall be started")
	}
	timer.nextValue = 3
	c.Act(config)
	if timer.setCalledWith != false {
		t.Error("Timer output shall be stopped")
	}
}
