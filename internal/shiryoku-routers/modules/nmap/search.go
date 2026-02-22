package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

// TODO: Later, use search engine interface

// Gets the last scans
func SearchNmapScans(c *gin.Context) {
	var params models.SearchParams

	if err := c.ShouldBindJSON(&params); err != nil {
		utils.ParseJSONError(c, err)
		return
	}

	if !utils.ValidateAndRespond(c, params, utils.SearchSchema) {
		return
	}

	// To prevent having "per_page=0"
	params.SetDefaults()

	// Get data from the database
	nmapResult, err := nmap.GetLastNmapScans(params)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, nmapResult)
}
