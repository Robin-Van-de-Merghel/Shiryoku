package osdb

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

// NmapDBIface is the interface that both NmapDB and DummyNmapDB implement.
type NmapDBIface interface {
	Search(ctx context.Context, params *models.SearchParams) (*NmapSearchResult, error)
	Insert(ctx context.Context, nmapData []models.NmapDocument) ([]string, error)
}

// DummyNmapDB is an in-memory implementation of NmapDBIface for unit tests.
type DummyNmapDB struct {
	store map[string]models.NmapDocument
}

func NewDummyNmapDB() NmapDBIface {
	return &DummyNmapDB{
		store: make(map[string]models.NmapDocument),
	}
}

func (d *DummyNmapDB) Search(_ context.Context, params *models.SearchParams) (*NmapSearchResult, error) {
	results := []models.NmapDocument{}
	for _, doc := range d.store {
		if matchesAll(doc, params.Search) {
			results = append(results, doc)
		}
	}
	return &NmapSearchResult{Total: uint64(len(results)), Results: results}, nil
}

func (d *DummyNmapDB) Insert(_ context.Context, docs []models.NmapDocument) ([]string, error) {
	ids := make([]string, 0, len(docs))
	for _, doc := range docs {
		id := fmt.Sprintf("%s:%d", doc.Host, doc.Port)
		d.store[id] = doc
		ids = append(ids, id)
	}
	return ids, nil
}

func matchesAll(doc models.NmapDocument, specs []models.SearchSpec) bool {
	for _, spec := range specs {
		if spec.Scalar != nil && !matchesScalar(doc, spec.Scalar) {
			return false
		}
	}
	return true
}

func matchesScalar(doc models.NmapDocument, spec *models.ScalarSearchSpec) bool {
	fieldVal := getDocField(doc, spec.Parameter)
	switch spec.Operator {
	case models.OpEq:
		return fmt.Sprintf("%v", fieldVal) == fmt.Sprintf("%v", spec.Value)
	case models.OpNeq:
		return fmt.Sprintf("%v", fieldVal) != fmt.Sprintf("%v", spec.Value)
	default:
		return true
	}
}

func getDocField(doc models.NmapDocument, field string) any {
	switch field {
	case "host":
		return doc.Host
	case "port":
		return doc.Port
	case "status":
		return doc.HostStatus
	case "service_name":
		return doc.ServiceName
	case "service_version":
		return doc.ServiceVersion
	default:
		return nil
	}
}
