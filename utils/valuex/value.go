package valuex

import "time"

func GetOrDefault[T comparable](value, defaultValue T) T {
	var zero T
	if value != zero {
		return value
	}
	return defaultValue
}

func GetDurationOrDefault(value int, defaultValue time.Duration) time.Duration {
	if value > 0 {
		return time.Duration(value) * time.Second
	}
	return defaultValue
}
