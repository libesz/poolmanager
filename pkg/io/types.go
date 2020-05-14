package io

type Input interface {
	Name() string
	Type() string
	Degree() string
	Value() float64
}

type Output interface {
	Name() string
	Set(bool) bool
	Get() bool
}

type Haltable interface {
	Halt()
}
