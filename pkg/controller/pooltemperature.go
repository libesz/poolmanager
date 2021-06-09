package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/libesz/poolmanager/pkg/io"
	"github.com/prometheus/client_golang/prometheus"
)

func NewPoolTempController(
	heaterFactor float64,
	tempSensor io.Input,
	heaterOutput io.Output,
	pumpOutput io.Output,
	pollDuration time.Duration,
	now func() time.Time) PoolTempController {
	calculatedDesiredTempGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "poolmanager",
		Subsystem: "input",
		Name:      "calculatedDesiredTemperature",
		Help:      "Temperature",
	})
	prometheus.MustRegister(calculatedDesiredTempGauge)

	return PoolTempController{
		heaterFactor:               heaterFactor,
		tempSensor:                 tempSensor,
		heaterOutput:               heaterOutput,
		pumpOutput:                 pumpOutput,
		now:                        now,
		pollDuration:               pollDuration,
		pendingOperationReady:      make(chan struct{}, 1),
		calculatedDesiredTempGauge: calculatedDesiredTempGauge,
	}
}

type PoolTempController struct {
	heaterFactor               float64
	tempSensor                 io.Input
	heaterOutput               io.Output
	pumpOutput                 io.Output
	pendingOperation           bool
	pendingOperationReady      chan struct{}
	pollDuration               time.Duration
	now                        func() time.Time
	calculatedDesiredTempGauge prometheus.Gauge
}

const (
	configKeyEnabled = "Enabled"
	configKeyTemp    = "Desired temperature"
	configKeyStart   = "Start hour"
	configKeyEnd     = "End hour"
)

func (c PoolTempController) GetConfigProperties() ConfigProperties {
	return ConfigProperties{
		Toggles: []ConfigToggleProperties{
			{
				Name: configKeyEnabled,
			},
		},
		Ranges: []ConfigRangeProperties{
			{
				Name:   configKeyTemp,
				Degree: "°C",
				Min:    20.0,
				Max:    30.0,
				Step:   0.5,
			},
			{
				Name:   configKeyStart,
				Degree: "H",
				Min:    0,
				Max:    23,
				Step:   1,
			},
			{
				Name:   configKeyEnd,
				Degree: "H",
				Min:    0,
				Max:    23,
				Step:   1,
			},
		},
	}
}

func (c PoolTempController) GetDefaultConfig() Config {
	return Config{
		Toggles: map[string]bool{
			configKeyEnabled: false,
		},
		Ranges: map[string]float64{
			configKeyTemp:  26,
			configKeyStart: 10,
			configKeyEnd:   16,
		},
	}
}

func (c PoolTempController) ValidateConfig(config Config) error {
	_, ok := config.Toggles[configKeyEnabled]
	if !ok {
		return fmt.Errorf("Enabled toggle not found in config")
	}
	temp, ok := config.Ranges[configKeyTemp]
	if !ok {
		return fmt.Errorf("Temperature is not found in config")
	}
	if temp < 20.0 || temp > 30.0 {
		return fmt.Errorf("Temperature is outside of the allowed range")
	}
	if int(temp*10)%5 != 0 {
		return fmt.Errorf("Temperature is not a valid step (.0 or .5 required)")
	}
	start, ok := config.Ranges[configKeyStart]
	if !ok {
		return fmt.Errorf("Start time it not found in config")
	}
	if start < 0 || start > 23 {
		return fmt.Errorf("Start time is outside of the allowed range")
	}
	end, ok := config.Ranges[configKeyEnd]
	if !ok {
		return fmt.Errorf("Start time it not found in config")
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
	log.Printf("PoolTempController: delayedOperation set %s to: %t\n", c.output.Name(), c.setTo)
	c.output.Set(c.setTo)
	close(c.pendingOperationReady)
	return nil
}

func (c *delayedOperation) GetName() string {
	return "delayedOperation"
}

func (c delayedOperation) GetConfigProperties() ConfigProperties {
	return ConfigProperties{}
}

func (c delayedOperation) ValidateConfig(Config) error {
	return nil
}

func (c delayedOperation) GetDefaultConfig() Config {
	return EmptyConfig()
}

func (c *PoolTempController) shutdown() []EnqueueRequest {
	c.calculatedDesiredTempGauge.Set(0)
	if c.heaterOutput.Set(false) || c.pumpOutput.Get() {
		c.pendingOperation = true
		pending := delayedOperation{setTo: false, output: c.pumpOutput, pendingOperationReady: c.pendingOperationReady}
		return []EnqueueRequest{{Controller: &pending, After: 6 * time.Second}}
	}
	return []EnqueueRequest{}
}

func (c *PoolTempController) startup() []EnqueueRequest {
	if c.pumpOutput.Set(true) || !c.heaterOutput.Get() {
		c.pendingOperation = true
		pending := delayedOperation{setTo: true, output: c.heaterOutput, pendingOperationReady: c.pendingOperationReady}
		return []EnqueueRequest{{Controller: &pending, After: 6 * time.Second}}
	}
	return []EnqueueRequest{}
}

func (c *PoolTempController) Act(config Config) []EnqueueRequest {
	if c.pendingOperation {
		select {
		case <-c.pendingOperationReady:
			c.pendingOperation = false
			c.pendingOperationReady = make(chan struct{}, 1)
		default:
			log.Printf("PoolTempController: pending operation is in progress")
			return []EnqueueRequest{{Controller: c, Config: config, After: 6 * time.Second}}
		}
	}
	if !config.Toggles[configKeyEnabled] {
		log.Printf("PoolTempController: controller is disabled, shutting down outputs\n")
		return c.shutdown()
	}
	currentTemp := c.tempSensor.Value()
	if currentTemp == io.InputError {
		log.Printf("PoolTempController: temperature value is not available, shutting down outputs for safety\n")
		return append(c.shutdown(), EnqueueRequest{Controller: c, Config: config, After: c.pollDuration})
	}
	now := c.now()
	nextStart := time.Date(now.Year(), now.Month(), now.Day(), int(config.Ranges[configKeyStart]), 0, 0, 0, now.Local().Location())
	nextStop := time.Date(now.Year(), now.Month(), now.Day(), int(config.Ranges[configKeyEnd]), 0, 0, 0, now.Local().Location())

	var thisManyHoursUntilNextStart float64
	if now.After(nextStart) {
		if now.Before(nextStop) {
			nextStart = now
		} else {
			nextStart = nextStart.Add(24 * time.Hour)
		}
	}
	thisManyHoursUntilNextStart = nextStart.Sub(now).Hours()
	calculatedDesiredTemp := config.Ranges[configKeyTemp] - thisManyHoursUntilNextStart*c.heaterFactor
	if calculatedDesiredTemp >= currentTemp {
		c.calculatedDesiredTempGauge.Set(calculatedDesiredTemp)
		log.Printf("PoolTempController: hours back: %.2f. Actual temperature: %.2f %s. Calculated desired temperature: %.2f %s. Need more heat.\n", thisManyHoursUntilNextStart, currentTemp, c.tempSensor.Degree(), calculatedDesiredTemp, c.tempSensor.Degree())
		return append(c.startup(), EnqueueRequest{Controller: c, Config: config, After: c.pollDuration})
	}
	log.Printf("PoolTempController: the actual temperature is already fine: %.2f %s\n", currentTemp, c.tempSensor.Degree())
	return append(c.shutdown(), EnqueueRequest{Controller: c, Config: config, After: c.pollDuration})
}

func (c *PoolTempController) GetName() string {
	return "Temperature controller"
}
