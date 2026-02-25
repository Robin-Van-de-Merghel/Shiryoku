package postgres

import (
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDB establishes a connection to PostgreSQL and performs schema migrations
func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// AutoMigrate creates tables in dependency order
	// Order matters: tables without foreign keys first, then tables that reference them
	if err := db.AutoMigrate(
		&models.NmapScan{},
		&models.NmapHost{},
		&models.Service{},
		&models.ScanResult{},
		&models.NmapScriptResult{},
		&widgets.WidgetDashboardScan{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate schema: %w", err)
	}

	// Create unique index on Service signature (ServiceName + Product + Version + ExtraInfo + Protocol + Tunnel)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_service_signature 
		ON services(service_name, service_product, service_version, service_extra_info, protocol, service_tunnel)
		WHERE service_name IS NOT NULL
	`).Error; err != nil {
		return nil, fmt.Errorf("failed to create service signature index: %w", err)
	}

	// Create composite unique index on ScanResult (ScanID + HostID + ServiceID + Port)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_scan_result_unique 
		ON scan_results(scan_id, host_id, service_id, port)
	`).Error; err != nil {
		return nil, fmt.Errorf("failed to create scan result unique index: %w", err)
	}

	// Create composite unique index for dashboard table (ScanID + HostID)
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_dashboard_scans_unique
		ON widget_dashboard_scans(scan_id, host_id)
	`).Error; err != nil {
		return nil, fmt.Errorf("failed to create dashboard unique index: %w", err)
	}

	// Optional: index on scan_start for fast queries
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_dashboard_scans_scan_start 
		ON widget_dashboard_scans(scan_start DESC)
	`).Error; err != nil {
		return nil, fmt.Errorf("failed to create dashboard scan_start index: %w", err)
	}

	return db, nil
}
