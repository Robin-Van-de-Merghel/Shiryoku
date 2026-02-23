package osdb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
)

const nmapIndex = "nmap-scans"

type NmapSearchResult struct {
	Total   uint64                `json:"total"`
	Results []models.NmapDocument `json:"results"`
}

type NmapDB struct {
	osClient *OpenSearchClient
}

func NewNmapDB(osClient *OpenSearchClient) *NmapDB {
	return &NmapDB{osClient: osClient}
}

// Search returns flat NmapDocument results
func (db *NmapDB) Search(ctx context.Context, params *models.SearchParams) (*NmapSearchResult, error) {
	searchResult, err := db.osClient.Search(ctx, nmapIndex, params)
	if err != nil {
		return nil, err
	}

	results := make([]models.NmapDocument, 0, len(searchResult.Results))
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

// Insert explodes NmapData into one document per port and bulk-inserts them
func (db *NmapDB) Insert(ctx context.Context, docs []models.NmapDocument) ([]string, error) {
	ids := make([]string, 0, len(docs))

	for _, doc := range docs {

		// Verify that host is defined
		if doc.Host == "" {
			// TODO: Logs
			continue
		}

		// ID = host:port to allow upsert
		id := fmt.Sprintf("%s:%d", doc.Host, doc.Port)
		result, err := db.osClient.Insert(ctx, nmapIndex, id, doc)
		if err != nil {
			return ids, fmt.Errorf("failed to insert document for port %d: %w", doc.Port, err)
		}
		ids = append(ids, result.ID)
	}

	return ids, nil
}

func parseNmapDocument(rawData map[string]any) (models.NmapDocument, error) {
	// Re-marshal and unmarshal into the struct — clean and safe
	b, err := json.Marshal(rawData)
	if err != nil {
		return models.NmapDocument{}, fmt.Errorf("failed to marshal raw data: %w", err)
	}

	var doc models.NmapDocument
	if err := json.Unmarshal(b, &doc); err != nil {
		return models.NmapDocument{}, fmt.Errorf("failed to unmarshal nmap document: %w", err)
	}

	return doc, nil
}
