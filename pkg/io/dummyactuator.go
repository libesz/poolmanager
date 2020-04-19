package io

import "fmt"

type DummyActuator struct {
	Name string
}

func (a *DummyActuator) Switch(value bool) {
	fmt.Printf("Dummy %s set to: %t\n", a.Name, value)
}
