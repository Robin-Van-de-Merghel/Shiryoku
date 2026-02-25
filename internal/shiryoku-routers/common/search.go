package common

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

// Search is a generic handler factory that returns a Gin handler for any SearchableRepository[T]
func Search[T any](repo postgres.SearchableRepository[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.SearchParams
		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		// Validate fields and stop if invalid
		if !utils.ValidateAndRespond(c, &params, utils.SearchSchema) {
			return
		}

		params.SetDefaults()
		result, err := common.Search[T](c.Request.Context(), repo, &params)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
	}
}
