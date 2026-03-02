package repositories

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/models"
)

// SearchableRepository is a generic interface for any resource that supports search and pagination
type SearchableRepository[T any] interface {
	Search(ctx context.Context, params *models.SearchParams) (uint64, []T, error)
}
