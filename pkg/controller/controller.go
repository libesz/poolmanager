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

func IsEqualConfig(a, b Config) bool {
	if len(a.Toggles) != len(b.Toggles) || len(a.Ranges) != len(b.Ranges) {
		return false
	}
	for i, v := range a.Toggles {
		valueInB, ok := b.Toggles[i]
		if !ok {
			return false
		}
		if valueInB != v {
			return false
		}
	}
	for i, v := range a.Ranges {
		valueInB, ok := b.Ranges[i]
		if !ok {
			return false
		}
		if valueInB != v {
			return false
		}
	}
	return true
}
