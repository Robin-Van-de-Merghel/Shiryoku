package setup

import (
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/modules/nmap"
	"github.com/gin-gonic/gin"
)

// SetupNmapRoutes registers all nmap module routes
func SetupNmapRoutes(nmap_group *gin.RouterGroup, repos *shiryoku_db.Repositories) {
	search_group := nmap_group.Group("/search")
	search_group.POST("", nmap.SearchNmapScans(repos.Nmap))
	search_group.POST("/ports", nmap.SearchNmapPorts(repos.Nmap))
	nmap_group.POST("/batch", nmap.InsertNmapScans(repos.Nmap))
}
