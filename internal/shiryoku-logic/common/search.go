package logic_common

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
)

type SearchResult[T any] struct {
	Total   uint64                    `json:"total"`
	Results []T `json:"results"`
}


// Search returns results matching the given params (and Type)
func Search[T any](ctx context.Context, client osdb.OpenSearchClient, params *models.SearchParams, index string) (*SearchResult[T], error) {
	searchResult, err := client.Search(ctx, index, params)
	if err != nil {
		return nil, err
	}

	results := make([]T, 0, len(searchResult.Results))
	for _, rawResult := range searchResult.Results {
		doc, err := parseDocument[T](rawResult)
		if err != nil {
			fmt.Printf("failed to parse result: %v\n", err)
			continue
		}
		results = append(results, *doc)
	}

	return &SearchResult[T]{
		Total:   searchResult.Total,
		Results: results,
	}, nil
}

func parseDocument[T any](rawData map[string]any) (*T, error) {
	// Re-marshal and unmarshal into the struct — clean and safe
	b, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal raw data: %w", err)
	}

	var doc T 
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal document: %w", err)
	}

	return &doc, nil
}

