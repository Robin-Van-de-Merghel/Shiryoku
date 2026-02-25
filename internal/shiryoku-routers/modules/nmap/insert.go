package routers_modules_nmap

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	logicnmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/Ullaakut/nmap/v4"
	"github.com/gin-gonic/gin"
)

// InsertNmapScans inserts multiple nmap scans
func InsertNmapScans(nmapRepo postgres.NmapRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		var nmapResults *nmap.Run

		// Read raw XML body from the request
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("failed to read request body: %v", err)})
			return
		}

		// Unmarshal XML data into nmap.Run struct
		err = xml.Unmarshal(body, &nmapResults)
		if err != nil {
			c.JSON(400, gin.H{"error": fmt.Sprintf("invalid XML format: %v", err)})
			return
		}

		// Continue as before, now you have `nmapResults` unmarshalled from XML
		ids, err := logicnmap.SaveNmapScans(c.Request.Context(), nmapResults, nmapRepo)
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
