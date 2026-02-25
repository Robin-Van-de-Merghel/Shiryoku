package common

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
)

type SearchResult[T any] struct {
	Total   uint64 `json:"total"`
	Results []T    `json:"results"`
}

// Search is a generic wrapper that works with any SearchableRepository[T]
func Search[T any](ctx context.Context, repo postgres.SearchableRepository[T], params *models.SearchParams) (*SearchResult[T], error) {
	total, results, err := repo.Search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	return &SearchResult[T]{
		Total:   total,
		Results: results,
	}, nil
}
