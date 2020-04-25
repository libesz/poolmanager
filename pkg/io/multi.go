package io

type MultiOutput struct {
	realOutputs []Output
}

func NewMultiOutput(realOutputs []Output) MultiOutput {
	return MultiOutput{realOutputs: realOutputs}
}

func (m MultiOutput) Switch(value bool) {
	for _, output := range m.realOutputs {
		output.Switch(value)
	}
}

type OrOutputMember struct {
	id     int
	master OrOutput
}

type OrOutput struct {
	realOutput   Output
	memberStates map[int]bool
}

func (o *OrOutputMember) Switch(value bool) bool {
	return o.master.setMemberState(o.id, value)
}

func (o *OrOutput) setMemberState(i int, value bool) bool {
	o.memberStates[i] = value
	for _, memberState := range o.memberStates {
		if memberState {
			return o.realOutput.Switch(true)
		}
	}
	return o.realOutput.Switch(false)
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
