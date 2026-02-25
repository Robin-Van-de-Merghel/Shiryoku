package shiryoku_db

import (
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
)

// Repositories holds all database repositories for different data sources
type Repositories struct {
	// Nmap scan data repository
	Nmap postgres.NmapRepository
}

// InitDB initializes the database connection and returns all configured repositories
func InitDB(dsn string) (*Repositories, error) {
	db, err := postgres.NewDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Repositories{
		Nmap: postgres.NewNmapRepository(db),
	}, nil
}
