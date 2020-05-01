package controller

import (
	"fmt"
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
		heaterFactor:          heaterFactor,
		tempSensor:            tempSensor,
		heaterOutput:          heaterOutput,
		pumpOutput:            pumpOutput,
		now:                   now,
		pendingOperationReady: make(chan struct{}, 1),
	}
}

type PoolTempController struct {
	heaterFactor          float64
	tempSensor            io.Input
	heaterOutput          io.Output
	pumpOutput            io.Output
	pendingOperation      bool
	pendingOperationReady chan struct{}
	now                   func() time.Time
}

const (
	configKeyTemp  = "desired temperature"
	configKeyStart = "start hour"
	configKeyEnd   = "end hour"
)

func (c PoolTempController) GetConfig() ConfigProperties {
	return ConfigProperties{
		configKeyTemp: ConfigRange{
			Default: 25.0,
			Min:     20.0,
			Max:     30.0,
			Step:    0.5,
		},
		configKeyStart: ConfigRange{
			Default: 13,
			Min:     0,
			Max:     23,
			Step:    1,
		},
		configKeyEnd: ConfigRange{
			Default: 16,
			Min:     0,
			Max:     23,
			Step:    1,
		},
	}
}

func (c PoolTempController) ValidateConfig(config Config) error {
	temp, ok := config[configKeyTemp].(float64)
	if !ok {
		return fmt.Errorf("Temperature is not a float")
	}
	if temp < 20.0 || temp > 30.0 {
		return fmt.Errorf("Temperature is outside of the allowed range")
	}
	if int(temp*10)%5 != 0 {
		return fmt.Errorf("Temperature is not a valid step (.0 or .5 required)")
	}
	start, ok := config[configKeyStart].(int)
	if !ok {
		return fmt.Errorf("Start time it not an int")
	}
	if start < 0 || start > 23 {
		return fmt.Errorf("Start time is outside of the allowed range")
	}
	end, ok := config[configKeyEnd].(int)
	if !ok {
		return fmt.Errorf("Start time it not an int")
	}
	if end < 0 || end > 23 {
		return fmt.Errorf("Start time is outside of the allowed range")
	}
	if start > end {
		return fmt.Errorf("Start time set to earlier than end time")
	}
	return nil
}

type delayedOperation struct {
	setTo                 bool
	output                io.Output
	pendingOperationReady chan struct{}
}

func (c *delayedOperation) Act(config Config) []EnqueueRequest {
	log.Printf("delayedOperation set heater to: %t\n", c.setTo)
	c.output.Set(c.setTo)
	close(c.pendingOperationReady)
	return nil
}

func (c *delayedOperation) GetName() string {
	return "delayedOperation"
}

func (c delayedOperation) GetConfig() ConfigProperties {
	return ConfigProperties{}
}

func (c delayedOperation) ValidateConfig(Config) error {
	return nil
}

func (c *PoolTempController) Act(config Config) []EnqueueRequest {
	if c.pendingOperation {
		select {
		case <-c.pendingOperationReady:
			c.pendingOperation = false
			c.pendingOperationReady = make(chan struct{}, 1)
		default:
			log.Printf("Pending heater operation is ongoing.")
			return []EnqueueRequest{{Controller: c, After: 5 * time.Second}}
		}
	}
	desiredTemp := config[configKeyTemp].(float64)
	currentTemp := c.tempSensor.Value()
	now := c.now()
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), config[configKeyStart].(int), 0, 0, 0, now.Local().Location())
	nextStop := time.Date(now.Year(), now.Month(), now.Day(), config[configKeyEnd].(int), 0, 0, 0, now.Local().Location())

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
		if c.pumpOutput.Set(true) {
			c.pendingOperation = true
			pending := delayedOperation{setTo: true, output: c.heaterOutput, pendingOperationReady: c.pendingOperationReady}
			return []EnqueueRequest{{Controller: c, After: 5 * time.Second}, {Controller: &pending, After: 6 * time.Second}}
		}
		c.heaterOutput.Set(true)
		return []EnqueueRequest{{Controller: c, After: 5 * time.Second}}
	}
	log.Printf("The temperature is already fine\n")
	if c.heaterOutput.Set(false) {
		c.pendingOperation = true
		pending := delayedOperation{setTo: false, output: c.pumpOutput, pendingOperationReady: c.pendingOperationReady}
		return []EnqueueRequest{{Controller: c, After: 5 * time.Second}, {Controller: &pending, After: 6 * time.Second}}
	}
	c.pumpOutput.Set(false)
	return []EnqueueRequest{{Controller: c, After: 5 * time.Second}}
}

func (c *PoolTempController) GetName() string {
	return "PoolTempController"
}
