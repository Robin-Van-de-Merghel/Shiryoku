package config

import (
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	routers_utils_setup "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils/setup"
)

func GetDefaultModules(client osdb.OpenSearchClient) []APIModule {
	return []APIModule{
		{
			Name:              "NMap Module",
			Description:       "Use to parse nmap results and query them.",
			Path:              "/nmap",
			OSDB:              client,
			SetupModuleRoutes: routers_utils_setup.SetupNmapRoutes,
		},
	}
}
