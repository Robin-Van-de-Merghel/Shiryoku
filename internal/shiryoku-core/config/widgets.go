package config

import (
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	routers_utils_setup "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils/setup"
)

func GetDefaultWidgets(client osdb.OpenSearchClient) []APIModule {
	return []APIModule{
		{
			Name: "Dashboard",
			Description: "Dashboard shown on the first page",
			Path: "/dashboard",
			OSDB: client,
			SetupModuleRoutes: routers_utils_setup.SetupWidgetsDashboardRoutes,
		},
	}
}

