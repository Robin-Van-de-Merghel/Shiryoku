package main

import (
	"fmt"
	"log"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/config"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/db"
	"github.com/gin-gonic/gin"
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
	provider, err := db.InitDB(dsn)
	if err != nil {
		log.Fatalf("couldn't initialize DB connection: %v", err)
	}

	// TODO: Import external modules

	// Pass to router
	engine := gin.New()
	if err := routers.SetupRoutes(engine, provider); err != nil {
		log.Fatalf("couldn't setup routes: %v", err)
	}

	err = engine.Run(
		fmt.Sprintf(":%d", serverConfig.Port),
	)

	if err != nil {
		// Rather kill instantly
		log.Fatalf("an error occurred with the server: %v", err)
	}
}
