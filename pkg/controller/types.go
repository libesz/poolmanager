package controller

import "time"

type Config map[string]float64

type EnqueueRequest struct {
	Controller Controller
	After      time.Duration
}

type Controller interface {
	Act(Config) []EnqueueRequest
	GetConfigKeys() []string
	GetName() string
}
