package routers_widgets_dashboard

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	logic_widgets_dashboard "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/widgets/dashboard"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

func GetDashboardData(dashboardRepo postgres.DashboardRepository) func(c *gin.Context) {
	return func(c *gin.Context) {
		const cacheKey = "dashboard_data"

		var params models.SearchParams

		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		if !utils.ValidateAndRespond(c, params, utils.SearchSchema) {
			return
		}

		params.SetDefaults()

		results, err := logic_widgets_dashboard.GetLatestWidgetScans(
			c.Request.Context(),
			dashboardRepo,
			&params,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, results)

	}
}
