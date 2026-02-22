package osdb

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

// NmapDBIface is the interface that both NmapDB and DummyNmapDB implement.
// Use this type everywhere instead of *NmapDB to allow test injection.
type NmapDBIface interface {
	Search(ctx context.Context, params *models.SearchParams) (*NmapSearchResult, error)
	Insert(ctx context.Context, nmapData *models.NmapData) ([]string, error)
	InsertBatch(ctx context.Context, nmapDataList []models.NmapData) ([]string, error)
}

// DummyNmapDB is a no-op implementation of NmapDBIface for unit tests.
// Always returns empty results and no errors.
type DummyNmapDB struct{}

func NewDummyNmapDB() NmapDBIface {
	return &DummyNmapDB{}
}

func (d *DummyNmapDB) Search(_ context.Context, _ *models.SearchParams) (*NmapSearchResult, error) {
	return &NmapSearchResult{
		Total:   0,
		Results: []models.NmapDocument{},
	}, nil
}

func (d *DummyNmapDB) Insert(_ context.Context, _ *models.NmapData) ([]string, error) {
	return []string{}, nil
}

func (d *DummyNmapDB) InsertBatch(_ context.Context, _ []models.NmapData) ([]string, error) {
	return []string{}, nil
}
