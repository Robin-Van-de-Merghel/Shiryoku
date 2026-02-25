package dashboard

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/shiryoku-routers/utils"
	"github.com/gin-gonic/gin"
)

func GetDashboardData(dashboardRepo postgres.DashboardRepository) gin.HandlerFunc {
	return common.Search[widgets.WidgetDashboardScan](dashboardRepo, utils.WidgetDashboardScanFields)
}
