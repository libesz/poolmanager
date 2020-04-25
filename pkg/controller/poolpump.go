package controller

import (
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

type PoolPumpController struct {
	timer io.InputOutput
}

const configKey = "desired runtime per day"

func NewPoolPumpController(timer io.InputOutput) PoolPumpController {
	return PoolPumpController{timer: timer}
}

func (c *PoolPumpController) GetConfigKeys() []string {
	return []string{configKey}
}

func (c *PoolPumpController) Act(config Config) time.Duration {
	c.timer.Switch(config[configKey] > c.timer.Value())
	return 5 * time.Second
}

func (c *PoolPumpController) GetName() string {
	return "PoolPumpController"
}
