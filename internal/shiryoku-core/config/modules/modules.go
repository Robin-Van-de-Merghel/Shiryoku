package modules

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/modules/nmap"
)

func GetDefaultModules(client osdb.OpenSearchClient) []config.APIModule {
	var modules []config.APIModule

	modules = append(modules, config.APIModule{
		Name:              "NMap Module",
		Description:       "Use to parse nmap results and query them.",
		Path:              "/nmap",
		OSDB:              client,
		SetupModuleRoutes: nmap.SetupNmapRoutes,
	})

	return modules
}
