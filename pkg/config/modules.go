package config

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
	"github.com/gin-gonic/gin"
)

// Modules to add modularity: they are representing api routes
type Module interface {
	Name() string
	Description() string
	SetupRoutes(router *gin.RouterGroup, provider repositories.RepositoryProvider) error
}

// Just aliases for clarity
type APIModule = Module
type Widget = Module
