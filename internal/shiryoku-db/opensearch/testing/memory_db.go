package opensearch_testing

import (
	"context"
	"encoding/json"
	"sync"
)

// MemoryDB is a generic in-memory store that simulates OpenSearch insert/search.
// Thread-safe. Stores documents as raw map[string]any.
type MemoryDB struct {
	mu   sync.RWMutex
	docs map[string]map[string]any
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		docs: make(map[string]map[string]any),
	}
}

// Store inserts or replaces a document by ID.
func (m *MemoryDB) Store(_ context.Context, id string, document any) error {
	b, err := json.Marshal(document)
	if err != nil {
		return err
	}

	var raw map[string]any
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.docs[id] = raw
	return nil
}

// All returns all stored documents as a slice of map[string]any.
func (m *MemoryDB) All() []map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]map[string]any, 0, len(m.docs))
	for _, doc := range m.docs {
		results = append(results, doc)
	}
	return results
}

// Len returns the number of stored documents.
func (m *MemoryDB) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.docs)
}

// Reset clears all stored documents.
func (m *MemoryDB) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.docs = make(map[string]map[string]any)
}
