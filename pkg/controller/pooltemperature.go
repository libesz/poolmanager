package controller

import (
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

func NewPoolTempController(
	heaterFactor float64,
	tempSensor io.Input,
	heaterOutput io.Output,
	pumpOutput io.Output,
	now func() time.Time) PoolTempController {
	return PoolTempController{
		heaterFactor: heaterFactor,
		tempSensor:   tempSensor,
		heaterOutput: heaterOutput,
		pumpOutput:   pumpOutput,
		now:          now,
	}
}

type PoolTempController struct {
	heaterFactor                float64
	tempSensor                  io.Input
	heaterOutput                io.Output
	pumpOutput                  io.Output
	pendingHeaterOperation      bool
	pendingHeaterOperationReady chan struct{}
	now                         func() time.Time
}

const (
	configKeyTemp  = "desired temperature"
	configKeyStart = "start hour"
	configKeyEnd   = "end hour"
)

func (c PoolTempController) GetConfigKeys() []string {
	return []string{
		configKeyTemp,
		configKeyStart,
		configKeyEnd,
	}
}

type delayedHeaterOperation struct {
	setTo  bool
	heater io.Output
}

func (c *delayedHeaterOperation) Act(config Config) []EnqueueRequest {
	log.Printf("DelayedHeaterOperation set heater to: %t\n", c.setTo)
	c.heater.Switch(c.setTo)
	return nil
}

func (c *delayedHeaterOperation) GetName() string {
	return "DelayedHeaterOperation"
}

func (c delayedHeaterOperation) GetConfigKeys() []string {
	return []string{}
}

func (c *PoolTempController) Act(config Config) []EnqueueRequest {
	if c.pendingHeaterOperation {
		select {
		case <-c.pendingHeaterOperationReady:
			c.pendingHeaterOperation = false
		default:
		}
	}
	desiredTemp := config[configKeyTemp]
	currentTemp := c.tempSensor.Value()
	now := c.now()
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), int(config[configKeyStart]), 0, 0, 0, now.Local().Location())
	nextStop := time.Date(now.Year(), now.Month(), now.Day(), int(config[configKeyEnd]), 0, 0, 0, now.Local().Location())

	var thisManyHoursUntilNextStart float64
	if now.After(nextStart) {
		if now.Before(nextStop) {
			nextStart = now
		} else {
			nextStart = nextStart.Add(24 * time.Hour)
		}
	}
	thisManyHoursUntilNextStart = nextStart.Sub(now).Hours()
	calculatedDesiredTemp := desiredTemp - thisManyHoursUntilNextStart*c.heaterFactor
	if calculatedDesiredTemp >= currentTemp {
		log.Printf("Hours until the next active period: %f. Calculated desired temperature: %f, need more heat\n", thisManyHoursUntilNextStart, calculatedDesiredTemp)
		if c.pumpOutput.Switch(true) {
			pending := delayedHeaterOperation{setTo: true, heater: c.heaterOutput}
			return []EnqueueRequest{{Controller: c, After: 5 * time.Second}, {Controller: &pending, After: 5 * time.Second}}
		}
		c.heaterOutput.Switch(true)
		return []EnqueueRequest{{Controller: c, After: 5 * time.Second}}
	}
	log.Printf("The temperature is already fine\n")
	if c.pumpOutput.Switch(false) {
		pending := delayedHeaterOperation{setTo: false, heater: c.heaterOutput}
		return []EnqueueRequest{{Controller: c, After: 5 * time.Second}, {Controller: &pending, After: 5 * time.Second}}
	}
	c.heaterOutput.Switch(false)
	return []EnqueueRequest{{Controller: c, After: 5 * time.Second}}
}

func (c *PoolTempController) GetName() string {
	return "PoolTempController"
}
