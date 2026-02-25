package config

import (
	"os"

	config_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config/common"
	"github.com/gin-gonic/gin"
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

// DB Config for a single DB (SQL, opensearch)
type DBConfig struct {
	// IP/domain
	Host string

	// Port number
	Port uint16

	// Creds
	Username string
	Password string

	// Database name
	Database string
}

// Modules to add modularity
type APIModule struct {
	// Group name
	Name string

	// Description (if needed)
	Description string

	// URL path, e.g. /nmap
	Path string

	// Callback to setup the API routes
	// It needs a subgroup. This subgroup separates the endpoints for each module.
	SetupModuleRoutes func(*gin.RouterGroup)
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
		Port:     config_common.GetEnvUint16("PORT", 8080),
		LogLevel: logLevel,
		DBConfig: DBConfig{
			Host:     config_common.GetEnv("DB_HOST", "localhost"),
			Port:     config_common.GetEnvUint16("DB_PORT", 5432),
			Username: config_common.GetEnv("DB_USERNAME", "shiryoku"),
			Password: config_common.GetEnv("DB_PASSWORD", "shiryoku"),
			Database: config_common.GetEnv("DB_NAME", "shiryoku"),
		},
		Modules: []APIModule{},
		Widgets: []APIModule{},
	}
}
