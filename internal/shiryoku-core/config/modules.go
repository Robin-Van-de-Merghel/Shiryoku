package config

import (
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	routers_utils_setup "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils/setup"
	"github.com/gin-gonic/gin"
)

// GetDefaultModules returns the default API modules with configured routes
func GetDefaultModules(repos *shiryoku_db.Repositories) []APIModule {
	return []APIModule{
		{
			Name:        "nmap",
			Description: "Nmap scan results",
			Path:        "/nmap",
			SetupModuleRoutes: func(group *gin.RouterGroup) {
				routers_utils_setup.SetupNmapRoutes(group, repos)
			},
		},
	}
}
