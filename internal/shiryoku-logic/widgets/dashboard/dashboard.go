package logic_widgets_dashboard

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	logic_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
)

// GetLatestWidgetScans fetches the latest scans with host and port info from materialized view
func GetLatestWidgetScans(ctx context.Context, dashboardRepo postgres.DashboardRepository, input models_widgets.WidgetDashboardInput) (*logic_common.SearchResult[models_widgets.WidgetDashboardScan], error) {
	input.SetDefaults()

	// Build search params for the materialized view
	searchParams := &models.SearchParams{
		PerPage: input.PerPage,
		Page:    input.Page,
		Sort: []models.SortSpec{
			{Parameter: "scan_start", Direction: models.SortDirection(input.SortDirection)},
		},
	}

	total, results, err := dashboardRepo.GetDashboardScans(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest scans: %w", err)
	}

	return &logic_common.SearchResult[models_widgets.WidgetDashboardScan]{
		Total:   total,
		Results: results,
	}, nil
}
