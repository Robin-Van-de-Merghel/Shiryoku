package nmap

import (
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logicnmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/Ullaakut/nmap/v4"
	"github.com/gin-gonic/gin"
)

// InsertNmapScans inserts multiple nmap scans
func InsertNmapScans(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		var nmapResults *nmap.Run

		if err := c.ShouldBindJSON(&nmapResults); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		ids, err := logicnmap.SaveNmapScans(c.Request.Context(), nmapResults, nmapDB)
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
