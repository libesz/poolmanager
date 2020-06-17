package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

type PoolPumpController struct {
	timer      io.Input
	pumpOutput io.Output
	now        func() time.Time
}

const configKeyRuntime = "Desired runtime per day in hours"

func NewPoolPumpController(timer io.Input, pumpOutput io.Output, now func() time.Time) PoolPumpController {
	return PoolPumpController{timer: timer, pumpOutput: pumpOutput, now: now}
}

func (c *PoolPumpController) GetConfigProperties() ConfigProperties {
	return ConfigProperties{
		Toggles: []ConfigToggleProperties{
			{Name: configKeyEnabled},
		},
		Ranges: []ConfigRangeProperties{
			{
				Name:   configKeyRuntime,
				Degree: "H",
				Min:    1,
				Max:    8,
				Step:   1,
			},
		},
	}
}

func (c PoolPumpController) GetDefaultConfig() Config {
	return Config{
		Toggles: map[string]bool{
			configKeyEnabled: true,
		},
		Ranges: map[string]float64{
			configKeyRuntime: 2,
		},
	}
}

func (c *PoolPumpController) ValidateConfig(config Config) error {
	time, ok := config.Ranges[configKeyRuntime]
	if !ok {
		return fmt.Errorf("Configured runtime missing")
	}
	if time < 1 || time > 8 {
		return fmt.Errorf("Configured type is outside of the allowed range")
	}
	return nil
}

func (c *PoolPumpController) Act(config Config) []EnqueueRequest {
	if !config.Toggles[configKeyEnabled] {
		log.Printf("PoolPumpController: controller is disabled, shutting down outputs\n")
		c.pumpOutput.Set(false)
		return []EnqueueRequest{}
	}
	now := c.now()
	dailyTimerValue := c.timer.Value()
	log.Printf("PoolPumpController: daily timer value: %.2f %s\n", dailyTimerValue, c.timer.Degree())
	durationFromPrevMidnightUntilNextStart := time.Duration(int(24-config.Ranges[configKeyRuntime])) * time.Hour
	prevMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Local().Location())
	nextStart := prevMidnight.Add(durationFromPrevMidnightUntilNextStart)
	nextStop := prevMidnight.Add(24 * time.Hour)
	if config.Ranges[configKeyRuntime] <= dailyTimerValue {
		nextStart = nextStart.Add(24 * time.Hour)
		log.Printf("PoolPumpController: pump time is enough for today, turning off and scheduling turn on in: %s\n", nextStart.Sub(now).String())
		c.pumpOutput.Set(false)
		return []EnqueueRequest{{Controller: c, Config: config, After: nextStart.Sub(now)}}
	}
	nextStart.Add(time.Duration(dailyTimerValue) * time.Hour)
	if now.Before(nextStart) {
		log.Printf("PoolPumpController: we still have time before pump need to run, turning off and scheduling turn on in: %s\n", nextStart.Sub(now).String())
		c.pumpOutput.Set(false)
		return []EnqueueRequest{{Controller: c, Config: config, After: nextStart.Sub(now)}}
	}
	c.pumpOutput.Set(true)
	log.Printf("PoolPumpController: pump time, turning on and scheduling turn off in: %s\n", nextStop.Sub(now).String())
	return []EnqueueRequest{{Controller: c, Config: config, After: nextStop.Sub(now)}}
}

func (c *PoolPumpController) GetName() string {
	return "Pump controller"
}
