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

const configKeyRuntime = "Desired runtime per day in hours"

func NewPoolPumpController(timer io.Input, pumpOutput io.Output) PoolPumpController {
	return PoolPumpController{timer: timer, pumpOutput: pumpOutput}
}

func (c *PoolPumpController) GetConfigProperties() ConfigProperties {
	return ConfigProperties{
		Toggles: []ConfigToggleProperties{
			{Name: configKeyEnabled},
		},
		Ranges: []ConfigRangeProperties{
			{
				Name: configKeyRuntime,
				Min:  0,
				Max:  8,
				Step: 1,
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
	if time < 0 || time > 8 {
		return fmt.Errorf("Configured type is outside of the allowed range")
	}
	return nil
}

func (c *PoolPumpController) Act(config Config) []EnqueueRequest {
	task := config.Ranges[configKeyRuntime] > (c.timer.Value())
	if c.pumpOutput.Set(task) {
		log.Printf("PoolPumpController: changed pump state to: %t", task)
	}
	return []EnqueueRequest{{Controller: c, Config: config, After: 5 * time.Second}}
}

func (c *PoolPumpController) GetName() string {
	return "Pump controller"
}
