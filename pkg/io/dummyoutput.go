package io

import "fmt"

type DummyOutput struct {
	Name string
}

func (a *DummyOutput) Switch(value bool) {
	fmt.Printf("Dummy %s set to: %t\n", a.Name, value)
}
