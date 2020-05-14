package controller

import (
	"testing"
	"time"
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

func (t *MockTimer) Name() string {
	return "dummy"
}

func TestPumpDefaultConfig(t *testing.T) {
	timer := &MockTimer{}
	c := NewPoolPumpController(timer, timer, tenpm)
	config := c.GetDefaultConfig()
	if err := c.ValidateConfig(config); err != nil {
		t.Fatal("Default config validation error", err.Error())
	}
}

func tenpm() time.Time {
	return time.Date(2020, 04, 15, 22, 0, 0, 0, time.Local)
}

func eightpm() time.Time {
	return time.Date(2020, 04, 15, 20, 0, 0, 0, time.Local)
}

func twelvepm() time.Time {
	return time.Date(2020, 04, 16, 0, 0, 0, 0, time.Local)
}

func TestNormal(t *testing.T) {
	timer := &MockTimer{}
	c := NewPoolPumpController(timer, timer, eightpm)
	config := Config{Toggles: map[string]bool{configKeyEnabled: true}, Ranges: map[string]float64{configKeyRuntime: 2}}
	result := c.Act(config)
	if timer.setCalledWith != false {
		t.Fatal("Timer output not shall be started")
	}
	if len(result) != 1 {
		t.Fatal("Result must contain one enqueue request")
	}
	if result[0].After.Hours() != 2.0 {
		t.Fatalf("Result must return two hours duration, returned: %s\n", result[0].After.String())
	}

	c.now = tenpm
	result = c.Act(config)
	if timer.setCalledWith != true {
		t.Fatal("Timer output shall be started")
	}
	if len(result) != 1 {
		t.Fatal("Result must contain one enqueue request")
	}
	if result[0].After.Hours() != 2.0 {
		t.Fatalf("Result must return 2 hours duration, returned: %s\n", result[0].After.String())
	}

	timer.nextValue = 0
	c.now = twelvepm
	result = c.Act(config)
	if timer.setCalledWith != false {
		t.Fatal("Timer output shall be stopped")
	}
	if len(result) != 1 {
		t.Fatal("Result must contain one enqueue request")
	}
	if result[0].After.Hours() != 22.0 {
		t.Fatalf("Result must return 22 hours duration, returned: %s\n", result[0].After.String())
	}
}

func TestAlreadyOverThePlannedDuration(t *testing.T) {
	timer := &MockTimer{}
	timer.nextValue = 3
	c := NewPoolPumpController(timer, timer, eightpm)
	config := Config{Toggles: map[string]bool{configKeyEnabled: true}, Ranges: map[string]float64{configKeyRuntime: 2}}
	result := c.Act(config)
	if timer.setCalledWith != false {
		t.Fatal("Timer output not shall be started")
	}
	if len(result) != 1 {
		t.Fatal("Result must contain one enqueue request")
	}
	if result[0].After.Hours() != 26.0 {
		t.Fatalf("Result must return 26 hours duration, returned: %s\n", result[0].After.String())
	}

	c.now = tenpm
	result = c.Act(config)
	if timer.setCalledWith != false {
		t.Fatal("Timer output not shall be started")
	}
	if len(result) != 1 {
		t.Fatal("Result must contain one enqueue request")
	}
	if result[0].After.Hours() != 24.0 {
		t.Fatalf("Result must return 24 hours duration, returned: %s\n", result[0].After.String())
	}
}
