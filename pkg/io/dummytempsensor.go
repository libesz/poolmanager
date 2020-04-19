package io

type DummyTempSensor struct {
	Temperature float64
}

func (s *DummyTempSensor) Type() string {
	return "Temperature"
}

func (s *DummyTempSensor) Degree() string {
	return "°C"
}

func (s *DummyTempSensor) Value() float64 {
	return s.Temperature
}
