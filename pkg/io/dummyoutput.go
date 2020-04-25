package io

import "log"

type DummyOutput struct {
	Name  string
	Value bool
}

func (a *DummyOutput) Switch(value bool) {
	a.Value = value
	log.Printf("Dummy %s set to: %t\n", a.Name, value)
}
