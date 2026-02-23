package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logic_nmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	router_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/common"
	"github.com/gin-gonic/gin"
)

// Multiple scanners are available: scans, hosts and ports

func SearchNmapScans(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
	return router_common.SearchOpenSearch[models.NmapScanDocument](nmapDB, logic_nmap.NMAP_SCANS_INDEX)
}

func SearchNmapPorts(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
	return router_common.SearchOpenSearch[models.NmapPortDocument](nmapDB, logic_nmap.NMAP_PORTS_INDEX)
}

func SearchNmapHosts(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
	return router_common.SearchOpenSearch[models.NmapHostDocument](nmapDB, logic_nmap.NMAP_HOSTS_INDEX)
}
