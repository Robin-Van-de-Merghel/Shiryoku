package models

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

// NmapDocument is a flat (easy to store) document for indexing nmap results (one per port).
type NmapDocument struct {
	// Scan metadata
	ScanID      string `json:"scan_id,omitempty"`      // optional unique ID per scan
	ScanStart   int64  `json:"scan_start,omitempty"`   // epoch timestamp
	ScanArgs    string `json:"scan_args,omitempty"`    // command line args
	NmapVersion string `json:"nmap_version,omitempty"` // e.g. "7.94"

	// Host information
	Host       string   `json:"host"`                  // primary IP
	Hostnames  []string `json:"hostnames,omitempty"`   // DNS names
	HostStatus string   `json:"host_status,omitempty"` // up / down
	MACAddress string   `json:"mac_address,omitempty"`
	MACVendor  string   `json:"mac_vendor,omitempty"`

	// OS detection
	OSName     string   `json:"os_name,omitempty"`
	OSAccuracy int      `json:"os_accuracy,omitempty"`
	OSCPE      []string `json:"os_cpe,omitempty"`

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
	ServiceCPE       []string `json:"service_cpe,omitempty"`        // CPE strings

	// NSE scripts
	Scripts []NmapScriptResult `json:"scripts,omitempty"`
}
