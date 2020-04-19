package controller

import (
	"fmt"
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

type PoolTempController struct {
	HeaterFactor  float64
	TempSensor    io.Sensor
	HeaterOutputs []io.Actuator
	Now           func() time.Time
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

func (c PoolTempController) SetHeater(value bool) {
	for _, item := range c.HeaterOutputs {
		item.Switch(value)
	}
}

func (c PoolTempController) Act(config Config) {
	desiredTemp := config[configKeyTemp]
	currentTemp := c.TempSensor.Value()
	now := c.Now()
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), int(config[configKeyStart]), 0, 0, 0, now.Local().Location())
	nextStop := time.Date(now.Year(), now.Month(), now.Day(), int(config[configKeyEnd]), 0, 0, 0, now.Local().Location())

	if now.After(nextStart) {
		if now.Before(nextStop) {
			if desiredTemp >= currentTemp {
				fmt.Printf("We are actually in the active period, and need more heat\n")
				c.SetHeater(true)
				return
			}
			fmt.Printf("We are actually in the active period, the temperature is already fine\n")
			c.SetHeater(false)
			return
		}
		nextStart = nextStart.Add(24 * time.Hour)
	}
	thisManyHoursUntilNextStart := now.Sub(nextStart).Hours()
	calculatedDesiredTemp := desiredTemp + thisManyHoursUntilNextStart/c.HeaterFactor
	fmt.Printf("We are not in the active period. Hours until the next one: %f. Calculated desired temperature: %f\n", thisManyHoursUntilNextStart, calculatedDesiredTemp)
	if calculatedDesiredTemp >= currentTemp {
		fmt.Printf("Need nore heat\n")
		c.SetHeater(true)
		return
	}
	fmt.Printf("The temperature is already fine\n")
	c.SetHeater(false)
	return
}
