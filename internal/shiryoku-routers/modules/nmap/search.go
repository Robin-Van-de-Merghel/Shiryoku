package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/gin-gonic/gin"
)

// TODO: Later, use search engine interface

// Gets the last scans
func GetNmapScans(c *gin.Context) {
	// Get data from the database
	nmapResult, _ := nmap.GetLastNmapScans()

	c.JSON(200, nmapResult)
}
