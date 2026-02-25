package routers_modules_nmap

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	postgres_testing "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres/testing"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(mockRepo *postgres_testing.MockNmapRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(utils.ErrorRecoveryMiddleware())

	api_group := r.Group("/api")
	{
		modules_group := api_group.Group("/modules")
		{
			nmap_group := modules_group.Group("/nmap")
			nmap_group.POST("/search", SearchNmapScans(mockRepo))
		}
	}
	return r
}

// Helper to make requests
func makeRequest(router *gin.Engine, payload any) *httptest.ResponseRecorder {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/modules/nmap/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestSearchNmapScansSuccess(t *testing.T) {
	mockRepo := &postgres_testing.MockNmapRepository{
		SearchFn: func(ctx context.Context, params *models.SearchParams) (uint64, []models.NmapScan, error) {
			return 0, []models.NmapScan{}, nil
		},
	}
	router := setupRouter(mockRepo)

	tests := []struct {
		name           string
		payload        string // Raw JSON strings
		expectedStatus int
		description    string
	}{
		{
			name: "Valid scalar search - eq operator",
			payload: `{
				"parameters": ["hostname", "port"],
				"search": [
					{
						"parameter": "status",
						"operator": "eq",
						"value": "open"
					}
				],
				"sort": [
					{
						"parameter": "hostname",
						"direction": "asc"
					}
				],
				"distinct": false
			}`,
			expectedStatus: 200,
			description:    "Should accept valid scalar search with eq operator",
		},
		{
			name: "Valid vector search - in operator",
			payload: `{
				"parameters": ["hostname"],
				"search": [
					{
						"parameter": "port",
						"operator": "in",
						"values": [80, 443, 8080]
					}
				]
			}`,
			expectedStatus: 200,
			description:    "Should accept valid vector search with in operator",
		},
		{
			name: "Mixed scalar and vector search",
			payload: `{
				"search": [
					{
						"parameter": "hostname",
						"operator": "like",
						"value": "%.example.com"
					},
					{
						"parameter": "protocol",
						"operator": "not in",
						"values": ["icmp", "igmp"]
					}
				]
			}`,
			expectedStatus: 200,
			description:    "Should accept mixed scalar and vector operators",
		},
		{
			name: "Regex operator",
			payload: `{
				"search": [
					{
						"parameter": "hostname",
						"operator": "regex",
						"value": "^[a-z0-9]+(\\.[a-z0-9]+)*\\.com$"
					}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle regex patterns",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/modules/nmap/search", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)
		})
	}
}

func TestSearchNmapScansValidationErrors(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		description    string
	}{
		{
			name:           "Empty JSON body",
			payload:        "{}",
			expectedStatus: 200, // Empty search is allowed (no filters)
			description:    "Should accept empty search (returns all results)",
		},
		{
			name:           "Invalid JSON format",
			payload:        `{"search": [{"operator": "eq"`,
			expectedStatus: 400,
			description:    "Should return 400 for malformed JSON",
		},
		{
			name:           "Missing required field - parameter",
			payload:        `{"search": [{"operator": "eq", "value": "test"}]}`,
			expectedStatus: 422,
			description:    "Should fail validation when parameter is missing",
		},
		{
			name:           "Missing value in scalar search",
			payload:        `{"search": [{"parameter": "hostname", "operator": "eq"}]}`,
			expectedStatus: 422,
			description:    "Should fail when value is missing in scalar search",
		},
		{
			name:           "Missing values in vector search",
			payload:        `{"search": [{"parameter": "port", "operator": "in"}]}`,
			expectedStatus: 422,
			description:    "Should fail when values is missing in vector search",
		},
		{
			name:           "Invalid operator value",
			payload:        `{"search": [{"parameter": "test", "operator": "invalid_op", "value": "test"}]}`,
			expectedStatus: 400,
			description:    "Should fail with 400 when unknown operator during unmarshal",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/modules/nmap/search", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)

			// For 422 errors, check response structure
			if tc.expectedStatus == 422 {
				var errorResponse utils.ValidationErrorResponse
				json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.Equal(t, 422, errorResponse.Code)
				assert.NotEmpty(t, errorResponse.Errors)
			}
		})
	}
}

func TestSearchNmapScansEdgeCases(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name           string
		payload        string
		expectedStatus int
		description    string
	}{
		{
			name: "Very large numeric value",
			payload: `{
				"search": [
					{
						"parameter": "port",
						"operator": "gt",
						"value": 65535
					}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle large numeric values",
		},
		{
			name: "Special characters in string value",
			payload: `{
				"search": [
					{
						"parameter": "hostname",
						"operator": "like",
						"value": "%@#$%^&*()"
					}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle special characters in values",
		},
		{
			name: "Empty array in vector search",
			payload: `{
				"search": [
					{
						"parameter": "port",
						"operator": "in",
						"values": []
					}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle empty values array",
		},
		{
			name: "Multiple sorts",
			payload: `{
				"search": [
					{
						"parameter": "status",
						"operator": "eq",
						"value": "open"
					}
				],
				"sort": [
					{"parameter": "hostname", "direction": "asc"},
					{"parameter": "port", "direction": "desc"},
					{"parameter": "protocol", "direction": "asc"}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle multiple sort specifications",
		},
		{
			name: "All scalar operators",
			payload: `{
				"search": [
					{"parameter": "a", "operator": "eq", "value": 1},
					{"parameter": "b", "operator": "neq", "value": 2},
					{"parameter": "c", "operator": "gt", "value": 3},
					{"parameter": "d", "operator": "lt", "value": 4},
					{"parameter": "e", "operator": "like", "value": "test%"},
					{"parameter": "f", "operator": "not like", "value": "%test"}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle all scalar operators",
		},
		{
			name: "All vector operators",
			payload: `{
				"search": [
					{"parameter": "a", "operator": "in", "values": [1, 2, 3]},
					{"parameter": "b", "operator": "not in", "values": ["x", "y", "z"]}
				]
			}`,
			expectedStatus: 200,
			description:    "Should handle all vector operators",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/modules/nmap/search", bytes.NewBufferString(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code, tc.description)
		})
	}
}
