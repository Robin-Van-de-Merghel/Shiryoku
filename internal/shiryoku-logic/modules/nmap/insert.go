package logic_nmap

import (
	"context"

	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/Ullaakut/nmap/v4"
)

// SaveNmapScans saves nmap scans (explodes into one doc per port)
func SaveNmapScans(ctx context.Context, nmapData *nmap.Run, nmapDB osdb.OpenSearchClient) ([]*osdb.InsertResult, error) {
	// Convert to documents
	bulkItems := ConvertFullScanIntoDocuments(nmapData)

	osdb.InsertBulk(ctx, &nmapDB, bulkItems.Hosts)
	osdb.InsertBulk(ctx, &nmapDB, bulkItems.Ports)
	osdb.InsertOne(ctx, &nmapDB, bulkItems.Scan)

	// FIXME
	return nil, nil
}
