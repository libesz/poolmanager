package io

type MultiOutput struct {
	realOutputs []Output
}

func NewMultiOutput(realOutputs []Output) MultiOutput {
	return MultiOutput{realOutputs: realOutputs}
}

func (m *MultiOutput) Switch(value bool) {
	for _, output := range m.realOutputs {
		output.Switch(value)
	}
}
