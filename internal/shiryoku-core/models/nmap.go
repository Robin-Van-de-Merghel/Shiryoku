package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Adapted from: https://github.com/Ullaakut/nmap/blob/master/xml.go

// PortStatus represents the state of a port.
type PortStatus string

const (
	NMAP_PORT_OPEN            PortStatus = "open"
	NMAP_PORT_CLOSED          PortStatus = "closed"
	NMAP_PORT_FILTERED        PortStatus = "filtered"
	NMAP_PORT_UNFILTERED      PortStatus = "unfiltered"
	NMAP_PORT_OPEN_FILTERED   PortStatus = "open|filtered"
	NMAP_PORT_CLOSED_FILTERED PortStatus = "closed|filtered"
)

func (ps PortStatus) IsValid() bool {
	switch ps {
	case NMAP_PORT_OPEN, NMAP_PORT_CLOSED, NMAP_PORT_FILTERED,
		NMAP_PORT_UNFILTERED, NMAP_PORT_OPEN_FILTERED, NMAP_PORT_CLOSED_FILTERED:
		return true
	default:
		return false
	}
}

// Protocol represents the port protocol.
type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
	// TODO: What is it?
	ProtocolSCTP Protocol = "sctp"
)

func (p Protocol) IsValid() bool {
	switch p {
	case ProtocolTCP, ProtocolUDP, ProtocolSCTP:
		return true
	default:
		return false
	}
}

// NmapScriptResult represents a single NSE script result for a scan result.
type NmapScriptResult struct {
	NmapScriptResultID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"nmap_script_result_id"`
	// TODO: Parsing?
	ScanResultID uuid.UUID `gorm:"type:uuid;index" json:"-"`
	ScriptID     string    `gorm:"type:varchar(255)" json:"id"`
	ScriptOutput string    `gorm:"type:text" json:"output"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"-"`
}

func (NmapScriptResult) TableName() string {
	return "nse_scripts"
}

// Service represents a discovered service
// Unicity : (ServiceName + Product + Version + ServiceExtraInfo + Protocol + ServiceTunnel)
// Usable by other scripts than nmap (nuclei, etc.)
type Service struct {
	ServiceID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"service_id"`
	ServiceName      string    `gorm:"type:varchar(255)" json:"service_name,omitempty"`
	ServiceProduct   string    `gorm:"type:varchar(255)" json:"service_product,omitempty"`
	ServiceVersion   string    `gorm:"type:varchar(255)" json:"service_version,omitempty"`
	ServiceExtraInfo string    `gorm:"type:text" json:"service_extra_info,omitempty"`
	// tcp / udp / sctp
	Protocol      string    `gorm:"type:varchar(10)" json:"protocol,omitempty"`
	ServiceTunnel string    `gorm:"type:varchar(50)" json:"service_tunnel,omitempty"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"-"`
}

func (Service) TableName() string {
	return "services"
}

type NmapScan struct {
	ScanID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"scan_id"`
	// epoch timestamp
	ScanStart   time.Time `gorm:"index:idx_scan_start" json:"scan_start"`
	// command line args
	ScanArgs    string    `gorm:"type:text" json:"scan_args,omitempty"`
	// e.g. "7.94"
	NmapVersion string    `gorm:"type:varchar(50)" json:"nmap_version,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`

	// Relations
	ScanResults []ScanResult `gorm:"foreignKey:ScanID" json:"scan_results,omitempty"`
	Hosts       []NmapHost   `gorm:"many2many:scan_hosts;" json:"hosts,omitempty"`
}

func (NmapScan) TableName() string {
	return "scans"
}

// NmapHost is dedicated to storing host info
// A host can appear in multiple scans
type NmapHost struct {
	HostID     uuid.UUID      `gorm:"type:uuid;primaryKey" json:"host_id"`
	// takes first address
	Host       string         `gorm:"index:idx_host;type:varchar(255)" json:"host"`
	Addresses  pq.StringArray `gorm:"type:text[]" json:"addresses,omitempty"`
	// DNS names
	Hostnames  pq.StringArray `gorm:"type:text[]" json:"hostnames,omitempty"`
	// up / down
	HostStatus string `gorm:"type:varchar(20)" json:"host_status,omitempty"`
	OSName     string `gorm:"type:varchar(255)" json:"os_name,omitempty"`
	OSAccuracy int    `json:"os_accuracy,omitempty"`
	// See nmap doc
	Comment   string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`

	ScanResults []ScanResult `gorm:"foreignKey:HostID" json:"scan_results,omitempty"`
}

func (NmapHost) TableName() string {
	return "hosts"
}

// ScanResult represents a discovery of a service in a scan on a specific host and port
// Composite unique key: (ScanID, HostID, ServiceID, Port)
// One ScanResult = one port + one scan + one host + one service
// Represents: In this scan, this host had this service on this port
type ScanResult struct {
	ScanResultID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"scan_result_id"`
	ScanID       uuid.UUID `gorm:"type:uuid;index:idx_scan_host_service;index:idx_scan_host" json:"scan_id"`
	HostID       uuid.UUID `gorm:"type:uuid;index:idx_scan_host_service;index:idx_scan_host" json:"host_id"`
	ServiceID    uuid.UUID `gorm:"type:uuid;index:idx_scan_host_service" json:"service_id"`
	// Port information
	Port      uint16    `gorm:"index:idx_scan_host_service" json:"port"`
	PortState string    `gorm:"type:varchar(20)" json:"port_state,omitempty"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`

	// Relations - Service loaded manually, no GORM FK constraint
	Scripts  []NmapScriptResult `gorm:"foreignKey:ScanResultID" json:"scripts,omitempty"`
}

func (ScanResult) TableName() string {
	return "scan_results"
}
