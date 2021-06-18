package controller

import "time"

type ConfigRangeProperties struct {
	Name   string  `json:"name"`
	Degree string  `json:"degree"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	Step   float64 `json:"step"`
}

type ConfigToggleProperties struct {
	Name string `json:"name"`
}

type ConfigProperties struct {
	Ranges  []ConfigRangeProperties  `json:"ranges"`
	Toggles []ConfigToggleProperties `json:"toggles"`
}

func EmptyProperties() ConfigProperties {
	return ConfigProperties{
		Ranges:  []ConfigRangeProperties{},
		Toggles: []ConfigToggleProperties{},
	}
}

type Config struct {
	Ranges  map[string]float64 `json:"ranges"`
	Toggles map[string]bool    `json:"toggles"`
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
