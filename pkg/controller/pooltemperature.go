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

type PendingHeaterOperation int

const (
	None PendingHeaterOperation = 0
	On   PendingHeaterOperation = 1
	Off  PendingHeaterOperation = 2
)

type PoolTempController struct {
	heaterFactor           float64
	tempSensor             io.Input
	heaterOutput           io.Output
	pumpOutput             io.Output
	pendingHeaterOperation PendingHeaterOperation
	now                    func() time.Time
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

func (c *PoolTempController) Act(config Config) time.Duration {
	if c.pendingHeaterOperation != None {
		if c.pendingHeaterOperation == On {
			c.heaterOutput.Switch(true)
		} else if c.pendingHeaterOperation == Off {
			c.heaterOutput.Switch(false)
		}
		c.pendingHeaterOperation = None
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
			c.pendingHeaterOperation = On
			return 5 * time.Second
		}
		c.heaterOutput.Switch(true)
		return 5 * time.Second
	}
	log.Printf("The temperature is already fine\n")
	if c.pumpOutput.Switch(false) {
		c.pendingHeaterOperation = Off
		return 5 * time.Second
	}
	c.heaterOutput.Switch(false)
	return 5 * time.Second
}

func (c *PoolTempController) GetName() string {
	return "PoolTempController"
}
