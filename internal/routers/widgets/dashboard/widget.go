package dashboard

import (
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/internal/db/postgres"
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
	"github.com/gin-gonic/gin"
)

type DashboardWidget struct {
	dasboardRepo postgres.DashboardRepository
}

func (w *DashboardWidget) Name() string {
	return "dashboard"
}

func (w *DashboardWidget) Description() string {
	return "Dashboard for the WUI"
}

func (w *DashboardWidget) SetupRoutes(dashboard_group *gin.RouterGroup, provider repositories.RepositoryProvider) error {
	repo := provider.GetRepository(repositories.DASHBOARD_REPOSITORY)
	if repo == nil {
		return fmt.Errorf("couldn't import repository %s from provider", repositories.DASHBOARD_REPOSITORY)
	}

	nmapRepo, ok := repo.(postgres.DashboardRepository)
	if !ok {
		return fmt.Errorf("repository %s is not an DashboardRepository", repositories.DASHBOARD_REPOSITORY)
	}

	w.dasboardRepo = nmapRepo

	dashboard_group.POST("/search", w.getDashboardData())

	return nil
}
