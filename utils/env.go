package utils

import (
	"os"
)

func GetEnvDefault(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
