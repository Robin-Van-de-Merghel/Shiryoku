package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/common"
	"github.com/gin-gonic/gin"
)

// SearchNmapScans returns a handler for searching nmap scans
func SearchNmapScans(nmapRepo postgres.NmapRepository) gin.HandlerFunc {
	return common.Search[models.NmapScan](nmapRepo)
}

// SearchNmapPorts returns a handler for searching nmap ports
func SearchNmapPorts(nmapRepo postgres.NmapRepository) gin.HandlerFunc {
	return common.Search[models.NmapScan](nmapRepo)
}
