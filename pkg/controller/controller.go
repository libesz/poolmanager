package controller

func CopyConfig(orig Config) Config {
	result := Config{
		Ranges:  make(map[string]float64),
		Toggles: make(map[string]bool),
	}
	for key, value := range orig.Ranges {
		result.Ranges[key] = value
	}
	for key, value := range orig.Toggles {
		result.Toggles[key] = value
	}
	return result
}

func IsEmptyConfig(config Config) bool {
	return len(config.Ranges) == 0 && len(config.Toggles) == 0
}
