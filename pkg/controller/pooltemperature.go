package controller

import (
	"fmt"
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

func NewPoolTempController(
	heaterFactor float64,
	tempSensor io.Input,
	heaterOutput io.Output,
	now func() time.Time) PoolTempController {
	return PoolTempController{
		heaterFactor: heaterFactor,
		tempSensor:   tempSensor,
		heaterOutput: heaterOutput,
		now:          now,
	}
}

type PoolTempController struct {
	heaterFactor float64
	tempSensor   io.Input
	heaterOutput io.Output
	now          func() time.Time
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

func (c PoolTempController) Act(config Config) time.Duration {
	desiredTemp := config[configKeyTemp]
	currentTemp := c.tempSensor.Value()
	now := c.now()
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), int(config[configKeyStart]), 0, 0, 0, now.Local().Location())
	nextStop := time.Date(now.Year(), now.Month(), now.Day(), int(config[configKeyEnd]), 0, 0, 0, now.Local().Location())

	if now.After(nextStart) {
		if now.Before(nextStop) {
			if desiredTemp >= currentTemp {
				fmt.Printf("In the active period, the temparature is %f, need more heat to reach %f\n", currentTemp, desiredTemp)
				c.heaterOutput.Switch(true)
				return 5 * time.Second
			}
			fmt.Printf("In the active period, the temperature is %f, already fine\n", currentTemp)
			c.heaterOutput.Switch(false)
			return 5 * time.Second
		}
		nextStart = nextStart.Add(24 * time.Hour)
	}
	thisManyHoursUntilNextStart := nextStart.Sub(now).Hours()
	calculatedDesiredTemp := desiredTemp - thisManyHoursUntilNextStart*c.heaterFactor
	if calculatedDesiredTemp >= currentTemp {
		fmt.Printf("Not in the active period. Hours until the next one: %f. Calculated desired temperature: %f, need more heat\n", thisManyHoursUntilNextStart, calculatedDesiredTemp)
		c.heaterOutput.Switch(true)
		return 5 * time.Second
	}
	fmt.Printf("The temperature is already fine\n")
	c.heaterOutput.Switch(false)
	return 5 * time.Second
}

func (c *PoolTempController) GetName() string {
	return "PoolTempController"
}
