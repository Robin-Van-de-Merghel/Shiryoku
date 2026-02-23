package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logicnmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

// // InsertNmapScan inserts a single nmap scan (one doc per port)
// func InsertNmapScan(nmapDB osdb.NmapDBIface) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		var nmapData models.NmapDocument
//
// 		if err := c.ShouldBindJSON(&nmapData); err != nil {
// 			utils.ParseJSONError(c, err)
// 			return
// 		}
//
// 		if nmapData.Host == "" {
// 			c.JSON(400, gin.H{"error": "host field is required"})
// 			return
// 		}
//
// 		// ids, err := logicnmap.SaveNmapScan(c.Request.Context(), &nmapData, nmapDB)
// 		ids, err := logicnmap.SaveNmapScan(c.Request.Context(), &nmapData, nmapDB)
// 		if err != nil {
// 			c.JSON(500, gin.H{"error": err.Error()})
// 			return
// 		}
//
// 		c.JSON(201, gin.H{
// 			"ids":     ids,
// 			"count":   len(ids),
// 			"message": "nmap scan inserted successfully",
// 		})
// 	}
// }

// InsertNmapScans inserts multiple nmap scans
func InsertNmapScans(nmapDB osdb.NmapDBIface) func(c *gin.Context) {
	return func(c *gin.Context) {
		var nmapDocumentList []models.NmapDocument

		if err := c.ShouldBindJSON(&nmapDocumentList); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		ids, err := logicnmap.SaveNmapScans(c.Request.Context(), nmapDocumentList, nmapDB)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"ids":     ids,
			"count":   len(ids),
			"message": "nmap scans inserted successfully",
		})
	}
}
