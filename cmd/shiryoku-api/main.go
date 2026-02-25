package main

import (
	"fmt"
	"log"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	shiryoku_routers "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers"
)

func main() {
	serverConfig := config.NewServerConfig()

	// URL to the db
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		serverConfig.DBConfig.Username,
		serverConfig.DBConfig.Password,
		serverConfig.DBConfig.Host,
		serverConfig.DBConfig.Port,
		serverConfig.DBConfig.Database,
	)

	// Create OpenSearch client
	repos, err := shiryoku_db.InitDB(dsn)
	if err != nil {
		log.Fatalf("couldn't initialize DB connection: %v", err)
	}

	// Get the modules and widgets
	serverConfig.Modules = config.GetDefaultModules(repos)
	serverConfig.Widgets = config.GetDefaultWidgets(repos)
	// TODO: Import external ones

	// Pass to router
	router := shiryoku_routers.GetFilledRouter(*serverConfig, repos)

	err = router.Run(
		fmt.Sprintf(":%d", serverConfig.Port),
	)

	if err != nil {
		// Rather kill instantly
		log.Fatalf("an error occurred with the server: %v", err)
	}
}
