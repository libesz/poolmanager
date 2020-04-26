package controller

import (
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

func (c *PoolPumpController) GetConfigKeys() []string {
	return []string{configKey}
}

func (c *PoolPumpController) Act(config Config) []EnqueueRequest {
	task := config[configKey] > c.timer.Value()
	if c.pumpOutput.Switch(task) {
		log.Printf("PoolPumpController: changed pump state to: %t", task)
	}
	return []EnqueueRequest{{Controller: c, After: 5 * time.Second}}
}

func (c *PoolPumpController) GetName() string {
	return "PoolPumpController"
}
