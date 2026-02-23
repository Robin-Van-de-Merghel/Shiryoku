package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	"github.com/gin-gonic/gin"
)

func SetupNmapRoutes(nmap_group *gin.RouterGroup, nmapModuleConfig config.APIModule) {
	// Routes
	search_group := nmap_group.Group("/search")

	search_group.POST("/scans", SearchNmapScans(nmapModuleConfig.OSDB))
	search_group.POST("/hosts", SearchNmapHosts(nmapModuleConfig.OSDB))
	search_group.POST("/ports", SearchNmapPorts(nmapModuleConfig.OSDB))

	// FIXME: Refactorize
	nmap_group.POST("/batch", InsertNmapScans(nmapModuleConfig.OSDB))
}
