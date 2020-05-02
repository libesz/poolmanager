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
}

const configKey = "desired runtime per day"

func NewPoolPumpController(timer io.Input, pumpOutput io.Output) PoolPumpController {
	return PoolPumpController{timer: timer, pumpOutput: pumpOutput}
}

func (c *PoolPumpController) GetConfigProperties() ConfigProperties {
	return ConfigProperties{
		configKey: ConfigRange{
			Default: 2,
			Min:     0,
			Max:     8,
			Step:    1,
		},
	}
}

func (c *PoolPumpController) ValidateConfig(config Config) error {
	time, ok := config[configKey].(float64)
	if !ok {
		return fmt.Errorf("Configured type is not int")
	}
	if time < 0 || time > 8 {
		return fmt.Errorf("Configured type is outside of the allowed range")
	}
	return nil
}

func (c *PoolPumpController) Act(config Config) []EnqueueRequest {
	task := config[configKey].(float64) > (c.timer.Value())
	if c.pumpOutput.Set(task) {
		log.Printf("PoolPumpController: changed pump state to: %t", task)
	}
	return []EnqueueRequest{{Controller: c, Config: config, After: 5 * time.Second}}
}

func (c *PoolPumpController) GetName() string {
	return "PoolPumpController"
}
