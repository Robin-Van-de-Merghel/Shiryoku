package logic_widgets_dashboard

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	logic_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
)

// GetLatestWidgetScans fetches the latest scans with host and port info
func GetLatestWidgetScans(ctx context.Context, nmapRepo postgres.NmapRepository, input models_widgets.WidgetDashboardInput) (*logic_common.SearchResult[models_widgets.WidgetDashboardScan], error) {
	// Fetch ALL scans without pagination (we'll paginate results after grouping)
	searchParams := &models.SearchParams{
		PerPage: 1000, // Large limit to get all scans
		Page:    1,
		Sort: []models.SortSpec{
			{Parameter: "scan_start", Direction: input.SortDirection},
		},
	}

	_, scans, err := nmapRepo.Search(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest scans: %w", err)
	}

	if len(scans) == 0 {
		return &logic_common.SearchResult[models_widgets.WidgetDashboardScan]{
			Total:   0,
			Results: []models_widgets.WidgetDashboardScan{},
		}, nil
	}

	// Build all dashboard scans
	allResults := buildDashboardScans(scans, nmapRepo, ctx)

	// Apply pagination on the results
	total := uint64(len(allResults))
	page := input.Page
	if page == 0 {
		page = 1
	}
	
	offset := (page - 1) * input.PerPage
	var paginatedResults []models_widgets.WidgetDashboardScan
	
	if offset < uint64(len(allResults)) {
		end := offset + input.PerPage
		if end > uint64(len(allResults)) {
			end = uint64(len(allResults))
		}
		paginatedResults = allResults[offset:end]
	}

	return &logic_common.SearchResult[models_widgets.WidgetDashboardScan]{
		Total:   total,
		Results: paginatedResults,
	}, nil
}

// buildDashboardScans converts scan models into flat WidgetDashboardScan objects
// Groups scan results by (ScanID, HostID) and loads host data
func buildDashboardScans(scans []models.NmapScan, nmapRepo postgres.NmapRepository, ctx context.Context) []models_widgets.WidgetDashboardScan {
	var results []models_widgets.WidgetDashboardScan

	// Map to group scan results by (ScanID, HostID)
	type scanHostKey struct {
		scanID string
		hostID string
	}
	scanHostGroups := make(map[scanHostKey]*models_widgets.WidgetDashboardScan)
	hostIDsToLoad := make(map[string]bool) // Track unique host IDs to load

	// Collect all unique scan-host combinations and their ports
	for _, scan := range scans {
		for _, scanResult := range scan.ScanResults {
			key := scanHostKey{
				scanID: scan.ScanID.String(),
				hostID: scanResult.HostID.String(),
			}

			// Initialize widget entry if not exists
			if _, exists := scanHostGroups[key]; !exists {
				scanHostGroups[key] = &models_widgets.WidgetDashboardScan{
					ScanID:    scan.ScanID.String(),
					ScanStart: scan.ScanStart,
					HostID:    scanResult.HostID.String(),
					Ports:     []uint16{},
				}
				hostIDsToLoad[scanResult.HostID.String()] = true
			}

			// Append port to this scan-host combination
			scanHostGroups[key].Ports = append(scanHostGroups[key].Ports, scanResult.Port)
		}
	}

	// Load host data for all unique hosts
	hostDataMap := make(map[string]*models.NmapHost)
	for hostID := range hostIDsToLoad {
		// For now, we'll populate from the scans if available
		// Better solution: query hosts by IDs, but that requires repo method
		hostDataMap[hostID] = nil
	}

	// Populate host data - query each scan for its hosts
	for _, scan := range scans {
		hosts, err := nmapRepo.GetHosts(ctx, scan.ScanID.String())
		if err == nil {
			for _, host := range hosts {
				hostDataMap[host.HostID.String()] = &host
			}
		}
	}

	// Apply host data to widgets
	for _, widget := range scanHostGroups {
		if host, exists := hostDataMap[widget.HostID]; exists && host != nil {
			widget.Host = host.Host
			widget.HostNames = host.Hostnames
		}
	}

	// Convert map to slice
	for _, widget := range scanHostGroups {
		results = append(results, *widget)
	}

	return results
}
