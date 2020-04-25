package io

import "fmt"

type DummyOutput struct {
	Name  string
	Value bool
}

func (a *DummyOutput) Switch(value bool) {
	a.Value = value
	fmt.Printf("Dummy %s set to: %t\n", a.Name, value)
}
