package controller

import (
	"testing"
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

type MockHeater struct {
	calledWith bool
	setReturns bool
	getReturns bool
}

func (m *MockHeater) Name() string {
	return "dummy"
}

func (m *MockHeater) Set(value bool) bool {
	m.calledWith = value
	return m.setReturns
}

func (m *MockHeater) Get() bool {
	return m.getReturns
}

type MockPump struct {
	calledWith bool
	setReturn  bool
	getReturn  bool
}

func (m *MockPump) Set(value bool) bool {
	m.calledWith = value
	return m.setReturn
}

func (m *MockPump) Get() bool {
	return m.getReturn
}

func (m *MockPump) Name() string {
	return "dummy"
}

type MockTempSensor struct {
	Temperature float64
}

func (s *MockTempSensor) Name() string {
	return "dummy"
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

func TestDefaultTempConfigPasses(t *testing.T) {
	heater := &MockHeater{}
	pumpOutput := &MockPump{}
	tempSensor := &MockTempSensor{}
	c := PoolTempController{
		heaterFactor: 0.5,
		heaterOutput: heater,
		pumpOutput:   pumpOutput,
		tempSensor:   tempSensor,
	}
	config := c.GetDefaultConfig()
	if err := c.ValidateConfig(config); err != nil {
		t.Fatal("Default config validation error", err.Error())
	}
}

func TestTempDisabled(t *testing.T) {
	heater := &MockHeater{}
	pumpOutput := &MockPump{}
	tempSensor := &MockTempSensor{}
	c := PoolTempController{
		heaterFactor: 0.5,
		heaterOutput: heater,
		pumpOutput:   pumpOutput,
		tempSensor:   tempSensor,
	}
	config := c.GetDefaultConfig()
	config.Toggles[configKeyEnabled] = false
	result := c.Act(config)
	if len(result) > 0 {
		t.Fatalf("Controller is disabled. Act shall not result in any new task, returned %d", len(result))
	}
}

func TestTempDisabledWhileActive(t *testing.T) {
	heater := &MockHeater{}
	pumpOutput := &MockPump{}
	tempSensor := &MockTempSensor{Temperature: 20}
	c := PoolTempController{
		heaterFactor: 0.5,
		heaterOutput: heater,
		pumpOutput:   pumpOutput,
		tempSensor:   tempSensor,
		now: func() time.Time {
			return time.Date(2020, 04, 15, 12, 0, 0, 0, time.Local)
		},
		pendingOperationReady: make(chan struct{}),
	}
	config := c.GetDefaultConfig()
	config.Toggles[configKeyEnabled] = true
	result := c.Act(config)
	if len(result) != 2 {
		t.Fatalf("Controller is enabled. Act shall result in two new tasks, returned %d", len(result))
	}
	result[0].Controller.Act(result[0].Config)

	heater.setReturns = true
	config.Toggles[configKeyEnabled] = false
	result = c.Act(config)
	if len(result) != 1 {
		t.Fatalf("Controller is disabled while heater is on. Act shall result in one new task, but returned %d", len(result))
	}
	if result[0].Controller.GetName() != "delayedOperation" {
		t.Fatalf("Controller is disabled while heater is on. Act shall result in new delayedOperation task, returned %s", result[0].Controller.GetName())
	}

	result = c.Act(config)
	if len(result) != 1 {
		t.Fatalf("Controller is disabled, but there is a delayedOperation in progress. Act shall result in one new task, but returned %d", len(result))
	}
	if result[0].Controller.GetName() != "Temperature controller" {
		t.Fatalf("Controller is disabled, but there is a delayedOperation in progress. Act shall result in new Temperature controller task, returned %s", result[0].Controller.GetName())
	}
}

func TestTemp(t *testing.T) {
	heater := &MockHeater{}
	pumpOutput := &MockPump{}
	tempSensor := &MockTempSensor{}
	c := PoolTempController{
		heaterFactor:          0.5,
		heaterOutput:          heater,
		pumpOutput:            pumpOutput,
		tempSensor:            tempSensor,
		pendingOperationReady: make(chan struct{}),
	}
	config := Config{Toggles: map[string]bool{configKeyEnabled: true}, Ranges: map[string]float64{configKeyTemp: 28.0, configKeyStart: 10, configKeyEnd: 13}}

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
	// Controller shall start the pump, and return a delayed operation to start the heater a bit later
	tempSensor.Temperature = 25
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 9, 0, 0, 0, time.Local)
	}
	pumpOutput.setReturn = true
	results := c.Act(config)
	if len(results) != 2 {
		t.Fatalf("Controller shall return two controllers to enqueue, returned %d", len(results))
	}
	if heater.calledWith != false {
		t.Fatal("Heater output shall not be started before the pump")
	}
	if pumpOutput.calledWith != true {
		t.Fatal("Pump output shall be started")
	}

	if results[0].Controller.GetName() != "delayedOperation" {
		t.Fatalf("Controller shall return a delayed operation as well, returned %s", results[1].Controller.GetName())
	}
	results[0].Controller.Act(results[0].Config)
	if heater.calledWith != true {
		t.Fatal("Heater output shall be started now by the delayed operation, as the pump already running")
	}

	// Expected time, temperature under the expected. Still heating.
	tempSensor.Temperature = 27.5
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 11, 0, 0, 0, time.Local)
	}
	results = c.Act(config)
	if len(results) != 2 {
		t.Fatalf("Controller shall return two controllers to enqueue, returned %d", len(results))
	}
	if pumpOutput.calledWith != true {
		t.Fatal("Pump output shall be started")
	}

	if results[0].Controller.GetName() != "delayedOperation" {
		t.Fatalf("Controller shall return a delayed operation as well, returned %s", results[1].Controller.GetName())
	}
	results[0].Controller.Act(results[0].Config)
	if heater.calledWith != true {
		t.Fatal("Heater output shall be started now by the delayed operation, as the pump already running")
	}

	// Expected time still, temperature above the expected. Stop heating and pump.
	tempSensor.Temperature = 28.5
	heater.setReturns = true
	pumpOutput.getReturn = true
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 12, 0, 0, 0, time.Local)
	}
	results = c.Act(config)
	if len(results) != 2 {
		t.Fatalf("Controller shall return two controllers to enqueue, returned %d", len(results))
	}
	if heater.calledWith != false {
		t.Fatal("Heater output shall be stopped")
	}
	if results[0].Controller.GetName() != "delayedOperation" {
		t.Fatalf("Controller shall return a delayed operation as well, returned %s", results[1].Controller.GetName())
	}
	results[0].Controller.Act(results[0].Config)
	if pumpOutput.calledWith != false {
		t.Fatal("Heater output shall be started now by the delayed operation, as the pump already running")
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

func TestTemperatureError(t *testing.T) {
	heater := &MockHeater{}
	pumpOutput := &MockPump{}
	tempSensor := &MockTempSensor{}
	c := PoolTempController{
		heaterFactor:          0.5,
		heaterOutput:          heater,
		pumpOutput:            pumpOutput,
		tempSensor:            tempSensor,
		pendingOperationReady: make(chan struct{}),
	}
	config := Config{Toggles: map[string]bool{configKeyEnabled: true}, Ranges: map[string]float64{configKeyTemp: 28.0, configKeyStart: 10, configKeyEnd: 13}}

	tempSensor.Temperature = io.InputError
	c.now = func() time.Time {
		return time.Date(2020, 04, 15, 2, 0, 0, 0, time.Local)
	}
	results := c.Act(config)
	if heater.calledWith != false {
		t.Fatal("Heater output shall be stopped")
	}
	if len(results) != 1 {
		t.Fatalf("Controller shall return one controller to enqueue, returned %d", len(results))
	}
	if results[0].Controller.GetName() != "Temperature controller" {
		t.Fatalf("Controller shall return itself, returned %s", results[1].Controller.GetName())
	}

}
