package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	"github.com/gin-gonic/gin"
)

func SetupNmapRoutes(nmap_group *gin.RouterGroup, nmapModuleConfig config.APIModule) {
	// Routes
	nmap_group.POST("/search", SearchNmapScans(nmapModuleConfig.NMapDB))

	// FIXME: Refactorize
	// nmap_group.POST("/", InsertNmapScan(nmapModuleConfig.NMapDB))
	nmap_group.POST("/batch", InsertNmapScans(nmapModuleConfig.NMapDB))
}
