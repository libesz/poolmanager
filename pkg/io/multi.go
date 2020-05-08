package io

import "strconv"

type MultiOutput struct {
	name        string
	realOutputs []Output
}

func NewMultiOutput(name string, realOutputs []Output) MultiOutput {
	return MultiOutput{realOutputs: realOutputs}
}

func (m *MultiOutput) Name() string {
	return m.name
}

func (m MultiOutput) Set(value bool) {
	for _, output := range m.realOutputs {
		output.Set(value)
	}
}

type OrOutputMember struct {
	id     int
	master OrOutput
}

type OrOutput struct {
	name         string
	realOutput   Output
	memberStates map[int]bool
}

func (o *OrOutputMember) Set(value bool) bool {
	return o.master.setMemberState(o.id, value)
}

func (o *OrOutputMember) Get() bool {
	return o.master.realOutput.Get()
}

func (o *OrOutputMember) Name() string {
	return o.master.name + " member " + strconv.Itoa(o.id)
}

func (o *OrOutput) setMemberState(i int, value bool) bool {
	o.memberStates[i] = value
	for _, memberState := range o.memberStates {
		if memberState {
			return o.realOutput.Set(true)
		}
	}
	return o.realOutput.Set(false)
}

func NewOrOutput(realOutput Output, amount int) []OrOutputMember {
	result := OrOutput{realOutput: realOutput}
	result.memberStates = make(map[int]bool)
	members := []OrOutputMember{}
	for i := 0; i < amount; i++ {
		members = append(members, OrOutputMember{id: i, master: result})
	}
	return members
}
