package repositories

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/utils"
)

// NmapRepository defines database operations specific to nmap scan data
type NmapRepository interface {
	SearchableRepository[models.NmapScan]

	// GetScan retrieves a single scan with all its scan results
	GetScan(ctx context.Context, scanID string) (*models.NmapScan, error)

	// GetHosts retrieves all unique hosts from a scan
	GetHosts(ctx context.Context, scanID string) ([]models.NmapHost, error)

	// GetScanResults retrieves all scan results (ports discovered) for a specific scan and host
	GetScanResults(ctx context.Context, scanID, hostID string) ([]models.ScanResult, error)

	// GetOrCreateService retrieves or creates a service by its signature
	// (ServiceName + Product + Version + ExtraInfo + Protocol + Tunnel)
	GetOrCreateService(ctx context.Context, service *models.Service) (*models.Service, error)

	// InsertScan inserts a new scan
	InsertScan(ctx context.Context, scan *models.NmapScan) error

	// InsertHosts inserts or updates hosts (upsert by host IP)
	InsertHosts(ctx context.Context, hosts []models.NmapHost) error

	// InsertScanResults inserts scan result records (ports discovered in a scan on a host with a service)
	InsertScanResults(ctx context.Context, results []models.ScanResult) error

	// InsertScripts inserts NSE script results
	InsertScripts(ctx context.Context, scripts []models.NmapScriptResult) error

	ReadyCheck() utils.Checker
}
