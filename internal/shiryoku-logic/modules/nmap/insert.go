package logic_nmap

import (
	"context"
	"fmt"
	
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/Ullaakut/nmap/v4"
)

// SaveNmapScans saves nmap scans (explodes into one doc per port)
func SaveNmapScans(ctx context.Context, nmapData *nmap.Run, nmapDB osdb.OpenSearchClient) ([]*osdb.InsertResult, error) {
	bulkItems := ConvertFullScanIntoDocuments(nmapData)
	
	fmt.Printf("DEBUG SaveNmapScans: hosts=%d, ports=%d\n", len(bulkItems.Hosts), len(bulkItems.Ports))
	
	if len(bulkItems.Hosts) > 0 {
		_, err := osdb.InsertBulk(ctx, &nmapDB, bulkItems.Hosts)
		if err != nil {
			return nil, fmt.Errorf("failed to insert hosts: %w", err)
		}
		fmt.Printf("DEBUG: Inserted %d hosts\n", len(bulkItems.Hosts))
	}
	
	if len(bulkItems.Ports) > 0 {
		results, err := osdb.InsertBulk(ctx, &nmapDB, bulkItems.Ports)
		if err != nil {
			return nil, fmt.Errorf("failed to insert ports: %w", err)
		}
		fmt.Printf("DEBUG: Inserted %d ports\n", len(bulkItems.Ports))
		for i, result := range results {
			fmt.Printf("DEBUG: Port result[%d]: ID=%s, Index=%s, Version=%d\n", i, result.ID, result.Index, result.Version)
		}
	}
	
	result, err := osdb.InsertOne(ctx, &nmapDB, bulkItems.Scan)
	if err != nil {
		return nil, fmt.Errorf("failed to insert scan: %w", err)
	}
	
	return []*osdb.InsertResult{result}, nil
}
