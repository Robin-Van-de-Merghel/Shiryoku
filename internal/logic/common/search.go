package common

import (
	"context"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
)

// Search is a generic wrapper that works with any SearchableRepository[T]
func Search[T any](ctx context.Context, repo repositories.SearchableRepository[T], params *models.SearchParams) (*models.SearchResult[T], error) {
	total, results, err := repo.Search(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}
	return &models.SearchResult[T]{
		Total:   total,
		Results: results,
	}, nil
}
