package db

import (
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
)

func InitDB(dsn string) (*repositories.DefaultRepositoryProvider, error) {
	db, err := postgres.NewPostgresDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	provider := repositories.NewRepositoryProvider()
	provider.RegisterRepository(repositories.NMAP_REPOSITORY, postgres.NewNmapRepository(db))
	// TODO: See if we call it from init (as it's internal)
	provider.RegisterRepository(repositories.DASHBOARD_REPOSITORY, postgres.NewDashboardRepository(db))

	return provider, nil
}
