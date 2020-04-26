package controller

import (
	"testing"
	"time"
)

type MockHeater struct {
	calledWith    bool
	switchReturns bool
}

func (m *MockHeater) Switch(value bool) bool {
	m.calledWith = value
	return m.switchReturns
}

type MockPump struct {
	calledWith    bool
	switchReturns bool
}

func (m *MockPump) Switch(value bool) bool {
	m.calledWith = value
	return m.switchReturns

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
	pumpOutput := &MockPump{}
	tempSensor := &MockTempSensor{}
	c := PoolTempController{
		heaterFactor: 0.5,
		heaterOutput: heater,
		pumpOutput:   pumpOutput,
		tempSensor:   tempSensor,
	}
	config := Config{"desired temperature": 28.0, "start hour": 10, "end hour": 13}

	// Early morning, no heat
	tempSensor.Temperature = 25
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 2, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != false {
		t.Fatal("Heater output shall be stopped")
	}

	// Almost the expected time, according to the heater factor, it shall already work.
	// For first invocation it shall start the pump, for second it shall start the heater
	tempSensor.Temperature = 25
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 9, 0, 0, 0, time.Local)
	}
	pumpOutput.switchReturns = true
	results := c.Act(config)
	if len(results) != 2 {
		t.Fatal("Controller shall return two controllers to enqueue")
	}
	if heater.calledWith != false {
		t.Fatal("Heater output shall not be started before the pump")
	}
	if pumpOutput.calledWith != true {
		t.Fatal("Pump output shall be started")
	}
	c.pendingHeaterOperation = false

	pumpOutput.switchReturns = false
	c.Act(config)
	if heater.calledWith != true {
		t.Fatal("Heater output shall be started now, as the pump already running")
	}
	if pumpOutput.calledWith != true {
		t.Fatal("Pump output shall be also started")
	}

	// Expected time, temperature under the expected. Still heating.
	tempSensor.Temperature = 27.5
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 11, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != true {
		t.Fatal("Heater output shall be started")
	}

	// Expected time still, temperature above the expected. Stop heating.
	tempSensor.Temperature = 28.5
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 12, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != false {
		t.Fatal("Heater output shall be stopped")
	}

	// After expected time, expected temperature is calculated for tomorrow's start time. No heating.
	tempSensor.Temperature = 25
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 18, 0, 0, 0, time.Local)
	}
	c.Act(config)
	if heater.calledWith != false {
		t.Fatal("Heater output shall be stopped")
	}
}
