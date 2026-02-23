package models

import "time"

// Adapted from: https://github.com/Ullaakut/nmap/blob/master/xml.go
// Initially, I stored it in a nested document
// -> Not practical for querying
// -> Not practical for storing
// TODO: Explode multiple ports into multiple documents

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

// NmapScriptResult represents a single NSE script result for a port.
type NmapScriptResult struct {
	ID string `json:"id"`
	// TODO: Parsing?
	Output string `json:"output"`
}

// NmapHostDocument is dedicated to storing host info
type NmapHostDocument struct {
	// Host metadata
	Comment string `json:"comment,omitempty"` // See nmap doc

	// Host information
	HostID string `json:"host_id,omitempty"`
	Host       string   `json:"host"`                  // (takes first address) 
	Addresses []string `json:"addresses,omitempty"`
	Hostnames  []string `json:"hostnames,omitempty"`   // DNS names
	HostStatus string   `json:"host_status,omitempty"` // up / down

	// OS detection
	OSName     string   `json:"os_name,omitempty"`
	OSAccuracy int      `json:"os_accuracy,omitempty"`

	// Not exported
	docID string `json:"-"`
}


// NmapScanDocument stores only scan info (uuid,)
type NmapScanDocument struct {
	ScanID string `json:"scan_id,omitempty"`
	HostIDs []string `json:"host_id,omitempty"`
	ScanStart   time.Time  `json:"scan_start"`   // epoch timestamp
	ScanArgs    string `json:"scan_args,omitempty"`    // command line args
	NmapVersion string `json:"nmap_version,omitempty"` // e.g. "7.94"

	// Not exported
	docID string `json:"-"`

}

// NmapPortDocument aims at storing port results, without host and scan info (for storage sake) 
type NmapPortDocument struct {
	// ScanID and HostID: uuid
	ScanID string `json:"scan_id,omitempty"`
	HostID string `json:"host_id,omitempty"`

	// Port information
	Port      uint16     `json:"port"`
	Protocol  Protocol   `json:"protocol,omitempty"`   // tcp / udp / sctp
	PortState PortStatus `json:"port_state,omitempty"` // open / closed / filtered / ...

	// Service detection
	ServiceName      string   `json:"service_name,omitempty"`       // e.g. http, ssh
	ServiceProduct   string   `json:"service_product,omitempty"`    // product name
	ServiceVersion   string   `json:"service_version,omitempty"`    // version
	ServiceExtraInfo string   `json:"service_extra_info,omitempty"` // extra info string
	ServiceTunnel    string   `json:"service_tunnel,omitempty"`     // e.g. ssl
	// TODO: See if we need more

	// NSE scripts
	Scripts []NmapScriptResult `json:"scripts,omitempty"`

	// Not exported
	docID string `json:"-"`
}
