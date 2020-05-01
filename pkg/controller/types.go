package controller

import "time"

type ConfigRange struct {
	Min     float64
	Max     float64
	Step    float64
	Default float64
}

type ConfigToggle struct {
	Default bool
}

type ConfigProperties map[string]interface{}

type Config map[string]interface{}

type EnqueueRequest struct {
	Controller Controller
	After      time.Duration
}

type Controller interface {
	Act(Config) []EnqueueRequest
	GetConfigProperties() ConfigProperties
	ValidateConfig(Config) error
	GetName() string
}
