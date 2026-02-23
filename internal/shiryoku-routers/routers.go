package shiryoku_routers

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/status"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

func GetFilledRouter(nmapDB *osdb.OpenSearchClient, modules []config.APIModule) *gin.Engine {

	// Main router
	router := gin.Default()

	// Middlewares
	router.Use(utils.ErrorRecoveryMiddleware())

	// For docker-compose status
	router.GET("/ping", status.Ping)

	// API generic group
	api_group := router.Group("/api")

	{
		// Modules group
		modules_group := api_group.Group("/modules")

		for _, module := range modules {
			// Create the custom module (e.g. /nmap)
			current_group := modules_group.Group(module.Path)

			// Import all routes
			module.SetupModuleRoutes(current_group, module)
		}
	}

	return router
}
