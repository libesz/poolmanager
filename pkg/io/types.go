package io

type Input interface {
	Type() string
	Degree() string
	Value() float64
}

type Output interface {
	Switch(bool)
}

type InputOutput interface {
	Input
	Output
}
