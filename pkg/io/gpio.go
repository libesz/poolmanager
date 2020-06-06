package io

import (
	"log"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

type GPIOOutput struct {
	name      string
	pin       gpio.PinIO
	activeLow bool
	state     bool
}

func NewGPIOOutput(name, pin string, activeLow bool) *GPIOOutput {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	p := gpioreg.ByName(pin)
	if p == nil {
		log.Fatalf("GPIO: Failed to find %s\n", pin)
	}
	result := &GPIOOutput{name: name, pin: p, activeLow: activeLow}
	result.state = true // force initial state change
	result.Set(false)
	return result
}

func (g *GPIOOutput) Name() string {
	return g.name
}

func (g *GPIOOutput) Set(newState bool) bool {
	if g.state == newState {
		//log.Printf("GPIO: %s unchanged: %t\n", g.name, newState)
		return false
	}
	log.Printf("GPIO: %s set to: %t\n", g.name, newState)
	g.state = newState
	if (g.state && !g.activeLow) || (!g.state && g.activeLow) {
		_ = g.pin.Out(gpio.High)
	} else {
		_ = g.pin.Out(gpio.Low)
	}
	return true
}

func (d *GPIOOutput) Get() bool {
	return d.state
}

func (d *GPIOOutput) Halt() {
	log.Printf("GPIO: %s teardown", d.name)
	d.Set(false)
}
