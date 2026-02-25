package routers_common

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	logic_common "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/common"
	"github.com/gin-gonic/gin"
)

// Search is a generic handler factory that returns a Gin handler for any SearchableRepository[T]
func Search[T any](repo postgres.SearchableRepository[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.SearchParams
		if err := c.ShouldBindJSON(&params); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		params.SetDefaults()

		result, err := logic_common.Search(c.Request.Context(), repo, &params)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	}
}
