package routers_utils_setup

import (
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	routers_modules_nmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/modules/nmap"
	"github.com/gin-gonic/gin"
)

func SetupNmapRoutes(
	nmap_group *gin.RouterGroup,
	OSDB osdb.OpenSearchClient,
	path string,
) {
	// Routes
	search_group := nmap_group.Group("/search")

	search_group.POST("/scans", routers_modules_nmap.SearchNmapScans(OSDB))
	search_group.POST("/hosts", routers_modules_nmap.SearchNmapHosts(OSDB))
	search_group.POST("/ports", routers_modules_nmap.SearchNmapPorts(OSDB))

	// FIXME: Refactorize
	nmap_group.POST("/batch", routers_modules_nmap.InsertNmapScans(OSDB))
}
