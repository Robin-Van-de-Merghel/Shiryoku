package nmap

import "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"

// List last scans
func GetLastNmapScans(searchParamas models.SearchParams) ([]models.NmapData, error) {
	// TODO: Add db layer to store it
	nmapResult := []models.NmapData{
		{
			Host: "1.1.1.1",
			Ports: []models.NmapPort{
				{
					Port: 443,
					MetaData: models.NmapService{
						ServiceName:    "HTTPS",
						ServiceVersion: "XX",
					},
				},
			},
		},
	}

	return nmapResult, nil
}
