package utils

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

// Generate schema from struct
// Useful for data validation
var SearchSchema = GenerateSchema(models.SearchParams{})
