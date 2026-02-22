package models

// Port status from nmap scans
type PortStatus uint8

const (
	NMAP_PORT_OPEN = iota + 1
	NMAP_PORT_CLOSED
	NMAP_PORT_FILTERED
)

func (ps PortStatus) String() string {
	return [...]string{"OPEN", "CLOSE", "FILTERED"}[ps]
}

// Validate PortStatus
func (ps PortStatus) IsValid() bool {
	return ps >= NMAP_PORT_OPEN && ps <= NMAP_PORT_FILTERED
}

// NmapDocument is the flat document stored in OpenSearch (one per port)
type NmapDocument struct {
	// Host (IP or domain)
	Host string `json:"host"`
	// Port number
	Port uint16 `json:"port"`
	// Port status (open, close, filtered)
	Status string `json:"status,omitempty"`
	// Service name (http, ssh, etc.)
	ServiceName string `json:"service_name,omitempty"`
	// Service version
	ServiceVersion string `json:"service_version,omitempty"`
}

// NmapData is the input format (host + multiple ports) used in the API
type NmapData struct {
	// IP or domain name
	Host string `json:"host"`
	// Ports with their metadata
	Ports []NmapPort `json:"ports"`
}

// NmapService holds service metadata
type NmapService struct {
	ServiceName    string `json:"service_name"`
	ServiceVersion string `json:"service_version"`
}

// NmapPort is a port with its metadata (used in insert input)
type NmapPort struct {
	Port     uint16      `json:"port"`
	Status   PortStatus  `json:"status,omitempty"`
	MetaData NmapService `json:"metadata"`
}

// ToDocuments converts NmapData into a flat list of NmapDocument (one per port)
func (n *NmapData) ToDocuments() []NmapDocument {
	docs := make([]NmapDocument, 0, len(n.Ports))
	for _, p := range n.Ports {
		docs = append(docs, NmapDocument{
			Host:           n.Host,
			Port:           p.Port,
			Status:         p.Status.String(),
			ServiceName:    p.MetaData.ServiceName,
			ServiceVersion: p.MetaData.ServiceVersion,
		})
	}
	return docs
}
