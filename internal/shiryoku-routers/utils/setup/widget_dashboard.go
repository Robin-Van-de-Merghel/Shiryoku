package routers_utils_setup

import (
	shiryoku_db "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db"
	routers_widgets_dashboard "github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/widgets/dashboard"
	"github.com/gin-gonic/gin"
)

// SetupWidgetsDashboardRoutes registers all dashboard widget routes
func SetupWidgetsDashboardRoutes(dashboard_group *gin.RouterGroup, repos *shiryoku_db.Repositories) {
	dashboard_group.POST("/search", routers_widgets_dashboard.GetDashboardData(repos.Nmap))
}
