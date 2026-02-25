package postgres_testing

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

type MockNmapRepository struct {
	SearchFn      func(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error)
	GetScanFn     func(ctx context.Context, scanID string) (*models.NmapScan, error)
	GetHostsFn    func(ctx context.Context, scanID string) ([]models.NmapHost, error)
	GetPortsFn    func(ctx context.Context, scanID, hostID string) ([]models.NmapPort, error)
	InsertScanFn  func(ctx context.Context, scan *models.NmapScan) error
	InsertHostsFn func(ctx context.Context, hosts []models.NmapHost) error
	InsertPortsFn func(ctx context.Context, ports []models.NmapPort) error
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

func (m *MockNmapRepository) GetPorts(ctx context.Context, scanID, hostID string) ([]models.NmapPort, error) {
	if m.GetPortsFn != nil {
		return m.GetPortsFn(ctx, scanID, hostID)
	}
	return nil, nil
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

func (m *MockNmapRepository) InsertPorts(ctx context.Context, ports []models.NmapPort) error {
	if m.InsertPortsFn != nil {
		return m.InsertPortsFn(ctx, ports)
	}
	return nil
}
