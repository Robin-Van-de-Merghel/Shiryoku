package config

import (
	"os"
)

type LogLevelT uint8

const (
	LOG_LEVEL_DEBUG LogLevelT = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_ERROR
)

func (ll LogLevelT) IsValid() bool {
	return ll <= LOG_LEVEL_ERROR
}

// Server config, contains everything that can be modified
type ServerConfig struct {
	// Server port
	Port uint16

	// Logs
	LogLevel LogLevelT

	// Database (one SQL)
	DBConfig DBConfig

	// Modules are generic exposed API
	Modules []APIModule

	// Widgets are UI-specific
	Widgets []APIModule
}

// NewServerConfig creates a ServerConfig instance using environment variables.
func NewServerConfig() *ServerConfig {
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

	return &ServerConfig{
		Port:     GetEnvUint16("PORT", 8080),
		LogLevel: logLevel,
		DBConfig: DBConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnvUint16("DB_PORT", 5432),
			Username: GetEnv("DB_USERNAME", "shiryoku"),
			Password: GetEnv("DB_PASSWORD", "shiryoku"),
			Database: GetEnv("DB_NAME", "shiryoku"),
		},
		Modules: []APIModule{},
		Widgets: []APIModule{},
	}
}
