package controller

import "time"

type ConfigRangeProperties struct {
	Min     float64
	Max     float64
	Step    float64
	Default float64
}

type ConfigToggleProperties struct {
	Default bool
}

type ConfigProperties struct {
	Ranges  map[string]ConfigRangeProperties
	Toggles map[string]ConfigToggleProperties
}

func EmptyProperties() ConfigProperties {
	return ConfigProperties{
		Ranges:  make(map[string]ConfigRangeProperties),
		Toggles: make(map[string]ConfigToggleProperties),
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
	ValidateConfig(Config) error
	GetName() string
}
