package shiryoku_routers

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/status"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

func GetFilledRouter(serverConfig config.ServerConfig, repos *shiryoku_db.Repositories) *gin.Engine {

	// Main router
	router := gin.Default()

	// Middlewares
	router.Use(utils.ErrorRecoveryMiddleware())

	// For docker-compose status
	router.GET("/ping", status.Ping(
		repos.Dashboard.ReadyCheck(),
	))

	// API generic group
	api_group := router.Group("/api")

	{
		// Modules group
		modules_group := api_group.Group("/modules")

		for _, module := range serverConfig.Modules {
			// Create the custom module (e.g. /nmap)
			current_group := modules_group.Group(module.Path)

			// Import all routes
			module.SetupModuleRoutes(current_group)
		}
	}

	{
		// Widgets group
		dashboard_group := api_group.Group("/widgets")

		for _, module := range serverConfig.Widgets {
			// Create the custom widget (e.g. /last_scans)
			current_group := dashboard_group.Group(module.Path)

			// Import all routes
			module.SetupModuleRoutes(current_group)
		}
	}

	return router
}
