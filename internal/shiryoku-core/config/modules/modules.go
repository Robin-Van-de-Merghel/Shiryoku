package modules

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/modules/nmap"
	routers_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/widgets"
)

func GetDefaultModules(client osdb.OpenSearchClient) []config.APIModule {
	return []config.APIModule{
		{
			Name:              "NMap Module",
			Description:       "Use to parse nmap results and query them.",
			Path:              "/nmap",
			OSDB:              client,
			SetupModuleRoutes: nmap.SetupNmapRoutes,
		},
		{
			// TODO: See if we separate also here
			Name: "Widgets",
			Description: "All the widgets used by the UI",
			Path: "/widgets",
			OSDB: client,
			SetupModuleRoutes: routers_widgets.SetupWidgetsRoutes,
		},
	}
}
