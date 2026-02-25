package common

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

// Checker is a function that we can give to the api endpoint
// If it returns false or an error, the server will be set as "down"
type Checker func(ctx context.Context) (bool, error)
