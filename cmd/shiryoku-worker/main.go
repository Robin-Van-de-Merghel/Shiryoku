package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/workers"
)

func main() {
	// Load worker config from environment
	workerConfig := config.NewWorkerConfig()

	log.Printf("Starting %s", workerConfig.Name)
	log.Printf("Refresh frequency: %v", workerConfig.Frequency)
	log.Printf("Database: %s:%d/%s", workerConfig.DBConfig.Host, workerConfig.DBConfig.Port, workerConfig.DBConfig.Database)

	// Create worker
	worker, err := workers.NewNmapWorker(workerConfig)
	if err != nil {
		log.Fatalf("Failed to create worker: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker
	worker.Start(ctx)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Shutdown signal received")

	worker.Stop()
	log.Println("Worker stopped")
}
