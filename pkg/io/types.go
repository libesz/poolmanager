package io

type Sensor interface {
	Type() string
	Degree() string
	Value() float64
}

type Actuator interface {
	Switch(bool)
}

type SensingActuator interface {
	Sensor
	Actuator
}
