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

// Services
type NmapService struct {
	// Service name (e.g. HTTP)
	ServiceName string

	// Service Version
	ServiceVersion string
}

// Port with its MetaData
type NmapPort struct {
	// Max port: 65,535
	Port uint16

	// Status

	// More info about this port
	MetaData NmapService
}

// (not-)All information obtained from an Nmap scan
type NmapData struct {
	// IP or domain name
	Host string
	
	// Ports with their metadata
	Ports []NmapPort

	// TODO: More info (OS, traceroute, modules, etc.)
}
