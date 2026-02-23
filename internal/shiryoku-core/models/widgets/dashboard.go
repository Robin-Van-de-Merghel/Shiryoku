package models_widgets

import (
	"time"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

// WidgetDashboardScan represents a single scan with host and port info
type WidgetDashboardScan struct {
	ScanID    string    `json:"scan_id"`
	ScanStart time.Time `json:"scan_start"`
	Host      string    `json:"host"`
	HostID    string    `json:"host_id"`
	HostNames []string  `json:"hostnames,omitempty"`
	Ports     []uint16  `json:"ports,omitempty"`
}

// WidgetHostInfo is internal helper for host data
type WidgetHostInfo struct {
	Host      string
	HostNames []string
}

// WigetDashboardInput is being sent by the front
type WidgetDashboardInput struct {
	SortDirection       models.SortDirection   `json:"sort"`
	Page       uint64       `json:"page"`
	PerPage    uint64        `json:"per_page"`
}

func (wdi *WidgetDashboardInput) SetDefaults() {
	if wdi.PerPage == 0 {
		wdi.PerPage = models.DEFAULT_RESULTS_PER_PAGE 
	}

	if wdi.SortDirection == "" {
		wdi.SortDirection = models.DirDESC
	}
}
