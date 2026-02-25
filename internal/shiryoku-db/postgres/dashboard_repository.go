package postgres

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"gorm.io/gorm"
)

// DashboardRepositoryImpl implements DashboardRepository interface for dashboard views
type DashboardRepositoryImpl struct {
	db *gorm.DB
}

// NewDashboardRepository creates a new dashboard repository instance
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &DashboardRepositoryImpl{db: db}
}

func (d *DashboardRepositoryImpl) ReadyCheck() common.Checker {
	return func(ctx context.Context) (bool, error) {
		var exists bool
		if err := d.db.Raw(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables 
				WHERE table_schema = current_schema()
				  AND table_name = 'widget_dashboard_scans'
			)
		`).Scan(&exists).Error; err != nil {
			return false, err
		}
		return exists, nil
	}
}

// GetDashboardScans retrieves paginated scan-host combinations from materialized view
func (d *DashboardRepositoryImpl) Search(ctx context.Context, params *models.SearchParams) (uint64, []widgets.WidgetDashboardScan, error) {
		return Search[widgets.WidgetDashboardScan](ctx, d.db, params, 
	)
}

func (d *DashboardRepositoryImpl) RefreshMaterializedView(ctx context.Context) error {
	return d.db.WithContext(ctx).
		Exec("REFRESH MATERIALIZED VIEW CONCURRENTLY dashboard_scans").
		Error
}

// CreateDashboardScans inserts multiple dashboard rows in batch
func (d *DashboardRepositoryImpl) CreateDashboardScans(ctx context.Context, rows []widgets.WidgetDashboardScan) error {
	if len(rows) == 0 {
		return nil
	}

	if err := d.db.WithContext(ctx).CreateInBatches(rows, 500).Error; err != nil {
		return fmt.Errorf("failed to insert dashboard scans: %w", err)
	}

	return nil
}

// TruncateDashboard clears the dashboard table
func (d *DashboardRepositoryImpl) TruncateDashboard(ctx context.Context) error {
	// Truncate all rows, restart identity (auto-increment) and cascade if needed
	if err := d.db.WithContext(ctx).Exec("TRUNCATE TABLE widget_dashboard_scans RESTART IDENTITY CASCADE").Error; err != nil {
		return fmt.Errorf("failed to truncate dashboard table: %w", err)
	}
	return nil
}
