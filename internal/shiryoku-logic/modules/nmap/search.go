package nmap

import (
	"context"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
)

// SearchNmapScans returns nmap scan results matching the given params
func SearchNmapScans(ctx context.Context, searchParams *models.SearchParams, nmapDB osdb.NmapDBIface) (*osdb.NmapSearchResult, error) {
	return nmapDB.Search(ctx, searchParams)
}
