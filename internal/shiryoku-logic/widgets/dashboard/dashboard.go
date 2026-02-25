package logic_widgets_dashboard

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	logic_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
)

// GetLatestWidgetScans fetches the latest scans with host and port info from dashboard table
func GetLatestWidgetScans(
	ctx context.Context,
	dashboardRepo postgres.DashboardRepository,
	input *models.SearchParams,
) (*logic_common.SearchResult[models_widgets.WidgetDashboardScan], error) {
	input.SetDefaults()

	total, results, err := dashboardRepo.GetDashboardScans(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest scans: %w", err)
	}

	return &logic_common.SearchResult[models_widgets.WidgetDashboardScan]{
		Total:   total,
		Results: results,
	}, nil
}

// BuildDashboardScans rebuilds the dashboard table from the last `days` of scans
func BuildDashboardScans(
	ctx context.Context,
	nmapRepo postgres.NmapRepository,
	dashboardRepo postgres.DashboardRepository,
	days int,
) error {
	// Fetch all recent scans
	params := &models.SearchParams{
		Sort: []models.SortSpec{{Parameter: "scan_start", Direction: "DESC"}},
	}
	_, scans, err := nmapRepo.Search(ctx, params)
	if err != nil {
		return err
	}

	dashboardRows := make([]models_widgets.WidgetDashboardScan, 0, len(scans))

	for _, scan := range scans {
		hosts, err := nmapRepo.GetHosts(ctx, scan.ScanID.String())
		if err != nil {
			return err
		}

		for _, host := range hosts {
			results, err := nmapRepo.GetScanResults(ctx, scan.ScanID.String(), host.HostID.String())
			if err != nil {
				return err
			}

			// Convert scan results to []int
			ports := make([]int, 0, len(results))
			for _, r := range results {
				ports = append(ports, int(r.Port))
			}

			row := models_widgets.WidgetDashboardScan{
				ScanID:     scan.ScanID.String(),
				HostID:     host.HostID.String(),
				ScanStart:  scan.ScanStart,
				Host:       host.Host,
				HostNames:  host.Hostnames,
				Ports:      ports,
				PortNumber: len(ports),
			}

			dashboardRows = append(dashboardRows, row)
		}
	}

	// Flush table
	if err := dashboardRepo.TruncateDashboard(ctx); err != nil {
		return err
	}

	// Insert aggregated data in batch
	if len(dashboardRows) > 0 {
		if err := dashboardRepo.CreateDashboardScans(ctx, dashboardRows); err != nil {
			return err
		}
	}

	return nil
}
