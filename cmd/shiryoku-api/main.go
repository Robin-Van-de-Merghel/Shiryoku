package main

import (
	"fmt"
	"log"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	shiryoku_routers "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers"
	"github.com/opensearch-project/opensearch-go"
)

func main() {
	serverConfig := config.NewServerConfig()

	// Use struct field instead of map?
	osdbConfig := serverConfig.DBConfigs["OSDB"]

	// Create OpenSearch client
	osClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{
			// TODO: Use password and schema
			fmt.Sprintf("http://%s:%d", osdbConfig.Host, osdbConfig.Port),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Wrap in NmapDB
	osdbClient := osdb.NewOpenSearchClient(osClient)


	// Get the modules and widgets
	serverConfig.Modules = config.GetDefaultModules(*osdbClient)
	serverConfig.Widgets = config.GetDefaultWidgets(*osdbClient)
	// TODO: Import external ones

	// Pass to router
	router := shiryoku_routers.GetFilledRouter(*serverConfig)

	// FIXME: port from config
	err = router.Run(":8080")

	if err != nil {
		// Rather kill instantly
		panic(err)
	}
}
