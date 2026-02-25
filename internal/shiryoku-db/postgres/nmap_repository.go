package postgres

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NmapRepositoryImpl implements NmapRepository interface for nmap scan data persistence
type NmapRepositoryImpl struct {
	db *gorm.DB
}

func NewNmapRepository(db *gorm.DB) NmapRepository {
	return &NmapRepositoryImpl{db: db}
}

func (n *NmapRepositoryImpl) Search(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error) {
    return Search[models.NmapScan](ctx, n.db, params,
        Preload[models.NmapScan]{Association: "ScanResults", Fn: func(db *gorm.DB) *gorm.DB {
            return db.Preload("Scripts")
        }},
        Preload[models.NmapScan]{Association: "Hosts", Fn: nil},
    )
}

func (n *NmapRepositoryImpl) GetScan(ctx context.Context, scanID string) (*models.NmapScan, error) {
	var scan models.NmapScan
	if err := n.db.WithContext(ctx).
		Where("scan_id = ?", scanID).
		Preload("ScanResults", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Scripts")
		}).
		First(&scan).Error; err != nil {
		return nil, fmt.Errorf("failed to get scan: %w", err)
	}
	return &scan, nil
}

func (n *NmapRepositoryImpl) GetHosts(ctx context.Context, scanID string) ([]models.NmapHost, error) {
	var hosts []models.NmapHost
	if err := n.db.WithContext(ctx).
		Joins("JOIN scan_results ON hosts.host_id = scan_results.host_id").
		Where("scan_results.scan_id = ?", scanID).
		Distinct("hosts.*").
		Find(&hosts).Error; err != nil {
		return nil, fmt.Errorf("failed to get hosts: %w", err)
	}
	return hosts, nil
}

// GetScanResults fetches all results for a specific scan and host
func (n *NmapRepositoryImpl) GetScanResults(ctx context.Context, scanID, hostID string) ([]models.ScanResult, error) {
	var results []models.ScanResult
	if err := n.db.WithContext(ctx).
		Where("scan_id = ? AND host_id = ?", scanID, hostID).
		Preload("Scripts").
		Find(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get scan results: %w", err)
	}
	return results, nil
}

// GetOrCreateService finds or creates a service by its signature
func (n *NmapRepositoryImpl) GetOrCreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	result := &models.Service{}
	err := n.db.WithContext(ctx).
		Where("service_name = ? AND service_product = ? AND service_version = ? AND service_extra_info = ? AND protocol = ? AND service_tunnel = ?",
			service.ServiceName, service.ServiceProduct, service.ServiceVersion, service.ServiceExtraInfo, service.Protocol, service.ServiceTunnel).
		FirstOrCreate(result, service).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get or create service: %w", err)
	}
	return result, nil
}

// InsertScan inserts a complete scan
func (n *NmapRepositoryImpl) InsertScan(ctx context.Context, scan *models.NmapScan) error {
	return n.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Insert scan
		if err := tx.Create(scan).Error; err != nil {
			return fmt.Errorf("failed to insert scan: %w", err)
		}
		return nil
	})
}

// InsertHosts inserts or updates hosts (upsert by host_id)
func (n *NmapRepositoryImpl) InsertHosts(ctx context.Context, hosts []models.NmapHost) error {
	if len(hosts) == 0 {
		return nil
	}
	// Use OnConflict to handle duplicate hosts (merge addresses and hostnames)
	if err := n.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "host_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"addresses", "hostnames", "host_status", "os_name", "os_accuracy", "comment"}),
		}).
		CreateInBatches(hosts, 100).Error; err != nil {
		return fmt.Errorf("failed to insert hosts: %w", err)
	}
	return nil
}

// InsertScanResults inserts scan result records (ports discovered in a scan on a host with a service)
func (n *NmapRepositoryImpl) InsertScanResults(ctx context.Context, results []models.ScanResult) error {
	if len(results) == 0 {
		return nil
	}
	if err := n.db.WithContext(ctx).CreateInBatches(results, 100).Error; err != nil {
		return fmt.Errorf("failed to insert scan results: %w", err)
	}
	return nil
}

// InsertScripts inserts NSE script results
func (n *NmapRepositoryImpl) InsertScripts(ctx context.Context, scripts []models.NmapScriptResult) error {
	if len(scripts) == 0 {
		return nil
	}
	if err := n.db.WithContext(ctx).CreateInBatches(scripts, 100).Error; err != nil {
		return fmt.Errorf("failed to insert scripts: %w", err)
	}
	return nil
}
