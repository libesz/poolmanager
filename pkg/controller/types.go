package controller

import "time"

type ConfigRangeProperties struct {
	Name   string
	Degree string
	Min    float64
	Max    float64
	Step   float64
}

type ConfigToggleProperties struct {
	Name string
}

type ConfigProperties struct {
	Ranges  []ConfigRangeProperties
	Toggles []ConfigToggleProperties
}

func EmptyProperties() ConfigProperties {
	return ConfigProperties{
		Ranges:  []ConfigRangeProperties{},
		Toggles: []ConfigToggleProperties{},
	}
}

type Config struct {
	Ranges  map[string]float64
	Toggles map[string]bool
}

func EmptyConfig() Config {
	return Config{
		Ranges:  make(map[string]float64),
		Toggles: make(map[string]bool),
	}
}

type EnqueueRequest struct {
	Controller Controller
	Config     Config
	After      time.Duration
}

type Controller interface {
	Act(Config) []EnqueueRequest
	GetConfigProperties() ConfigProperties
	GetDefaultConfig() Config
	ValidateConfig(Config) error
	GetName() string
}
