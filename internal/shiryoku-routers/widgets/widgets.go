package routers_widgets

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/config"
	"github.com/gin-gonic/gin"
)

func SetupWidgetsRoutes(widget_group *gin.RouterGroup, widgetModuleConfig config.APIModule) {
	// Routes
	dashboard_group := widget_group.Group("/dashboard")

	dashboard_group.POST("/search", GetDashboardData(widgetModuleConfig.OSDB))
}
