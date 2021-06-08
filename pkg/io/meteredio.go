package io

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type meteredInput struct {
	input Input
	gauge prometheus.Gauge
}

func NewMeteredInput(input Input) *meteredInput {
	sensorData := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "poolmanager",
		Subsystem: "input",
		Name:      strings.ReplaceAll(input.Name(), " ", ""),
		Help:      input.Type(),
	})
	prometheus.MustRegister(sensorData)
	return &meteredInput{
		input: input,
		gauge: sensorData,
	}
}

func (m *meteredInput) Name() string {
	return m.input.Name()
}

func (m *meteredInput) Degree() string {
	return m.input.Degree()
}

func (m *meteredInput) Value() float64 {
	value := m.input.Value()
	m.gauge.Set(value)
	return value
}

func (m *meteredInput) Type() string {
	return m.input.Type()
}

type meteredOutput struct {
	output Output
	gauge  prometheus.Gauge
}

func NewMeteredOutput(output Output) *meteredOutput {
	outputData := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "poolmanager",
		Subsystem: "output",
		Name:      strings.ReplaceAll(output.Name(), " ", ""),
		Help:      "bool",
	})
	prometheus.MustRegister(outputData)
	return &meteredOutput{
		output: output,
		gauge:  outputData,
	}
}

func (m *meteredOutput) Name() string {
	return m.output.Name()
}

func (m *meteredOutput) Get() bool {

	return m.output.Get()
}

func (m *meteredOutput) Set(value bool) bool {
	if value {
		m.gauge.Set(1)
	} else {
		m.gauge.Set(0)
	}
	return m.output.Set(value)
}
