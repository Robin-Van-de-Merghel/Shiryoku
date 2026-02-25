package config_common

import (
	"context"
	"os"
	"strconv"
)

// Helper functions
func GetEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func GetEnvUint16(key string, defaultVal uint16) uint16 {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.ParseUint(value, 10, 16); err == nil {
			return uint16(v)
		}
	}
	return defaultVal
}

type Checker func(ctx context.Context) (bool, error)
