package config

import (
	"os"
	"strconv"
	"time"

	config_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config/common"
)

// WorkerConfig contains configuration for background workers
type WorkerConfig struct {
	// Worker name (e.g., "nmap-dashboard-worker")
	Name string
	// Database configuration
	DBConfig DBConfig
	// Refresh frequency in seconds
	Frequency time.Duration
	// Log level
	LogLevel LogLevelT
}

// NewWorkerConfig creates a worker config with defaults from environment variables
func NewWorkerConfig() *WorkerConfig {
	// Default frequency: 300 seconds (5 minutes)
	frequency := 300
	if freqStr := os.Getenv("VIEW_WORK_FREQUENCY"); freqStr != "" {
		if f, err := strconv.Atoi(freqStr); err == nil {
			frequency = f
		}
	}

	// Default log level: DEBUG
	logLevel := LOG_LEVEL_DEBUG
	if levelStr := os.Getenv("LOG_LEVEL"); levelStr != "" {
		switch levelStr {
		case "info":
			logLevel = LOG_LEVEL_INFO
		case "error":
			logLevel = LOG_LEVEL_ERROR
		}
	}

	return &WorkerConfig{
		Name: "nmap-dashboard-worker",
		DBConfig: DBConfig{
			Host:     config_common.GetEnv("DB_HOST", "localhost"),
			Port:     config_common.GetEnvUint16("DB_PORT", 5432),
			Username: config_common.GetEnv("DB_USERNAME", "shiryoku"),
			Password: config_common.GetEnv("DB_PASSWORD", "shiryoku"),
			Database: config_common.GetEnv("DB_NAME", "shiryoku"),
		},
		Frequency: time.Duration(frequency) * time.Second,
		LogLevel:  LogLevelT(logLevel),
	}
}
