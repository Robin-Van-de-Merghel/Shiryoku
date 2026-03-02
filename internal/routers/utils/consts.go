package utils

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/models"
)

// Generate schema from struct
// Useful for data validation
var SearchSchema = GenerateSchema(models.SearchParams{})

// Precompute field map
// Use for Search's parameters: verify that search: {"parameter": "azazazazaza"} exists
var NmapScanFields = buildFieldTypeMap(models.NmapScan{})
var NmapHostFields = buildFieldTypeMap(models.NmapHost{})
var NmapScriptResultFields = buildFieldTypeMap(models.NmapScriptResult{})
var WidgetDashboardScanFields = buildFieldTypeMap(widgets.WidgetDashboardScan{})
