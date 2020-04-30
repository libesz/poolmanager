package controller

import "time"

type ConfigProperty struct {
	Min     interface{}
	Max     interface{}
	Default interface{}
}

type ConfigProperties map[string]ConfigProperty

type Config map[string]interface{}

type EnqueueRequest struct {
	Controller Controller
	After      time.Duration
}

type Controller interface {
	Act(Config) []EnqueueRequest
	GetConfig() ConfigProperties
	ValidateConfig(Config) error
	GetName() string
}
