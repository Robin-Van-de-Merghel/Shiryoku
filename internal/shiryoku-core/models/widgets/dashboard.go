package models_widgets

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// WidgetDashboardScan represents a single scan-host combination with discovered ports
// Maps to the table: widget_dashboard_scans
type WidgetDashboardScan struct {
	ScanID    string         `gorm:"column:scan_id;primaryKey" json:"scan_id"`
	HostID    string         `gorm:"column:host_id;primaryKey" json:"host_id"`
	ScanStart time.Time      `gorm:"column:scan_start" json:"scan_start"`
	Host      string         `gorm:"column:host" json:"host"`
	PortNumber int `gorm:"column:port_number;type:integer" json:"-"`

	// Stored in DB as pq.IntArray or pq.StringArray
	PGPorts pq.Int64Array `gorm:"column:ports;type:integer[]" json:"-"`
	PGHostnames pq.StringArray `gorm:"column:hostnames;type:text[]" json:"-"`

	// Go-friendly helper slice, not persisted
	Ports []int `gorm:"-" json:"ports,omitempty"`
	HostNames []string       `gorm:"-" json:"hostnames,omitempty"`
}

// Hooks for syncing
func (w *WidgetDashboardScan) AfterFind(tx *gorm.DB) error {
    w.Ports = make([]int, len(w.PGPorts))
    for i, p := range w.PGPorts {
			w.Ports[i] = int(p)
    }

    w.HostNames = make([]string, len(w.PGHostnames))
    copy(w.HostNames, w.PGHostnames)

    return nil
}

func (w *WidgetDashboardScan) BeforeSave(tx *gorm.DB) error {
	w.PGPorts = make(pq.Int64Array, len(w.Ports))
	for i, p := range w.Ports {
		w.PGPorts[i] = int64(p)
	}

	w.PGHostnames = make(pq.StringArray, len(w.HostNames))
	copy(w.PGHostnames, w.HostNames)

	return nil
}

// TableName specifies the table name for GORM
func (WidgetDashboardScan) TableName() string {
	return "widget_dashboard_scans"
}
