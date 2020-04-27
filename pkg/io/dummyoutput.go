package io

import "log"

type DummyOutput struct {
	Name  string
	Value bool
}

func (a *DummyOutput) Set(value bool) bool {
	if a.Value == value {
		log.Printf("Dummy %s unchanged: %t\n", a.Name, value)
		return false
	}
	log.Printf("Dummy %s set to: %t\n", a.Name, value)
	a.Value = value
	return true
}
