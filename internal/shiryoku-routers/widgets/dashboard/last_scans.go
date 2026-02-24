package routers_widgets_dashboard

import (
	"time"

	models_widgets "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	logic_widgets_dashboard "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-logic/widgets/dashboard"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var dashboardCache = cache.New(5*time.Minute, 10*time.Minute)

func GetDashboardData(nmapDB osdb.OpenSearchClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		const cacheKey = "dashboard_data"

		// Try cache first
		if cached, found := dashboardCache.Get(cacheKey); found {
			c.JSON(200, cached)
			return
		}

		var params models_widgets.WidgetDashboardInput

		if err := c.ShouldBindJSON(&params); err != nil {
			utils.ParseJSONError(c, err)
			return
		}

		if !utils.ValidateAndRespond(c, params, utils.WidgetDashboardSchema) {
			return
		}

		params.SetDefaults()

		results, err := logic_widgets_dashboard.GetLatestWidgetScans(
			c.Request.Context(),
			nmapDB,
			params,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		// Store in cache
		dashboardCache.Set(cacheKey, results, cache.DefaultExpiration)

		c.JSON(200, results)

	}
}
