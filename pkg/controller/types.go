package controller

import "time"

type Config map[string]float64

type Controller interface {
	Act(Config) time.Duration
	GetConfigKeys() []string
	GetName() string
}
