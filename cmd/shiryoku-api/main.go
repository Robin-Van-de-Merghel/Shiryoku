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
	nmapDB := osdb.NewNmapDB(osdb.NewOpenSearchClient(osClient))

	// Get the modules
	default_modules := modules.GetDefaultModules(nmapDB)
	// TODO: Import external ones

	// Pass to router
	router := shiryoku_routers.GetFilledRouter(nmapDB, default_modules)

	// FIXME: port from config
	router.Run(":8080")
}
