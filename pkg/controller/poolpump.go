package controller

import "github.com/libesz/poolmanager/pkg/io"

type PoolPumpController struct {
	timer io.SensingActuator
}

const configKey = "desired runtime per day"

func NewPoolPumpController(timer io.SensingActuator) PoolPumpController {
	return PoolPumpController{timer: timer}
}

func (c *PoolPumpController) GetConfigKeys() []string {
	return []string{configKey}
}

func (c *PoolPumpController) Act(config Config) {
	c.timer.Switch(config[configKey] > c.timer.Value())
}
