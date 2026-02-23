package models_widgets

import (
	"time"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

// WidgetNmapHost is being used in the dashboard to display "last hosts scanned"
type WidgetDashboardOutput struct {
	Host      string       `json:"host,omitempty"`
	HostNames []string     `json:"hostnames,omitempty"`
	Scans     []WidgetDashboardHostScan   `json:"scans,omitempty"`
}

// HostScan represents a single scan for a host, including ports
type WidgetDashboardHostScan struct {
	ScanID    string   `json:"scan_id"`
	ScanStart time.Time `json:"scan_start"`
	Ports     []uint16 `json:"ports,omitempty"`
}

// WigetDashboardInput is being sent by the front
type WidgetDashboardInput struct {
	SortDirection       models.SortDirection   `json:"sort"`
	Page       uint64       `json:"page"`
	PerPage    uint8        `json:"per_page"`
}

func (wdi *WidgetDashboardInput) SetDefaults() {
	if wdi.PerPage == 0 {
		wdi.PerPage = models.DEFAULT_RESULTS_PER_PAGE 
	}

	if wdi.SortDirection == "" {
		wdi.SortDirection = models.DirDESC
	}
}
