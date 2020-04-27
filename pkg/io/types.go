package io

type Input interface {
	Type() string
	Degree() string
	Value() float64
}

type Output interface {
	Set(bool) bool
	Get() bool
}
