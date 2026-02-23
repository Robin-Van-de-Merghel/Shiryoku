package routers_widgets

import (
	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logic_widget "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

func GetDashboardData(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		var params models_widgets.WidgetDashboardInput

		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		if !utils.ValidateAndRespond(c, params, utils.WidgetDashboardSchema) {
			return
		}

		params.SetDefaults()

		results, err := logic_widget.GetLatestWidgetHosts(
			c.Request.Context(),
			nmapDB,
			params,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, results)

	}
}
