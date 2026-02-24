package config

import (
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/gin-gonic/gin"
)

type LogLevelT uint8

const (
	LOG_LEVEL_DEBUG = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_ERROR
)

func (ll LogLevelT) IsValid() bool {
	return ll <= LOG_LEVEL_ERROR
}

// DB Config for a single DB (SQL, opensearch)
type DBConfig struct {
	// IP/domain/port
	Host string
	Port uint16

	// Creds
	Username string
	Password string
}

// Modules to add modularity
type APIModule struct {
	// Group name
	Name string

	// Description (if needed)
	Description string

	// URL path, e.g. /nmap
	Path string

	// DB instance
	// FIXME: Use generic?
	OSDB osdb.OpenSearchClient

	// Callback to setup the API routes
	// It needs a subgroup. This subgroup separates the endpoints for each module.
	SetupModuleRoutes func(*gin.RouterGroup, osdb.OpenSearchClient, string)
}

// Server config, contains everything that can be modified
type ServerConfig struct {
	// Server port
	Port uint16

	// Logs
	LogLevel LogLevelT

	// Databases
	// [DBName] -> [Config]
	DBConfigs map[string]DBConfig

	// Modules are generic exposed API
	Modules []APIModule

	// Widgets are UI-specific
	Widgets []APIModule
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: 8080,
		LogLevel: LOG_LEVEL_DEBUG,
		DBConfigs: map[string]DBConfig{
			"OSDB": {
				Host: "localhost",
				Port: 9200,
			},
		},
		Modules: []APIModule{},
		Widgets: []APIModule{},
	}
}

