package io

import (
	"log"

	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/host"
)

type Servo struct {
	name                                string
	pin                                 gpio.PinIO
	offDutyPercentage, onDutyPercentage gpio.Duty
	state                               bool
}

func NewServo(name, pin, offDutyPercentage, onDutyPercentage string) *Servo {
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	p := gpioreg.ByName(pin)
	if p == nil {
		log.Fatalf("Servo: Failed to find %s\n", pin)
	}
	result := &Servo{name: name, pin: p}
	duty, err := gpio.ParseDuty(offDutyPercentage)
	if err != nil {
		log.Fatalf("Servo: Failed to convert offDutyPercentage: %s\n", offDutyPercentage)
	}
	result.offDutyPercentage = duty

	duty, err = gpio.ParseDuty(onDutyPercentage)
	if err != nil {
		log.Fatalf("Servo: Failed to convert onDutyPercentage: %s\n", onDutyPercentage)
	}
	result.onDutyPercentage = duty

	return result
}

func (s *Servo) Name() string {
	return s.name
}

func (s *Servo) Set(newState bool) bool {
	if newState == s.state {
		return false
	}
	s.state = newState
	var duty gpio.Duty
	if s.state {
		log.Println("Servo: set to ON")
		duty = s.onDutyPercentage
	} else {
		log.Println("Servo: set to OFF")
		duty = s.offDutyPercentage
	}
	if err := s.pin.PWM(duty, 200*physic.Hertz); err != nil {
		log.Printf("Servo: error setting new duty: %s", err.Error())
	}
	return true
}

func (s *Servo) Get() bool {
	return s.state
}

func (s *Servo) Halt() {
	log.Println("Servo: teardown")
	s.Set(false)
	if err := s.pin.Halt(); err != nil {
		log.Printf("Servo: teardown error: %s\n", err.Error())
	}
}
