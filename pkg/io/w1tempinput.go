package io

import (
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

type OneWireTemperatureInput struct {
	name, fullPath string
}

func NewOneWireTemperatureInput(name, deviceId string) *OneWireTemperatureInput {
	fullPath := "/sys/bus/w1/devices/28-" + deviceId + "/w1_slave"
	return &OneWireTemperatureInput{name: name, fullPath: fullPath}
}

func (o *OneWireTemperatureInput) Type() string {
	return "Temperature"
}

func (o *OneWireTemperatureInput) Name() string {
	return o.name
}

func (o *OneWireTemperatureInput) Degree() string {
	return "°C"
}

func (o *OneWireTemperatureInput) Value() float64 {
	raw, err := ioutil.ReadFile(o.fullPath)
	if err != nil {
		log.Printf("OneWireTemperatureInput: Failed to open onewire file: %s, error: %s", o.fullPath, err.Error())
		return InputError
	}
	rawMatched := regexp.MustCompile(` t=([0-9]{4,5})`).FindSubmatch(raw)
	if rawMatched == nil || len(rawMatched) != 2 {
		log.Printf("OneWireTemperatureInput: Failed to parse valid value from onewire file: %s", o.fullPath)
		return InputError
	}
	converted, err := strconv.Atoi(string(rawMatched[1]))
	if err != nil {
		log.Printf("OneWireTemperatureInput: Failed to convert string %s to int from onewire file: %s, error: %s", rawMatched[1], o.fullPath, err.Error())
		return InputError
	}
	final := float64(converted) / 1000

	return final
}
