package nmap

import "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"

const LAST_SCAN_NUMBER = 10

// List last scans
func GetLastNmapScans() ([]models.NmapData, error) {
	// TODO: Add db layer to store it 
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

	return nmapResult, nil
}
