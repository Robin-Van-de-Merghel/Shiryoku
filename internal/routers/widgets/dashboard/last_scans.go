package dashboard

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/core/models/widgets"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/common"
	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/routers/utils"
	"github.com/gin-gonic/gin"
)

func (w *DashboardWidget) getDashboardData() gin.HandlerFunc {
	return common.Search[widgets.WidgetDashboardScan](w.dasboardRepo, utils.WidgetDashboardScanFields)
}
