package controller

import (
	"testing"
	"time"
)

type MockHeater struct {
	calledWith bool
}

func (t *MockHeater) Switch(value bool) {
	t.calledWith = value
}

type MockTempSensor struct {
	Temperature float64
}

func (s *MockTempSensor) Type() string {
	return "Temperature"
}

func (s *MockTempSensor) Degree() string {
	return "°C"
}

func (s *MockTempSensor) Value() float64 {
	return s.Temperature
}

func TestTemp(t *testing.T) {
	heater := &MockHeater{}
	tempSensor := &MockTempSensor{}
	c := PoolTempController{
		HeaterFactor: 0.5,
		HeaterOutput: heater,
		TempSensor:   tempSensor,
	}
	config := Config{"desired temperature": 28, "start hour": 10, "end hour": 13}

	// Early morning, no heat
	tempSensor.Temperature = 20
	c.Now = func() time.Time {
		return time.Date(2020, 04, 15, 5, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != false {
		t.Error("Heater output shall be stopped")
	}

	// Almost the expected time, according to the heater factor, it shall already work
	tempSensor.Temperature = 25
	c.Now = func() time.Time {
		return time.Date(2020, 04, 15, 9, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != true {
		t.Error("Heater output shall be started")
	}

	// Expected time, temperature under the expected. Still heating.
	tempSensor.Temperature = 27.5
	c.Now = func() time.Time {
		return time.Date(2020, 04, 15, 11, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != true {
		t.Error("Heater output shall be started")
	}

	// Expected time still, temperature above the expected. Stop heating.
	tempSensor.Temperature = 28.5
	c.Now = func() time.Time {
		return time.Date(2020, 04, 15, 12, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != false {
		t.Error("Heater output shall be stopped")
	}

	// After expected time, expected temperature is calculated for tomorrow's start time. No heating.
	tempSensor.Temperature = 25
	c.Now = func() time.Time {
		return time.Date(2020, 04, 15, 18, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != false {
		t.Error("Heater output shall be stopped")
	}
}
