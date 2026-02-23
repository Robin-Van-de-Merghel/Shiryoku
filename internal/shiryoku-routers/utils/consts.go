package utils

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
)

// Generate schema from struct
// Useful for data validation
var SearchSchema = GenerateSchema(models.SearchParams{})
var WidgetDashboardSchema = GenerateSchema(models_widgets.WidgetDashboardInput{})
