package router_common

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logic_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

// SearchOpenSearch returns a gin handler injected with the given opensearchDB.
// Generic search for multiple db
func SearchOpenSearch[T any](opensearchClient osdb.OpenSearchClient, index string) func(c *gin.Context) {
	return func(c *gin.Context) {
		var params models.SearchParams

		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		if !utils.ValidateAndRespond(c, params, utils.SearchSchema) {
			return
		}

		params.SetDefaults()

		results, err := logic_common.Search[T](
			c.Request.Context(), 
			opensearchClient, &params, 
			index,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, results)
	}
}
