package nmap

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
)

// SaveNmapScan saves a single nmap scan (explodes into one doc per port)
func SaveNmapScan(ctx context.Context, nmapData *models.NmapData, nmapDB osdb.NmapDBIface) ([]string, error) {
	return nmapDB.Insert(ctx, nmapData)
}

// SaveNmapScans saves multiple nmap scan results
func SaveNmapScans(ctx context.Context, nmapDataList []models.NmapData, nmapDB osdb.NmapDBIface) ([]string, error) {
	return nmapDB.InsertBatch(ctx, nmapDataList)
}
