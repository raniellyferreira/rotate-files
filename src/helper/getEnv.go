package helper

import (
	"os"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func EnvExists(key string) bool {
	value, exists := os.LookupEnv(key)
	return exists && len(value) > 0
}
