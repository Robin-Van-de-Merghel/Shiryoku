package config

import (
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	routers_utils_setup "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils/setup"
	"github.com/gin-gonic/gin"
)

// GetDefaultWidgets returns the default UI widgets with configured routes
func GetDefaultWidgets(repos *shiryoku_db.Repositories) []APIModule {
	return []APIModule{
		{
			Name:        "Dashboard",
			Description: "Dashboard shown on the first page",
			Path:        "/dashboard",
			SetupModuleRoutes: func(group *gin.RouterGroup) {
				routers_utils_setup.SetupWidgetsDashboardRoutes(group, repos)
			},
		},
	}
}
