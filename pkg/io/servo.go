package io

import "log"

type Servo struct {
	pin, offDutyPercentage, onDutyPercentage int
}

func NewServo(pin, offDutyPercentage, onDutyPercentage int) *Servo {
	return &Servo{pin: pin, offDutyPercentage: offDutyPercentage, onDutyPercentage: onDutyPercentage}
}

func (s *Servo) Name() string {
	return "Servo"
}

func (s *Servo) Set(bool) bool {
	return false
}

func (s *Servo) Get() bool {
	return false
}

func (s *Servo) Halt() {
	log.Println("Servo teardown")
}
