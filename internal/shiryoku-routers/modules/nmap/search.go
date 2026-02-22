package nmap

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/gin-gonic/gin"
)

// TODO: Later, use search engine interface

func GetNmapScans(c *gin.Context) {
	// TODO: Add logic layer to fetch from DB
	nmapResult := []models.NmapData{
		{
			Host: "1.1.1.1",
			Ports: []models.NmapPort{
				{
					Port: 443,
					MetaData: models.NmapService{
						ServiceName: "HTTPS",
						ServiceVersion: "XX",
					},	
				},
			},
		},
	}

	c.JSON(200, nmapResult)
}
