package shiryoku_routers

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/status"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

func GetFilledRouter() *gin.Engine {
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

		{
			// Nmap group
			nmap_group := modules_group.Group("/nmap")

			// Routes
			nmap_group.POST("/search", nmap.SearchNmapScans)
		}
	}

	return router
}
