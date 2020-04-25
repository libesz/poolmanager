package controller

import (
	"testing"
)

type MockTimer struct {
	nextValue  float64
	calledWith bool
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

func (t *MockTimer) Switch(value bool) bool {
	t.calledWith = value
	return false
}

func TestNormal(t *testing.T) {
	timer := &MockTimer{}
	c := NewPoolPumpController(timer)
	config := Config{"desired runtime per day": 1}
	c.Act(config)
	if timer.calledWith != true {
		t.Error("Timer output shall be started")
	}
	timer.nextValue = 3
	c.Act(config)
	if timer.calledWith != false {
		t.Error("Timer output shall be stopped")
	}
}
