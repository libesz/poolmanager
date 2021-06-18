package controller

import (
	"time"

	"github.com/libesz/poolmanager/pkg/io"
)

type PollController struct {
	inputs   []io.Input
	outputs  []io.Output
	pollTime time.Duration
}

func NewPollController(inputs []io.Input, outputs []io.Output, pollTime time.Duration) PollController {
	return PollController{
		inputs:   inputs,
		outputs:  outputs,
		pollTime: pollTime,
	}
}

func (c *PollController) GetConfigProperties() ConfigProperties {
	return ConfigProperties{}
}

func (c PollController) GetDefaultConfig() Config {
	return Config{}
}

func (c *PollController) ValidateConfig(config Config) error {
	return nil
}

func (c *PollController) Act(config Config) []EnqueueRequest {
	for _, input := range c.inputs {
		_ = input.Value()
	}
	for _, input := range c.outputs {
		_ = input.Get()
	}
	return []EnqueueRequest{{Controller: c, Config: config, After: c.pollTime}}
}

func (c *PollController) GetName() string {
	return "IO Poll controller"
}
