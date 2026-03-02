package routers

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/status"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/utils"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/widgets/dashboard"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/config"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, provider repositories.RepositoryProvider) error {
	// Middlewares
	router.Use(utils.ErrorRecoveryMiddleware())

	// For docker-compose status
	dashboardRepo := provider.GetRepository("dashboard").(postgres.DashboardRepository)
	nmapRepo := provider.GetRepository("nmap").(repositories.NmapRepository)
	router.GET("/ping", status.Ping(
		dashboardRepo.ReadyCheck(),
		nmapRepo.ReadyCheck(),
	))

	// API generic group
	api_group := router.Group("/api")
	{
		// Modules group
		modules_group := api_group.Group("/modules")
		for _, module := range getDefaultModules() {
			// Create the custom module (e.g. /nmap)
			current_group := modules_group.Group(module.Name())
			// Import all routes
			if err := module.SetupRoutes(current_group, provider); err != nil {
				return err
			}
		}
	}
	{
		// Widgets group
		dashboard_group := api_group.Group("/widgets")
		for _, widget := range getDefaultWidgets() {
			// Create the custom widget (e.g. /last_scans)
			current_group := dashboard_group.Group(widget.Name())
			// Import all routes
			if err := widget.SetupRoutes(current_group, provider); err != nil {
				return err
			}
		}
	}

	return nil
}

// getDefaultModules returns the default API modules
func getDefaultModules() []config.APIModule {
	return []config.Module{
		&nmap.NmapModule{},
	}
}

// getDefaultWidgets returns the default widget modules
func getDefaultWidgets() []config.Widget {
	return []config.Module{
		&dashboard.DashboardWidget{},
	}
}
