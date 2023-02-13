package helper

import (
	"os"
	"strconv"
)

func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func GetEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		m, _ := strconv.Atoi(value)
		return m
	}
	return defaultVal
}

func EnvExists(key string) bool {
	value, exists := os.LookupEnv(key)
	return exists && len(value) > 0
}
