package io

import "testing"

func TestMultiOut(t *testing.T) {
	dummy1 := DummyOutput{Name: "dummy1"}
	dummy2 := DummyOutput{Name: "dummy2"}
	multi := NewMultiOutput([]Output{&dummy1, &dummy2})
	multi.Switch(true)
	if !dummy1.Value {
		t.Fatal("Multiplexer set to true. Dummy1 shall be also true")
	}
	if !dummy2.Value {
		t.Fatal("Multiplexer set to true. Dummy2 shall be also true")
	}
	multi.Switch(false)
	if dummy1.Value {
		t.Fatal("Multiplexer set to false. Dummy1 shall be also false")
	}
	if dummy2.Value {
		t.Fatal("Multiplexer set to false. Dummy2 shall be also false")
	}
}

func TestOr(t *testing.T) {
	dummy := DummyOutput{Name: "dummy"}
	or := NewOrOutput(&dummy, 2)
	orMembers := or.GetMembers()
	if dummy.Value {
		t.Fatal("Default state for dummy output shall be false")
	}
	orMembers[0].Switch(true)
	if !dummy.Value {
		t.Fatal("One OR member is true. Dummy output shall be true")
	}
	orMembers[1].Switch(true)
	if !dummy.Value {
		t.Fatal("Both OR members are true. Dummy output shall be true")
	}
	orMembers[0].Switch(false)
	orMembers[1].Switch(false)
	if dummy.Value {
		t.Fatal("Both OR members are false. Dummy output shall be false")
	}
}
