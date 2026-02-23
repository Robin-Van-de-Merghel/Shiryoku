package opensearch_testing

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
)

// DummyNmapDB implements osdb.NmapDBIface using a MemoryDB.
// Insert stores flat NmapDocuments. Search returns everything stored.
type DummyNmapDB struct {
	mem *MemoryDB
}

func NewDummyNmapDB() *DummyNmapDB {
	return &DummyNmapDB{mem: NewMemoryDB()}
}

// Insert explodes NmapData into flat NmapDocuments and stores them.
func (d *DummyNmapDB) Insert(ctx context.Context, docs []models.NmapDocument) ([]string, error) {
	ids := make([]string, 0, len(docs))
	for _, doc := range docs {
		id := fmt.Sprintf("%s:%d", doc.Host, doc.Port)
		if err := d.mem.Store(ctx, id, doc); err != nil {
			return ids, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// Search returns all stored documents as NmapDocuments (no filtering — it's a dummy).
func (d *DummyNmapDB) Search(_ context.Context, _ *models.SearchParams) (*osdb.NmapSearchResult, error) {
	raw := d.mem.All()

	results := make([]models.NmapDocument, 0, len(raw))
	for _, r := range raw {
		doc, err := toNmapDocument(r)
		if err != nil {
			continue
		}
		results = append(results, doc)
	}

	return &osdb.NmapSearchResult{
		Total:   uint64(len(results)),
		Results: results,
	}, nil
}

// Mem exposes the underlying MemoryDB for assertions in tests.
func (d *DummyNmapDB) Mem() *MemoryDB {
	return d.mem
}

func toNmapDocument(raw map[string]any) (models.NmapDocument, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return models.NmapDocument{}, err
	}
	var doc models.NmapDocument
	if err := json.Unmarshal(b, &doc); err != nil {
		return models.NmapDocument{}, err
	}
	return doc, nil
}
