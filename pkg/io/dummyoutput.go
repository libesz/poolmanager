package io

import "log"

type DummyOutput struct {
	Name_ string
	Value bool
}

func (d *DummyOutput) Name() string {
	return d.Name_
}

func (d *DummyOutput) Set(value bool) bool {
	if d.Value == value {
		log.Printf("Dummy %s unchanged: %t\n", d.Name_, value)
		return false
	}
	log.Printf("Dummy %s set to: %t\n", d.Name_, value)
	d.Value = value
	return true
}

func (d *DummyOutput) Get() bool {
	return d.Value
}
