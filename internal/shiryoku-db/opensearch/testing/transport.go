package opensearch_testing

// TODO: This file has been exceptionnally AI-generated (via claude)
// For later, check for bugs, and improve.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/opensearch-project/opensearch-go"
)

// Storage

type storedDoc struct {
	source  map[string]any
	version uint64
}

// MockTransport implements http.RoundTripper.
// It intercepts all calls made by opensearch-go and responds with in-memory data.
// Only the two endpoints actually used are implemented: _bulk and _search.
// TODO: Later implement more (deletion)

type MockTransport struct {
	mu   sync.RWMutex
	docs map[string]*storedDoc // key: "index/id"
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		docs: make(map[string]*storedDoc),
	}
}

// NewMockOpenSearchClient returns a real *opensearch.Client backed by MockTransport.
// Use it exactly like opensearch.NewClient(...) — just pass the result to osdb.NewOpenSearchClient.
// Helps keeping the logic 
func NewMockOpenSearchClient(t *MockTransport) (*opensearch.Client, error) {
	return opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://localhost:9200"},
		Transport: t,
	})
}

// RoundTrip handles requests
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	path := req.URL.Path

	switch {
	case req.Method == http.MethodPost && path == "/_bulk":
		return m.handleBulk(req)

	case req.Method == http.MethodPost && strings.HasSuffix(path, "/_search"):
		return m.handleSearch(req)

	default:
		// Any other call (e.g. HEAD / for ping) → acknowledge it
		// TODO: Should it be an error for other..? If not implemented
		return okResponse(map[string]any{"acknowledged": true})
	}
}

// handleBulk parses the NDJSON bulk body and upserts documents.
// Only the "index" action is supported (which is all InsertBulk uses).
//
// Format (two lines per doc):
//
//	{ "index": { "_index": "...", "_id": "..." } }
//	{ ...document... }
func (m *MockTransport) handleBulk(req *http.Request) (*http.Response, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return errResponse(http.StatusBadRequest, "cannot read body")
	}

	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines)%2 != 0 {
		return errResponse(http.StatusBadRequest, "malformed bulk body: odd number of lines")
	}

	type bulkResponseItem struct {
		Index struct {
			ID      string `json:"_id"`
			Index   string `json:"_index"`
			Version uint64 `json:"_version"`
			Result  string `json:"result"`
			Status  int    `json:"status"`
		} `json:"index"`
	}

	items := make([]bulkResponseItem, 0, len(lines)/2)

	for i := 0; i < len(lines); i += 2 {
		// Parse action line
		var action struct {
			Index struct {
				Index string `json:"_index"`
				ID    string `json:"_id"`
			} `json:"index"`
		}
		if err := json.Unmarshal([]byte(lines[i]), &action); err != nil {
			return errResponse(http.StatusBadRequest, fmt.Sprintf("invalid action line %d: %v", i, err))
		}

		// Parse document line
		var source map[string]any
		if err := json.Unmarshal([]byte(lines[i+1]), &source); err != nil {
			return errResponse(http.StatusBadRequest, fmt.Sprintf("invalid document line %d: %v", i+1, err))
		}

		index := action.Index.Index
		id := action.Index.ID
		key := index + "/" + id

		m.mu.Lock()
		existing, exists := m.docs[key]
		version := uint64(1)
		if exists {
			version = existing.version + 1
		}
		m.docs[key] = &storedDoc{source: source, version: version}
		m.mu.Unlock()

		result := "created"
		if exists {
			result = "updated"
		}

		var item bulkResponseItem
		item.Index.ID = id
		item.Index.Index = index
		item.Index.Version = version
		item.Index.Result = result
		item.Index.Status = http.StatusOK
		items = append(items, item)
	}

	return okResponse(map[string]any{
		"errors": false,
		"items":  items,
	})
}

// handleSearch returns all documents stored in the requested index.
// Filtering is intentionally not implemented — this is a mock.
// If you need filter testing, use DummyNmapDB (higher level) instead.
func (m *MockTransport) handleSearch(req *http.Request) (*http.Response, error) {
	// Extract index from path: /index/_search
	index := strings.Split(strings.Trim(req.URL.Path, "/"), "/")[0]

	m.mu.RLock()
	defer m.mu.RUnlock()

	type hit struct {
		Index  string         `json:"_index"`
		ID     string         `json:"_id"`
		Source map[string]any `json:"_source"`
	}

	hits := make([]hit, 0)
	for key, doc := range m.docs {
		parts := strings.SplitN(key, "/", 2)
		if parts[0] != index {
			continue
		}
		hits = append(hits, hit{
			Index:  parts[0],
			ID:     parts[1],
			Source: doc.source,
		})
	}

	return okResponse(map[string]any{
		"hits": map[string]any{
			"total": map[string]any{
				"value":    len(hits),
				"relation": "eq",
			},
			"hits": hits,
		},
	})
}

// Len returns the total number of stored documents across all indices.
func (m *MockTransport) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.docs)
}

// LenIndex returns the number of documents stored in a specific index.
func (m *MockTransport) LenIndex(index string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for key := range m.docs {
		if strings.HasPrefix(key, index+"/") {
			count++
		}
	}
	return count
}

// Reset clears all stored documents.
func (m *MockTransport) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.docs = make(map[string]*storedDoc)
}

func okResponse(body any) (*http.Response, error) {
	return jsonResponse(http.StatusOK, body)
}

func errResponse(status int, msg string) (*http.Response, error) {
	return jsonResponse(status, map[string]any{"error": msg})
}

func jsonResponse(status int, body any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("mock: failed to marshal response: %w", err)
	}
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(b)),
	}, nil
}
