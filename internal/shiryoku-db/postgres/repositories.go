package postgres

import (
	"context"
	"fmt"

	config_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"gorm.io/gorm"
)

// SearchableRepository is a generic interface for any resource that supports search and pagination
type SearchableRepository[T any] interface {
	Search(ctx context.Context, params *models.SearchParams) (uint64, []T, error)
}

// NmapRepository defines database operations specific to nmap scan data
type NmapRepository interface {
	SearchableRepository[models.NmapScan]

	// GetScan retrieves a single scan with all its scan results
	GetScan(ctx context.Context, scanID string) (*models.NmapScan, error)

	// GetHosts retrieves all unique hosts from a scan
	GetHosts(ctx context.Context, scanID string) ([]models.NmapHost, error)

	// GetScanResults retrieves all scan results (ports discovered) for a specific scan and host
	GetScanResults(ctx context.Context, scanID, hostID string) ([]models.ScanResult, error)

	// GetOrCreateService retrieves or creates a service by its signature
	// (ServiceName + Product + Version + ExtraInfo + Protocol + Tunnel)
	GetOrCreateService(ctx context.Context, service *models.Service) (*models.Service, error)

	// InsertScan inserts a new scan
	InsertScan(ctx context.Context, scan *models.NmapScan) error

	// InsertHosts inserts or updates hosts (upsert by host IP)
	InsertHosts(ctx context.Context, hosts []models.NmapHost) error

	// InsertScanResults inserts scan result records (ports discovered in a scan on a host with a service)
	InsertScanResults(ctx context.Context, results []models.ScanResult) error

	// InsertScripts inserts NSE script results
	InsertScripts(ctx context.Context, scripts []models.NmapScriptResult) error
}

// DashboardRepository defines operations specific to dashboard views
type DashboardRepository interface {
	// Search retrieves paginated scan-host combinations from materialized view
	Search(ctx context.Context, params *models.SearchParams) (uint64, []models_widgets.WidgetDashboardScan, error)

	// RefreshMaterializedView refreshes the dashboard materialized view
	RefreshMaterializedView(ctx context.Context) error

	// TruncateDashboard removes everything
	TruncateDashboard(ctx context.Context) error

	// Inserts dashboard scans
	CreateDashboardScans(ctx context.Context, rows []models_widgets.WidgetDashboardScan) error

	// Check health
	ReadyCheck() config_common.Checker
}

type Preload[T any] struct {
    Association string
    Fn          func(*gorm.DB) *gorm.DB
}

// Search is a generic function to query a simple table with SearchParams
// It accepts preloads fields, as some results may be nested
func Search[T any](
    ctx context.Context, 
    db *gorm.DB, 
    params *models.SearchParams, 
    preloads ...Preload[T],
) (uint64, []T, error) {
	var results []T

	params.SetDefaults()

	builder := NewSearchBuilder[T](db.WithContext(ctx))
	query, err := builder.Build(params)
	if err != nil {
		return 0, nil, err
	}

	var total int64
	if err := query.Model(new(T)).Count(&total).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to count records: %w", err)
	}

	// Handle pagination - page 0 defaults to 1
	page := params.Page
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * params.PerPage
	query = query.Offset(int(offset)).Limit(int(params.PerPage))

	// Apply preloads if provided
	for _, preload := range preloads {
		query = query.Preload(preload.Association, preload.Fn)
	}

	if err := query.Find(&results).Error; err != nil {
		return 0, nil, fmt.Errorf("failed to search records: %w", err)
	}

	return uint64(total), results, nil
}
