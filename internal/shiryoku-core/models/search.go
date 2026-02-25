package models

import (
	"encoding/json"
	"fmt"
)

// Consts
const MAX_RESULTS_PER_PAGE = 1000
const DEFAULT_RESULTS_PER_PAGE = 100

// ScalarOperator enum
type ScalarOperator string

const (
	OpEq      ScalarOperator = "eq"
	OpNeq     ScalarOperator = "neq"
	OpGt      ScalarOperator = "gt"
	OpLt      ScalarOperator = "lt"
	OpLike    ScalarOperator = "like"
	OpNotLike ScalarOperator = "not like"
	OpRegex   ScalarOperator = "regex"
)

func (s ScalarOperator) IsValid() bool {
	switch s {
	case OpEq, OpNeq, OpGt, OpLt, OpLike, OpNotLike, OpRegex:
		return true
	default:
		return false
	}
}

// VectorOperator enum
type VectorOperator string

const (
	OpIn    VectorOperator = "in"
	OpNotIn VectorOperator = "not in"
)

func (v VectorOperator) IsValid() bool {
	switch v {
	case OpIn, OpNotIn:
		return true
	default:
		return false
	}
}

// SortDirection enum
type SortDirection string

const (
	DirASC  SortDirection = "asc"
	DirDESC SortDirection = "desc"
)

func (sd SortDirection) IsValid() bool {
	switch sd {
	case DirASC, DirDESC:
		return true
	default:
		return false
	}
}

type SortSpec struct {
	Parameter string        `json:"parameter"`
	Direction SortDirection `json:"direction"`
}

// To fetch couple elements based on their values
type ScalarSearchSpec struct {
	Parameter string         `json:"parameter" validate:"required"`
	Operator  ScalarOperator `json:"operator" validate:"required"`
	Value     any            `json:"value" validate:"required"`
}

// Fetch with "IN" clauses
type VectorSearchSpec struct {
	Parameter string         `json:"parameter" validate:"required"`
	Operator  VectorOperator `json:"operator" validate:"required"`
	Values    any            `json:"values" validate:"required"`
}

// SearchSpec either Vector, either Scalar
type SearchSpec struct {
	Scalar *ScalarSearchSpec
	Vector *VectorSearchSpec
}

type SearchParams struct {
	Parameters []string     `json:"parameters,omitempty"`
	Search     []SearchSpec `json:"search" validate:"dive,required"`
	Sort       []SortSpec   `json:"sort"`
	Distinct   bool         `json:"distinct"`
	Page       uint64       `json:"page"`
	PerPage    uint64       `json:"per_page"`
}

func (s *SearchParams) SetDefaults() {
	if s.PerPage == 0 {
		s.PerPage = DEFAULT_RESULTS_PER_PAGE
	}

	if s.Sort == nil {
		s.Sort = []SortSpec{}
	}
}

// Custom unmarshaler for SearchSpec
func (s *SearchSpec) UnmarshalJSON(data []byte) error {
	// First pass: read the operator to discriminate
	var base struct {
		Operator string `json:"operator"`
	}
	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	// Second pass: unmarshal into correct type based on operator
	switch VectorOperator(base.Operator) {
	case OpIn, OpNotIn:
		var vector VectorSearchSpec
		if err := json.Unmarshal(data, &vector); err != nil {
			return err
		}
		if !vector.Operator.IsValid() {
			return fmt.Errorf("invalid vector operator: %s", vector.Operator)
		}
		s.Vector = &vector
	default:
		// Try scalar
		var scalar ScalarSearchSpec
		if err := json.Unmarshal(data, &scalar); err != nil {
			return err
		}
		if !scalar.Operator.IsValid() {
			return fmt.Errorf("invalid scalar operator: %s", scalar.Operator)
		}
		s.Scalar = &scalar
	}

	return nil
}
