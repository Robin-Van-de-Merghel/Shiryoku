package logic_nmap

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	"github.com/Ullaakut/nmap/v4"
)

// SaveNmapScans saves nmap scans with proper service deduplication and cascading relationships
func SaveNmapScans(ctx context.Context, nmapData *nmap.Run, nmapRepo postgres.NmapRepository) ([]string, error) {
	bulkItems := ConvertFullScanIntoDocuments(nmapData)

	// 1. Insert or reference hosts (upsert by host IP)
	if len(bulkItems.Hosts) > 0 {
		if err := nmapRepo.InsertHosts(ctx, bulkItems.Hosts); err != nil {
			return nil, fmt.Errorf("failed to insert hosts: %w", err)
		}
	}

	// 2. Get or create services (dedup by signature)
	serviceSignatureMap := make(map[string]*models.Service) // signature -> Service with populated ServiceID
	for i, service := range bulkItems.Services {
		createdService, err := nmapRepo.GetOrCreateService(ctx, &service)
		if err != nil {
			return nil, fmt.Errorf("failed to get or create service: %w", err)
		}
		// Store the service with its populated ServiceID
		signature := generateServiceKey(createdService)
		serviceSignatureMap[signature] = createdService
		bulkItems.Services[i] = *createdService
	}

	// 3. Update scan results with actual service IDs
	for i := range bulkItems.ScanResults {
		// ServiceID is already set from conversion, just need to populate from map
		// Find matching service by iterating through created services
		for _, service := range bulkItems.Services {
			if service.ServiceName == "" && service.ServiceProduct == "" {
				// Skip empty services
				continue
			}
			// Use the service ID from the created service
			bulkItems.ScanResults[i].ServiceID = service.ServiceID
			break
		}
	}

	// 4. Insert scan
	if err := nmapRepo.InsertScan(ctx, &bulkItems.Scan); err != nil {
		return nil, fmt.Errorf("failed to insert scan: %w", err)
	}

	// 5. Insert scan results
	if len(bulkItems.ScanResults) > 0 {
		if err := nmapRepo.InsertScanResults(ctx, bulkItems.ScanResults); err != nil {
			return nil, fmt.Errorf("failed to insert scan results: %w", err)
		}
	}

	// 6. Insert scripts (after scan results exist)
	if len(bulkItems.Scripts) > 0 {
		if err := nmapRepo.InsertScripts(ctx, bulkItems.Scripts); err != nil {
			return nil, fmt.Errorf("failed to insert scripts: %w", err)
		}
	}

	return []string{bulkItems.Scan.ScanID.String()}, nil
}
