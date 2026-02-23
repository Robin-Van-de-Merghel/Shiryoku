package main

import (
	"log"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config/modules"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	shiryoku_routers "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers"
	"github.com/opensearch-project/opensearch-go"
)

func main() {
	// Create OpenSearch client
	osClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"http://localhost:9200"}, // Or from config
	})
	if err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Wrap in NmapDB
	osdbClient := osdb.NewOpenSearchClient(osClient)

	// Get the modules
	default_modules := modules.GetDefaultModules(*osdbClient)
	// TODO: Import external ones

	// Pass to router
	router := shiryoku_routers.GetFilledRouter(osdbClient, default_modules)

	// FIXME: port from config
	err = router.Run(":8080")

	if err != nil {
		// Rather kill instantly
		panic(err)
	}
}
