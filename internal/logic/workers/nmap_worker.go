package workers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/logic/widgets/dashboard"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/config"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/db"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
)

// NmapWorker refreshes the dashboard materialized view with latest nmap scans
type NmapWorker struct {
	config   *config.WorkerConfig
	provider repositories.RepositoryProvider
	ticker   *time.Ticker
	done     chan bool
}

// NewNmapWorker creates a new nmap worker instance
func NewNmapWorker(workerConfig *config.WorkerConfig) (*NmapWorker, error) {
	// Initialize database
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		workerConfig.DBConfig.Username,
		workerConfig.DBConfig.Password,
		workerConfig.DBConfig.Host,
		workerConfig.DBConfig.Port,
		workerConfig.DBConfig.Database,
	)
	provider, err := db.InitDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &NmapWorker{
		config:   workerConfig,
		provider: provider,
		ticker:   time.NewTicker(workerConfig.Frequency),
		done:     make(chan bool),
	}, nil
}

// Start begins the worker's refresh loop
func (w *NmapWorker) Start(ctx context.Context) {
	log.Printf("[%s] Starting worker with frequency: %v", w.config.Name, w.config.Frequency)
	// Refresh immediately on start
	w.refreshView(ctx)
	// Then refresh on ticker
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.refreshView(ctx)
			case <-w.done:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop gracefully shuts down the worker
func (w *NmapWorker) Stop() {
	log.Printf("[%s] Stopping worker", w.config.Name)
	w.ticker.Stop()
	w.done <- true
}

// refreshView fetches latest scans and refreshes the materialized view
func (w *NmapWorker) refreshView(ctx context.Context) {
	start := time.Now()
	log.Printf("[%s] Refreshing dashboard table", w.config.Name)

	nmapRepo := w.provider.GetRepository("nmap").(repositories.NmapRepository)
	dashboardRepo := w.provider.GetRepository("dashboard").(postgres.DashboardRepository)

	if err := dashboard.BuildDashboardScans(ctx, nmapRepo, dashboardRepo, 7); err != nil {
		log.Printf("[%s] Error rebuilding dashboard: %v", w.config.Name, err)
		return
	}
	log.Printf("[%s] Dashboard table rebuilt in %v", w.config.Name, time.Since(start))
}
