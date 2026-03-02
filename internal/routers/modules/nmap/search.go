package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/utils"
	"github.com/gin-gonic/gin"
)

// SearchNmapScans returns a handler for searching nmap scans
func (m *NmapModule) searchNmapScans() gin.HandlerFunc {
	return common.Search(m.nmapRepo, utils.NmapScanFields)
}
