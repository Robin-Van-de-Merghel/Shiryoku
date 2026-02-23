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

// GetLatestWidgetHosts fetches hosts with all their scans
func GetLatestWidgetHosts(ctx context.Context, client osdb.OpenSearchClient, input models_widgets.WidgetDashboardInput) ([]models_widgets.WidgetDashboardOutput, error) {
	scans, err := fetchLatestScans(ctx, client, uint64(input.PerPage), input.Page, input.SortDirection)
	if err != nil {
		return nil, err
	}

	if len(scans) == 0 {
		return []models_widgets.WidgetDashboardOutput{}, nil
	}

	hostIDsSet := make(map[string]struct{})
	scanIDs := []string{}
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
						hostIDsSet[hs] = struct{}{} // fallback if no colon
					}
				}
			}
		}
	}
	hostIDs := make([]string, 0, len(hostIDsSet))
	for h := range hostIDsSet {
		hostIDs = append(hostIDs, h)
	}

	hostMap, err := fetchHosts(ctx, client, hostIDs)
	if err != nil {
		return nil, err
	}

	portResults, err := fetchPorts(ctx, client, scanIDs)
	if err != nil {
		return nil, err
	}

	return mergeWidgetHosts(scans, hostMap, portResults), nil
}

func mergeWidgetHosts(scanResults []map[string]any, hostMap map[string]models_widgets.WidgetDashboardOutput, portResults []map[string]any) []models_widgets.WidgetDashboardOutput {
	fmt.Printf("DEBUG mergeWidgetHosts: portResults has %d entries\n", len(portResults))
	fmt.Printf("DEBUG mergeWidgetHosts: hostMap has %d hosts\n", len(hostMap))
	fmt.Printf("DEBUG mergeWidgetHosts: scanResults has %d scans\n", len(scanResults))
	
	type scanKey struct {
		HostID string
		ScanID string
	}
	scanMap := make(map[scanKey]*models_widgets.WidgetDashboardHostScan)

	for _, hit := range portResults {
		hostIDComposite, _ := hit["host_id"].(string)
		scanID, _ := hit["scan_id"].(string)
		portFloat, _ := hit["port"].(float64)
		port := uint16(portFloat)

		// Extract just the host ID part (after the colon)
		hostID := hostIDComposite
		parts := strings.Split(hostIDComposite, ":")
		if len(parts) == 2 {
			hostID = parts[1]
		}

		key := scanKey{HostID: hostID, ScanID: scanID}
		if _, exists := scanMap[key]; !exists {
			var scanStart time.Time
			for _, s := range scanResults {
				if sID, _ := s["scan_id"].(string); sID == scanID {
					if tStr, ok := s["scan_start"].(string); ok {
						if t, err := time.Parse(time.RFC3339, tStr); err == nil {
							scanStart = t
						}
					}
					break
				}
			}
			scanMap[key] = &models_widgets.WidgetDashboardHostScan{
				ScanID:    scanID,
				ScanStart: scanStart,
				Ports:     []uint16{},
			}
		}
		scanMap[key].Ports = append(scanMap[key].Ports, port)
	}

	fmt.Printf("DEBUG mergeWidgetHosts: scanMap has %d entries\n", len(scanMap))

	// Assign scans to their hosts
	assignedCount := 0
	for key, scan := range scanMap {
		if host, ok := hostMap[key.HostID]; ok {
			host.Scans = append(host.Scans, *scan)
			hostMap[key.HostID] = host
			assignedCount++
			fmt.Printf("DEBUG: Assigned scan %s to host %s\n", key.ScanID, key.HostID)
		} else {
			fmt.Printf("DEBUG: Host %s not found in hostMap\n", key.HostID)
		}
	}
	fmt.Printf("DEBUG: Assigned %d scans to hosts\n", assignedCount)

	// Convert host map to slice
	widgets := make([]models_widgets.WidgetDashboardOutput, 0, len(hostMap))
	for _, h := range hostMap {
		widgets = append(widgets, h)
	}
	return widgets
}

// fetchLatestScans fetches scans from OS
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

// fetchHosts fetches host info from OS
func fetchHosts(ctx context.Context, client osdb.OpenSearchClient, hostIDs []string) (map[string]models_widgets.WidgetDashboardOutput, error) {
	if len(hostIDs) == 0 {
		return map[string]models_widgets.WidgetDashboardOutput{}, nil
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

	hostMap := make(map[string]models_widgets.WidgetDashboardOutput)
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
		hostMap[hostID] = models_widgets.WidgetDashboardOutput{
			Host:      hostName,
			HostNames: hostnames,
			Scans:     []models_widgets.WidgetDashboardHostScan{},
		}
	}
	return hostMap, nil
}

// fetchPorts fetches port info from OS
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
		PerPage: uint8(100),
	}
	res, err := client.Search(ctx, logic_nmap.NMAP_PORTS_INDEX, params)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ports: %w", err)
	}
	return res.Results, nil
}
