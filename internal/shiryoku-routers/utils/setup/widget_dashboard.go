package routers_utils_setup

import (
	osdb "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/opensearch"
	routers_widgets_dashboard "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/widgets/dashboard"
	"github.com/gin-gonic/gin"
)

func SetupWidgetsDashboardRoutes(
	dashboard_group *gin.RouterGroup,
	OSDB osdb.OpenSearchClient,
	path string,
) {
	dashboard_group.POST("/search", routers_widgets_dashboard.GetDashboardData(OSDB))
}
