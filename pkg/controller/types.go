package controller

import "time"

type Config map[string]interface{}

type EnqueueRequest struct {
	Controller Controller
	After      time.Duration
}

type Controller interface {
	Act(Config) []EnqueueRequest
	GetConfigKeys() []string
	GetName() string
}
