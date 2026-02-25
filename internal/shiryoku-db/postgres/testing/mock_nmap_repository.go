package testing

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

type MockNmapRepository struct {
	SearchFn             func(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error)
	GetScanFn            func(ctx context.Context, scanID string) (*models.NmapScan, error)
	GetHostsFn           func(ctx context.Context, scanID string) ([]models.NmapHost, error)
	GetScanResultsFn     func(ctx context.Context, scanID, hostID string) ([]models.ScanResult, error)
	GetOrCreateServiceFn func(ctx context.Context, service *models.Service) (*models.Service, error)
	InsertScanFn         func(ctx context.Context, scan *models.NmapScan) error
	InsertHostsFn        func(ctx context.Context, hosts []models.NmapHost) error
	InsertScanResultsFn  func(ctx context.Context, results []models.ScanResult) error
	InsertScriptsFn      func(ctx context.Context, scripts []models.NmapScriptResult) error
	SearchWithHostsFn    func(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error)
}

func (m *MockNmapRepository) Search(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error) {
	if m.SearchFn != nil {
		return m.SearchFn(ctx, params)
	}
	return 0, []models.NmapScan{}, nil
}

func (m *MockNmapRepository) GetScan(ctx context.Context, scanID string) (*models.NmapScan, error) {
	if m.GetScanFn != nil {
		return m.GetScanFn(ctx, scanID)
	}
	return nil, nil
}

func (m *MockNmapRepository) GetHosts(ctx context.Context, scanID string) ([]models.NmapHost, error) {
	if m.GetHostsFn != nil {
		return m.GetHostsFn(ctx, scanID)
	}
	return nil, nil
}

func (m *MockNmapRepository) GetScanResults(ctx context.Context, scanID, hostID string) ([]models.ScanResult, error) {
	if m.GetScanResultsFn != nil {
		return m.GetScanResultsFn(ctx, scanID, hostID)
	}
	return nil, nil
}

func (m *MockNmapRepository) GetOrCreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	if m.GetOrCreateServiceFn != nil {
		return m.GetOrCreateServiceFn(ctx, service)
	}
	return service, nil
}

func (m *MockNmapRepository) InsertScan(ctx context.Context, scan *models.NmapScan) error {
	if m.InsertScanFn != nil {
		return m.InsertScanFn(ctx, scan)
	}
	return nil
}

func (m *MockNmapRepository) InsertHosts(ctx context.Context, hosts []models.NmapHost) error {
	if m.InsertHostsFn != nil {
		return m.InsertHostsFn(ctx, hosts)
	}
	return nil
}

func (m *MockNmapRepository) InsertScanResults(ctx context.Context, results []models.ScanResult) error {
	if m.InsertScanResultsFn != nil {
		return m.InsertScanResultsFn(ctx, results)
	}
	return nil
}

func (m *MockNmapRepository) InsertScripts(ctx context.Context, scripts []models.NmapScriptResult) error {
	if m.InsertScriptsFn != nil {
		return m.InsertScriptsFn(ctx, scripts)
	}
	return nil
}

func (m *MockNmapRepository) SearchWithHosts(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error) {
	if m.SearchWithHostsFn != nil {
		return m.SearchWithHostsFn(ctx, params)
	}
	return 0, []models.NmapScan{}, nil
}
