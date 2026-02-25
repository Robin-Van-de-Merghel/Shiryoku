package routers_utils_setup

import (
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	routers_modules_nmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/modules/nmap"
	"github.com/gin-gonic/gin"
)

// SetupNmapRoutes registers all nmap module routes
func SetupNmapRoutes(nmap_group *gin.RouterGroup, repos *shiryoku_db.Repositories) {
	search_group := nmap_group.Group("/search")
	search_group.POST("", routers_modules_nmap.SearchNmapScans(repos.Nmap))
	search_group.POST("/ports", routers_modules_nmap.SearchNmapPorts(repos.Nmap))
	nmap_group.POST("/batch", routers_modules_nmap.InsertNmapScans(repos.Nmap))
}
