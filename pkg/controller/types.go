package controller

type Config map[string]float64

type Controller interface {
	Act(Config)
	GetConfigKeys() []string
}
