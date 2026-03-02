package common

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/logic/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/utils"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
	"github.com/gin-gonic/gin"
)

// Search is a generic Gin handler factory for any SearchableRepository[T].
// It validates field names + types, and maps to actual DB columns.
func Search[T any](repo repositories.SearchableRepository[T], allowedMaps ...map[string]utils.FieldTypeInfo) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params models.SearchParams

		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		if !utils.ValidateAndRespond(c, &params, utils.SearchSchema) {
			return
		}

		params.SetDefaults()

		// Validate all parameters exist in allowed maps and types match
		if err := utils.ValidateSearchParamTypesPrecomputed(&params, allowedMaps...); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Execute the search
		result, err := common.Search[T](c.Request.Context(), repo, &params)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, result)
	}
}
