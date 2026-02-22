package osdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// SearchResult wraps raw results with metadata
type SearchResult struct {
	Total   uint64           `json:"total"`
	Results []map[string]any `json:"results"`
}

// OpenSearchClient wraps the opensearch client for generic operations
type OpenSearchClient struct {
	client *opensearch.Client
}

// NewOpenSearchClient creates a new generic OpenSearch client
func NewOpenSearchClient(client *opensearch.Client) *OpenSearchClient {
	return &OpenSearchClient{client: client}
}

// Search performs a search on a given index
func (os *OpenSearchClient) Search(ctx context.Context, index string, params *models.SearchParams) (*SearchResult, error) {
	query := buildOpenSearchQuery(params)

	body := map[string]any{
		"query": query,
		"size":  params.PerPage,
	}

	if params.Page > 0 {
		body["from"] = (params.Page - 1) * uint64(params.PerPage)
	}

	if len(params.Parameters) > 0 {
		body["_source"] = params.Parameters
	}

	if len(params.Sort) > 0 {
		sorts := make([]map[string]any, len(params.Sort))
		for i, sort := range params.Sort {
			sorts[i] = map[string]any{
				sort.Parameter: map[string]any{
					"order": sort.Direction,
				},
			}
		}
		body["sort"] = sorts
	}

	bodyBytes, _ := json.Marshal(body)

	req := opensearchapi.SearchRequest{
		Index: []string{index},
		Body:  strings.NewReader(string(bodyBytes)),
	}

	res, err := req.Do(ctx, os.client)
	if err != nil {
		return nil, fmt.Errorf("opensearch search failed: %w", err)
	}
	defer res.Body.Close()

	var searchResponse struct {
		Hits struct {
			Total struct {
				Value uint64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	results := make([]map[string]any, len(searchResponse.Hits.Hits))
	for i, hit := range searchResponse.Hits.Hits {
		results[i] = hit.Source
	}

	if params.Distinct {
		results = deduplicateGeneric(results)
	}

	return &SearchResult{
		Total:   searchResponse.Hits.Total.Value,
		Results: results,
	}, nil
}

// buildOpenSearchQuery converts SearchParams into an OpenSearch bool query
// No nested logic needed — documents are flat
func buildOpenSearchQuery(params *models.SearchParams) map[string]any {
	if len(params.Search) == 0 {
		return map[string]any{
			"match_all": map[string]any{},
		}
	}

	must := []map[string]any{}
	mustNot := []map[string]any{}

	for _, spec := range params.Search {
		if spec.Scalar != nil {
			buildScalarQuery(spec.Scalar, &must, &mustNot)
		} else if spec.Vector != nil {
			buildVectorQuery(spec.Vector, &must, &mustNot)
		}
	}

	return map[string]any{
		"bool": map[string]any{
			"must":     must,
			"must_not": mustNot,
		},
	}
}

func buildScalarQuery(spec *models.ScalarSearchSpec, must, mustNot *[]map[string]any) {
	field := spec.Parameter
	value := spec.Value

	// For keyword/string fields, use .keyword suffix for exact match
	// For known numeric fields, use the field as-is
	fieldForTerm := toTermField(field)

	switch spec.Operator {
	case models.OpEq:
		*must = append(*must, map[string]any{
			"term": map[string]any{fieldForTerm: value},
		})
	case models.OpNeq:
		*mustNot = append(*mustNot, map[string]any{
			"term": map[string]any{fieldForTerm: value},
		})
	case models.OpGt:
		*must = append(*must, map[string]any{
			"range": map[string]any{field: map[string]any{"gt": value}},
		})
	case models.OpLt:
		*must = append(*must, map[string]any{
			"range": map[string]any{field: map[string]any{"lt": value}},
		})
	case models.OpLike:
		*must = append(*must, map[string]any{
			"wildcard": map[string]any{field: map[string]any{"value": fmt.Sprintf("*%v*", value)}},
		})
	case models.OpNotLike:
		*mustNot = append(*mustNot, map[string]any{
			"wildcard": map[string]any{field: map[string]any{"value": fmt.Sprintf("*%v*", value)}},
		})
	case models.OpRegex:
		*must = append(*must, map[string]any{
			"regexp": map[string]any{field: map[string]any{"value": value}},
		})
	}
}

func buildVectorQuery(spec *models.VectorSearchSpec, must, mustNot *[]map[string]any) {
	field := spec.Parameter
	values := spec.Values

	var valueSlice []any
	switch v := values.(type) {
	case []any:
		valueSlice = v
	default:
		valueSlice = []any{values}
	}

	switch spec.Operator {
	case models.OpIn:
		*must = append(*must, map[string]any{
			"terms": map[string]any{field: valueSlice},
		})
	case models.OpNotIn:
		*mustNot = append(*mustNot, map[string]any{
			"terms": map[string]any{field: valueSlice},
		})
	}
}

// toTermField returns the field name to use for term queries.
// Numeric fields don't need .keyword, string fields do (for exact match on dynamic mappings).
func toTermField(field string) string {
	numericFields := map[string]bool{
		"port": true,
	}
	// Get last part of dotted path
	parts := strings.Split(field, ".")
	last := parts[len(parts)-1]

	if numericFields[last] {
		return field
	}
	return field + ".keyword"
}

func deduplicateGeneric(results []map[string]any) []map[string]any {
	seen := make(map[string]bool)
	dedup := []map[string]any{}

	for _, result := range results {
		key, _ := json.Marshal(result)
		keyStr := string(key)
		if !seen[keyStr] {
			seen[keyStr] = true
			dedup = append(dedup, result)
		}
	}

	return dedup
}

// InsertResult wraps the insert response
type InsertResult struct {
	ID      string `json:"id"`
	Index   string `json:"index"`
	Version uint64 `json:"version"`
}

// Insert performs a generic insert/upsert and returns the document ID
func (os *OpenSearchClient) Insert(ctx context.Context, index string, id string, document any) (*InsertResult, error) {
	bodyBytes, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal document: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       strings.NewReader(string(bodyBytes)),
	}

	res, err := req.Do(ctx, os.client)
	if err != nil {
		return nil, fmt.Errorf("opensearch insert failed: %w", err)
	}
	defer res.Body.Close()

	var indexResponse struct {
		ID      string `json:"_id"`
		Index   string `json:"_index"`
		Version uint64 `json:"_version"`
	}

	if err := json.NewDecoder(res.Body).Decode(&indexResponse); err != nil {
		return nil, fmt.Errorf("failed to decode insert response: %w", err)
	}

	return &InsertResult{
		ID:      indexResponse.ID,
		Index:   indexResponse.Index,
		Version: indexResponse.Version,
	}, nil
}

// Upsert performs an insert or update
func (os *OpenSearchClient) Upsert(ctx context.Context, index string, id string, document any) (*InsertResult, error) {
	return os.Insert(ctx, index, id, document)
}
