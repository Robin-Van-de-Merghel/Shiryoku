package postgres

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
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

// GetDashboardScans retrieves paginated scan-host combinations from materialized view
func (d *DashboardRepositoryImpl) GetDashboardScans(ctx context.Context, params *models.SearchParams) (uint64, []models_widgets.WidgetDashboardScan, error) {
	var results []models_widgets.WidgetDashboardScan

	params.SetDefaults()

	// Query the materialized view directly
	query := d.db.WithContext(ctx)

	var total int64
	if err := query.Model(&models_widgets.WidgetDashboardScan{}).Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count dashboard scans: %w", err)
	}

	// Handle pagination - page 0 defaults to 1
	page := params.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * params.PerPage
	query = query.Offset(int(offset)).Limit(int(params.PerPage))

	// Apply sorting
	if len(params.Sort) > 0 {
		for _, sort := range params.Sort {
			direction := sort.Direction
			if direction == "" {
				direction = "ASC"
			}
			query = query.Order(fmt.Sprintf("\"%s\" %s", sort.Parameter, direction))
		}
	}

	// Execute query
	if err := query.Find(&results).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to fetch dashboard scans: %w", err)
	}

	return uint64(total), results, nil
}
