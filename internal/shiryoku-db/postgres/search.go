package postgres

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"gorm.io/gorm"
)

type SearchBuilder[T any] struct {
	db *gorm.DB
}

func NewSearchBuilder[T any](db *gorm.DB) *SearchBuilder[T] {
	return &SearchBuilder[T]{db: db}
}

func (sb *SearchBuilder[T]) Build(params *models.SearchParams) (*gorm.DB, error) {
	query := sb.db

	for _, spec := range params.Search {
		if spec.Scalar != nil {
			var err error
			query, err = sb.applyScalarFilter(query, spec.Scalar)
			if err != nil {
				return nil, err
			}
		}
		if spec.Vector != nil {
			var err error
			query, err = sb.applyVectorFilter(query, spec.Vector)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, sort := range params.Sort {
		query = query.Order(fmt.Sprintf("\"%s\" %s", sort.Parameter, sort.Direction))
	}

	return query, nil
}

func (sb *SearchBuilder[T]) applyScalarFilter(query *gorm.DB, spec *models.ScalarSearchSpec) (*gorm.DB, error) {
	switch spec.Operator {
	case models.OpEq:
		return query.Where(fmt.Sprintf("\"%s\" = ?", spec.Parameter), spec.Value), nil
	case models.OpNeq:
		return query.Where(fmt.Sprintf("\"%s\" != ?", spec.Parameter), spec.Value), nil
	case models.OpGt:
		return query.Where(fmt.Sprintf("\"%s\" > ?", spec.Parameter), spec.Value), nil
	case models.OpLt:
		return query.Where(fmt.Sprintf("\"%s\" < ?", spec.Parameter), spec.Value), nil
	case models.OpLike:
		escapedValue := strings.ReplaceAll(fmt.Sprint(spec.Value), "%", "\\%")
		escapedValue = strings.ReplaceAll(escapedValue, "_", "\\_")
		return query.Where(fmt.Sprintf("\"%s\" ILIKE ?", spec.Parameter), "%"+escapedValue+"%"), nil
	case models.OpNotLike:
		escapedValue := strings.ReplaceAll(fmt.Sprint(spec.Value), "%", "\\%")
		escapedValue = strings.ReplaceAll(escapedValue, "_", "\\_")
		return query.Where(fmt.Sprintf("\"%s\" NOT ILIKE ?", spec.Parameter), "%"+escapedValue+"%"), nil
	case models.OpRegex:
		if _, err := regexp.Compile(fmt.Sprint(spec.Value)); err != nil {
			return nil, fmt.Errorf("invalid regex: %w", err)
		}
		return query.Where(fmt.Sprintf("\"%s\" ~ ?", spec.Parameter), spec.Value), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", spec.Operator)
	}
}

func (sb *SearchBuilder[T]) applyVectorFilter(query *gorm.DB, spec *models.VectorSearchSpec) (*gorm.DB, error) {
	switch spec.Operator {
	case models.OpIn:
		return query.Where(fmt.Sprintf("\"%s\" IN ?", spec.Parameter), spec.Values), nil
	case models.OpNotIn:
		return query.Where(fmt.Sprintf("\"%s\" NOT IN ?", spec.Parameter), spec.Values), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", spec.Operator)
	}
}
