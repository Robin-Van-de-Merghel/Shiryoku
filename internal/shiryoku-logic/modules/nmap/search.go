package logic_nmap

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
)

const nmapIndex = "nmap-scans"

type NmapSearchResult struct {
	Total   uint64                `json:"total"`
	Results []models.NmapScanDocument `json:"results"`
}


// SearchNmapScans returns nmap scan results matching the given params
func Search(ctx context.Context, NmapDB osdb.OpenSearchClient, params *models.SearchParams) (*NmapSearchResult, error) {
	searchResult, err := NmapDB.Search(ctx, nmapIndex, params)
	if err != nil {
		return nil, err
	}

	results := make([]models.NmapScanDocument, 0, len(searchResult.Results))
	for _, rawResult := range searchResult.Results {
		doc, err := parseNmapDocument(rawResult)
		if err != nil {
			fmt.Printf("failed to parse nmap result: %v\n", err)
			continue
		}
		results = append(results, doc)
	}

	return &NmapSearchResult{
		Total:   searchResult.Total,
		Results: results,
	}, nil
}

func parseNmapDocument(rawData map[string]any) (models.NmapScanDocument, error) {
	// Re-marshal and unmarshal into the struct — clean and safe
	b, err := json.Marshal(rawData)
	if err != nil {
		return models.NmapScanDocument{}, fmt.Errorf("failed to marshal raw data: %w", err)
	}

	var doc models.NmapScanDocument
	if err := json.Unmarshal(b, &doc); err != nil {
		return models.NmapScanDocument{}, fmt.Errorf("failed to unmarshal nmap document: %w", err)
	}

	return doc, nil
}
