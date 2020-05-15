package io

import "testing"

func TestDistributor(t *testing.T) {
	dummy1 := DummyOutput{Name_: "dummy1"}
	dummy2 := DummyOutput{Name_: "dummy2"}
	multi := NewOutputDistributor("dummy", []Output{&dummy1, &dummy2})
	setResult := multi.Set(true)
	if !dummy1.Value {
		t.Fatal("Distributor set to true. Dummy1 shall be also true")
	}
	if !dummy2.Value {
		t.Fatal("Distributor set to true. Dummy2 shall be also true")
	}
	if !setResult {
		t.Fatal("Set result shall indicate the operation changed the output or not. Now it changed")
	}
	if !multi.Get() {
		t.Fatal("Distributor shall be on now")
	}

	setResult = multi.Set(false)
	if dummy1.Value {
		t.Fatal("Distributor set to false. Dummy1 shall be also false")
	}
	if dummy2.Value {
		t.Fatal("Distributor set to false. Dummy2 shall be also false")
	}
	if !setResult {
		t.Fatal("Set result shall indicate the operation changed the output or not. Now it changed")
	}
	if multi.Get() {
		t.Fatal("Distributor shall be on now")
	}
}

func TestOr(t *testing.T) {
	dummy := DummyOutput{Name_: "dummy"}
	orMembers := NewOrOutput("dummy or", &dummy, 2)
	if dummy.Value {
		t.Fatal("Default state for dummy output shall be false")
	}
	orMembers[0].Set(true)
	if !dummy.Value {
		t.Fatal("One OR member is true. Dummy output shall be true")
	}
	if !(orMembers[0].Get() == orMembers[1].Get() == true) {
		t.Fatal("Both OR members shall get the value of the final output, which now should be true")
	}
	orMembers[1].Set(true)
	if !dummy.Value {
		t.Fatal("Both OR members are true. Dummy output shall be true")
	}
	orMembers[0].Set(false)
	orMembers[1].Set(false)
	if dummy.Value {
		t.Fatal("Both OR members are false. Dummy output shall be false")
	}
}
