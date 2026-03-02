package postgres

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/utils"
)

// DashboardRepository defines operations specific to dashboard views
type DashboardRepository interface {
	// Search retrieves paginated scan-host combinations from materialized view
	Search(ctx context.Context, params *models.SearchParams) (uint64, []widgets.WidgetDashboardScan, error)

	// RefreshMaterializedView refreshes the dashboard materialized view
	RefreshMaterializedView(ctx context.Context) error

	// TruncateDashboard removes everything
	TruncateDashboard(ctx context.Context) error

	// Inserts dashboard scans
	CreateDashboardScans(ctx context.Context, rows []widgets.WidgetDashboardScan) error

	// Check health
	ReadyCheck() utils.Checker
}
