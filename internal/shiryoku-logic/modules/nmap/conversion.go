package logic_nmap

import (
	"time"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Ullaakut/nmap/v4"
	"github.com/google/uuid"
)

type FullScanResults struct {
	Scan        models.NmapScan
	Hosts       []models.NmapHost
	Services    []models.Service
	ScanResults []models.ScanResult
	Scripts     []models.NmapScriptResult
}

// Converts a full nmap scan into database-usable structs
func ConvertFullScanIntoDocuments(results *nmap.Run) *FullScanResults {
	fullScanResults := &FullScanResults{
		Hosts:       []models.NmapHost{},
		Services:    []models.Service{},
		ScanResults: []models.ScanResult{},
		Scripts:     []models.NmapScriptResult{},
	}

	scanInfo := convertScanInfoToModel(results)
	fullScanResults.Scan = scanInfo

	// Track services we've already seen (for deduplication)
	serviceMap := make(map[string]*models.Service)

	for _, host := range results.Hosts {
		hostItem := convertHostToModel(&host)
		fullScanResults.Hosts = append(fullScanResults.Hosts, hostItem)

		// Convert ports to scan results and services
		scanResultItems, serviceItems := convertHostPortsToModels(&host, hostItem.HostID, scanInfo.ScanID, serviceMap)

		fullScanResults.ScanResults = append(fullScanResults.ScanResults, scanResultItems...)
		fullScanResults.Services = append(fullScanResults.Services, serviceItems...)
	}

	return fullScanResults
}

// Convert all ports info to scan results and services
// Returns both scan results and any new services discovered
// Cf https://github.com/Ullaakut/nmap/blob/5b5552b95453ccf933110e2b48c58cf67160ce1c/xml.go#L175C1-L182C2
func convertHostPortsToModels(h *nmap.Host, hostID uuid.UUID, scanID uuid.UUID, serviceMap map[string]*models.Service) ([]models.ScanResult, []models.Service) {
	var scanResults []models.ScanResult
	var newServices []models.Service

	for _, port := range h.Ports {
		// Create or reference service
		service := models.Service{
			ServiceName:      port.Service.Name,
			ServiceProduct:   port.Service.Product,
			ServiceVersion:   port.Service.HighVersion,
			ServiceExtraInfo: port.Service.ExtraInfo,
			Protocol:         port.Protocol, // Detect protocol from nmap data
			ServiceTunnel:    port.Service.Tunnel,
		}

		// Generate service signature for deduplication
		serviceKey := generateServiceKey(&service)

		// Check if we've already added this service
		var servicePtr *models.Service
		if existingService, exists := serviceMap[serviceKey]; exists {
			servicePtr = existingService
		} else {
			// New service - will be created in DB later and we'll need to get its ID
			// For now, add to our tracking and to the results
			newServices = append(newServices, service)
			servicePtr = &newServices[len(newServices)-1]
			serviceMap[serviceKey] = servicePtr
		}

		// Create scan result
		scanResult := models.ScanResult{
			ScanID:      scanID,
			HostID:      hostID,
			ServiceID:   servicePtr.ServiceID, // Will be set after DB insert
			Port:        port.ID,
			PortState:   string(port.Status()),
		}

		scanResults = append(scanResults, scanResult)

		// Convert scripts for this port/scan result
		// Note: Scripts will be linked to ScanResultID after DB insert
		scripts := convertScripts(port.Scripts)
		// We'll handle script assignment in SaveNmapScans after ScanResults are inserted
		_ = scripts // placeholder until we refactor
	}

	return scanResults, newServices
}

// Generate a unique key for service deduplication
func generateServiceKey(s *models.Service) string {
	return s.ServiceName + "|" + s.ServiceProduct + "|" + s.ServiceVersion + "|" + s.ServiceExtraInfo + "|" + s.Protocol + "|" + s.ServiceTunnel
}

// Convert scan info into single model
// Cf https://github.com/Ullaakut/nmap/blob/5b5552b95453ccf933110e2b48c58cf67160ce1c/xml.go#L15C1-L39C2
func convertScanInfoToModel(si *nmap.Run) models.NmapScan {
	return models.NmapScan{
		ScanID:      uuid.New(),
		ScanArgs:    si.Args,
		NmapVersion: si.Version,
		ScanStart:   time.Time(si.Start),
	}
}

// Convert a host data to a single model
// Cf https://github.com/Ullaakut/nmap/blob/5b5552b95453ccf933110e2b48c58cf67160ce1c/xml.go#L101C1-L121C2
func convertHostToModel(h *nmap.Host) models.NmapHost {
	addresses := convertAddresses(h.Addresses)

	var host string
	if len(addresses) > 0 {
		host = addresses[0] // takes first address
	}

	osName, accuracy := convertOS(h.OS.Matches)

	return models.NmapHost{
		HostID:     uuid.New(),
		Host:       host,
		Addresses:  addresses,
		Hostnames:  convertHostnames(h.Hostnames),
		HostStatus: h.Status.State,
		Comment:    h.Comment,
		OSName:     osName,
		OSAccuracy: accuracy,
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
			ScriptID:     script.ID,
			ScriptOutput: script.Output,
		})
	}

	return scripts
}
