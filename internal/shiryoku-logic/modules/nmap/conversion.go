package logic_nmap

import (
	"fmt"
	"time"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/Ullaakut/nmap/v4"
	"github.com/google/uuid"
)

type FullScanResults struct {
	Hosts []osdb.BulkItem[models.NmapHostDocument]
	Ports []osdb.BulkItem[models.NmapPortDocument]
	Scan  osdb.BulkItem[models.NmapScanDocument]
}

// Converts a full nmap scan into an opensearch-usable struct
func ConvertFullScanIntoDocuments(results *nmap.Run) *FullScanResults {
	fullScanResults := &FullScanResults{
		Hosts: []osdb.BulkItem[models.NmapHostDocument]{},
		Ports: []osdb.BulkItem[models.NmapPortDocument]{},
	}

	scanInfo := convertScanInfoToDocuments(results)

	for _, host := range results.Hosts {
		hostItem := convertHostToDocument(&host, scanInfo.ID)
		fullScanResults.Hosts = append(fullScanResults.Hosts, osdb.BulkItem[models.NmapHostDocument]{
			Index: NMAP_HOSTS_INDEX,
			ID:    hostItem.ID,
			Doc:   hostItem.Doc,
		})

		scanInfo.Doc.HostIDs = append(scanInfo.Doc.HostIDs, hostItem.ID)

		portItems := convertHostPortsToDocuments(&host, scanInfo.ID, hostItem.ID)
		fmt.Printf("DEBUG: Host %s has %d ports\n", hostItem.ID, len(portItems))
		
		fullScanResults.Ports = append(fullScanResults.Ports, portItems...)
	}

	fmt.Printf("DEBUG: Total ports created: %d\n", len(fullScanResults.Ports))

	fullScanResults.Scan = osdb.BulkItem[models.NmapScanDocument]{
		Index: NMAP_SCANS_INDEX,
		ID:    scanInfo.ID,
		Doc:   scanInfo.Doc,
	}

	return fullScanResults
}

// Convert all ports info to multiple documents
// Cf https://github.com/Ullaakut/nmap/blob/5b5552b95453ccf933110e2b48c58cf67160ce1c/xml.go#L175C1-L182C2
func convertHostPortsToDocuments(h *nmap.Host, scanID string, hostID string) []osdb.BulkItem[models.NmapPortDocument] {
	var bulkItems []osdb.BulkItem[models.NmapPortDocument]

	for _, port := range h.Ports {
		var doc models.NmapPortDocument

		doc.ScanID = scanID
		doc.HostID = hostID
		doc.Port = port.ID
		doc.PortState = models.PortStatus(port.Status())
		doc.Scripts = convertScripts(port.Scripts)

		// Service metadata
		doc.ServiceName = port.Service.Name
		doc.ServiceVersion = port.Service.HighVersion
		doc.ServiceExtraInfo = port.Service.ExtraInfo
		doc.ServiceProduct = port.Service.Product
		doc.ServiceTunnel = port.Service.Tunnel

		bulkItems = append(bulkItems, osdb.BulkItem[models.NmapPortDocument]{
			ID:    fmt.Sprintf("%s:%s:%d", scanID, hostID, doc.Port),
			Index: NMAP_PORTS_INDEX,
			Doc:   doc,
		})
	}

	return bulkItems
}

// Convert scan info into single document
// Cf https://github.com/Ullaakut/nmap/blob/5b5552b95453ccf933110e2b48c58cf67160ce1c/xml.go#L15C1-L39C2
func convertScanInfoToDocuments(si *nmap.Run) osdb.BulkItem[models.NmapScanDocument] {
	var doc models.NmapScanDocument

	// Unique ID to find
	doc.ScanID = uuid.NewString()
	doc.HostIDs = []string{} // Filled later
	doc.ScanArgs = si.Args
	doc.NmapVersion = si.Version
	doc.ScanStart = time.Time(si.Start)

	return osdb.BulkItem[models.NmapScanDocument]{
		ID:    doc.ScanID,
		Index: NMAP_SCANS_INDEX,
		Doc: doc,
	}
}

// Convert a host data to a single document
// Cf https://github.com/Ullaakut/nmap/blob/5b5552b95453ccf933110e2b48c58cf67160ce1c/xml.go#L101C1-L121C2
func convertHostToDocument(h *nmap.Host, scanID string) osdb.BulkItem[models.NmapHostDocument] {
	var doc models.NmapHostDocument

	doc.Addresses = convertAddresses(h.Addresses)
	// Host = first address
	doc.HostID = uuid.NewString()
	doc.Host = doc.Addresses[0]
	doc.Hostnames = convertHostnames(h.Hostnames)
	doc.Addresses = convertAddresses(h.Addresses)
	doc.HostStatus = h.Status.State
	doc.Comment = h.Comment

	osName, accuracy := convertOS(h.OS.Matches)
	doc.OSName = osName
	doc.OSAccuracy = accuracy

	return osdb.BulkItem[models.NmapHostDocument]{
		Index: NMAP_HOSTS_INDEX,
		ID:    fmt.Sprintf("%s:%s", scanID, doc.HostID),
		Doc:   doc,
	}
}

// Get all ips
func convertAddresses(addresses []nmap.Address) []string {
	var ips []string

	for _, address := range addresses {
		ips = append(ips, address.Addr)
	}

	return ips
}

// Get all hostnames
func convertHostnames(hostnames []nmap.Hostname) []string {
	var names []string

	for _, hostname := range hostnames {
		names = append(names, hostname.Name)
	}

	return names
}

// Extracts top match. Returns the top OS, with its accuracy
func convertOS(matches []nmap.OSMatch) (string, int) {
	if len(matches) == 0 {
		return "", 0
	}

	maxAcc := matches[0].Accuracy
	maxAccName := matches[0].Name

	// FIXME: Remove first iteration
	for _, match := range matches {
		if maxAcc <= match.Accuracy {
			maxAcc = match.Accuracy
			maxAccName = match.Name
		}
	}

	return maxAccName, maxAcc
}

func convertScripts(rawScripts []nmap.Script) []models.NmapScriptResult {
	var scripts []models.NmapScriptResult

	for _, script := range rawScripts {
		scripts = append(scripts, models.NmapScriptResult{
			ID:     script.ID,
			Output: script.Output,
		})
	}

	return scripts
}
