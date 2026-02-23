package nmap

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
)

// SaveNmapScans saves nmap scans (explodes into one doc per port)
func SaveNmapScans(ctx context.Context, nmapData []models.NmapDocument, nmapDB osdb.NmapDBIface) ([]string, error) {
	return nmapDB.Insert(ctx, nmapData)
}
