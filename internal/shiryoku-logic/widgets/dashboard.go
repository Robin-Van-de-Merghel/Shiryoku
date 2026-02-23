package logic_widget

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logic_nmap "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/modules/nmap"
)

// GetLatestWidgetScans fetches the latest scans with host and port info
func GetLatestWidgetScans(ctx context.Context, client osdb.OpenSearchClient, input models_widgets.WidgetDashboardInput) ([]models_widgets.WidgetDashboardScan, error) {
	// Fetch latest scans
	scans, err := fetchLatestScans(ctx, client, uint64(input.PerPage), input.Page, input.SortDirection)
	if err != nil {
		return nil, err
	}

	if len(scans) == 0 {
		return []models_widgets.WidgetDashboardScan{}, nil
	}

	// Extract scan IDs and host IDs from scans
	scanIDs := make([]string, 0, len(scans))
	hostIDsSet := make(map[string]struct{})

	for _, s := range scans {
		scanID, _ := s["scan_id"].(string)
		scanIDs = append(scanIDs, scanID)

		if hosts, ok := s["host_id"].([]any); ok {
			for _, h := range hosts {
				if hs, ok := h.(string); ok {
					// Extract just the host ID part (after the colon)
					parts := strings.Split(hs, ":")
					if len(parts) == 2 {
						hostIDsSet[parts[1]] = struct{}{}
					} else {
						hostIDsSet[hs] = struct{}{}
					}
				}
			}
		}
	}

	// Convert host set to slice
	hostIDs := make([]string, 0, len(hostIDsSet))
	for h := range hostIDsSet {
		hostIDs = append(hostIDs, h)
	}

	// Fetch host info
	hostMap, err := fetchHosts(ctx, client, hostIDs)
	if err != nil {
		return nil, err
	}

	// Fetch ports for all scans
	portResults, err := fetchPorts(ctx, client, scanIDs)
	if err != nil {
		return nil, err
	}

	// Build results
	return buildDashboardScans(scans, hostMap, portResults), nil
}

// buildDashboardScans converts raw scan/host/port data into flat WidgetDashboardScan objects
func buildDashboardScans(scanResults []map[string]any, hostMap map[string]models_widgets.WidgetHostInfo, portResults []map[string]any) []models_widgets.WidgetDashboardScan {
	// Group ports by (scanID, hostID)
	type scanKey struct {
		ScanID string
		HostID string
	}
	portMap := make(map[scanKey][]uint16)

	for _, port := range portResults {
		hostIDComposite, _ := port["host_id"].(string)
		scanID, _ := port["scan_id"].(string)
		portFloat, _ := port["port"].(float64)

		// Extract just the host ID part (after the colon)
		hostID := hostIDComposite
		parts := strings.Split(hostIDComposite, ":")
		if len(parts) == 2 {
			hostID = parts[1]
		}

		key := scanKey{ScanID: scanID, HostID: hostID}
		portMap[key] = append(portMap[key], uint16(portFloat))
	}

	// Build scan objects
	results := make([]models_widgets.WidgetDashboardScan, 0, len(scanResults))
	for _, s := range scanResults {
		scanID, _ := s["scan_id"].(string)
		scanStartStr, _ := s["scan_start"].(string)

		// Parse scan start time
		var scanStart time.Time
		if scanStartStr != "" {
			if t, err := time.Parse(time.RFC3339, scanStartStr); err == nil {
				scanStart = t
			}
		}

		// Get hosts for this scan
		var hosts []any
		if h, ok := s["host_id"].([]any); ok {
			hosts = h
		}

		// Create a scan entry for each host in the scan
		for _, hostInterface := range hosts {
			hostIDComposite, ok := hostInterface.(string)
			if !ok {
				continue
			}

			// Extract just the host ID part
			hostID := hostIDComposite
			parts := strings.Split(hostIDComposite, ":")
			if len(parts) == 2 {
				hostID = parts[1]
			}

			// Get host info
			hostInfo, hostExists := hostMap[hostID]
			if !hostExists {
				continue
			}

			// Get ports for this scan+host combination
			key := scanKey{ScanID: scanID, HostID: hostID}
			ports := portMap[key]

			results = append(results, models_widgets.WidgetDashboardScan{
				ScanID:    scanID,
				ScanStart: scanStart,
				Host:      hostInfo.Host,
				HostID:    hostID,
				HostNames: hostInfo.HostNames,
				Ports:     ports,
			})
		}
	}

	return results
}

// fetchLatestScans fetches scans from OpenSearch
func fetchLatestScans(ctx context.Context, client osdb.OpenSearchClient, perPage, pageNumber uint64, sortDir models.SortDirection) ([]map[string]any, error) {
	params := &models.SearchParams{
		PerPage:    uint8(perPage),
		Page:       pageNumber,
		Sort:       []models.SortSpec{{Parameter: "scan_start", Direction: sortDir}},
		Parameters: []string{"scan_id", "scan_start", "host_id"},
	}
	res, err := client.Search(ctx, logic_nmap.NMAP_SCANS_INDEX, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest scans: %w", err)
	}
	return res.Results, nil
}

// WidgetHostInfo is used internally to map host data
type WidgetHostInfo struct {
	Host      string
	HostNames []string
}

// fetchHosts fetches host info from OpenSearch
func fetchHosts(ctx context.Context, client osdb.OpenSearchClient, hostIDs []string) (map[string]models_widgets.WidgetHostInfo, error) {
	if len(hostIDs) == 0 {
		return map[string]models_widgets.WidgetHostInfo{}, nil
	}

	params := &models.SearchParams{
		Parameters: []string{"host_id", "host", "hostnames"},
		Search: []models.SearchSpec{
			{
				Vector: &models.VectorSearchSpec{
					Parameter: "host_id",
					Operator:  models.OpIn,
					Values:    hostIDs,
				},
			},
		},
		PerPage: uint8(100),
	}
	results, err := client.Search(ctx, logic_nmap.NMAP_HOSTS_INDEX, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch host info: %w", err)
	}

	hostMap := make(map[string]models_widgets.WidgetHostInfo)
	for _, h := range results.Results {
		hostID, _ := h["host_id"].(string)
		hostName, _ := h["host"].(string)
		var hostnames []string
		if hns, ok := h["hostnames"].([]any); ok {
			for _, hn := range hns {
				if hs, ok := hn.(string); ok {
					hostnames = append(hostnames, hs)
				}
			}
		}
		hostMap[hostID] = models_widgets.WidgetHostInfo{
			Host:      hostName,
			HostNames: hostnames,
		}
	}
	return hostMap, nil
}

// fetchPorts fetches port info from OpenSearch
func fetchPorts(ctx context.Context, client osdb.OpenSearchClient, scanIDs []string) ([]map[string]any, error) {
	if len(scanIDs) == 0 {
		return []map[string]any{}, nil
	}

	params := &models.SearchParams{
		Parameters: []string{"scan_id", "host_id", "port"},
		Search: []models.SearchSpec{
			{
				Vector: &models.VectorSearchSpec{
					Parameter: "scan_id",
					Operator:  models.OpIn,
					Values:    scanIDs,
				},
			},
		},
		PerPage: uint8(255), // Maximum uint8 value to get all ports in one query
	}
	res, err := client.Search(ctx, logic_nmap.NMAP_PORTS_INDEX, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ports: %w", err)
	}
	return res.Results, nil
}
