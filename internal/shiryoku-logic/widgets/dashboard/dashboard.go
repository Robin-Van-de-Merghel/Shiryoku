package dashboard

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
)

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

	dashboardRows := make([]widgets.WidgetDashboardScan, 0, len(scans))

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

			row := widgets.WidgetDashboardScan{
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
