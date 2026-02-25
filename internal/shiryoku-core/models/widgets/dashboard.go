package models_widgets

import (
	"time"

	"github.com/lib/pq"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

// WidgetDashboardScan represents a single scan-host combination with discovered ports
// Maps to the materialized view: dashboard_scans
type WidgetDashboardScan struct {
	ScanID    string         `gorm:"column:scan_id;primaryKey" json:"scan_id"`
	HostID    string         `gorm:"column:host_id;primaryKey" json:"host_id"`
	ScanStart time.Time      `gorm:"column:scan_start" json:"scan_start"`
	Host      string         `gorm:"column:host" json:"host"`
	HostNames pq.StringArray `gorm:"column:hostnames;type:text[]" json:"hostnames,omitempty"`
	Ports     pq.Int64Array  `gorm:"column:ports;type:integer[]" json:"ports,omitempty"`
}

// TableName specifies the materialized view name
func (WidgetDashboardScan) TableName() string {
	return "dashboard_scans"
}

// WidgetDashboardInput represents the request parameters from the frontend
type WidgetDashboardInput struct {
	SortDirection string `json:"sort"`
	Page          uint64 `json:"page"`
	PerPage       uint64 `json:"per_page"`
}

// SetDefaults sets default values for pagination and sorting
func (wdi *WidgetDashboardInput) SetDefaults() {
	if wdi.PerPage == 0 {
		wdi.PerPage = models.DEFAULT_RESULTS_PER_PAGE
	}
	if wdi.SortDirection == "" {
		wdi.SortDirection = string(models.DirDESC)
	}
}
