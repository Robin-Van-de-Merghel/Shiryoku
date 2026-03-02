package nmap

import (
	"fmt"

	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/repositories"
	"github.com/gin-gonic/gin"
)

type NmapModule struct {
	nmapRepo repositories.NmapRepository
}

func (m *NmapModule) Name() string {
	return "nmap"
}

func (m *NmapModule) Description() string {
	return "Nmap scan results"
}

func (m *NmapModule) SetupRoutes(nmap_group *gin.RouterGroup, provider repositories.RepositoryProvider) error {
	repo := provider.GetRepository(repositories.NMAP_REPOSITORY)
	if repo == nil {
		return fmt.Errorf("couldn't import repository %s from provider", repositories.NMAP_REPOSITORY)
	}

	nmapRepo, ok := repo.(repositories.NmapRepository)
	if !ok {
		return fmt.Errorf("repository %s is not an NmapRepository", repositories.NMAP_REPOSITORY)
	}

	m.nmapRepo = nmapRepo

	search_group := nmap_group.Group("/search")
	search_group.POST("", m.searchNmapScans())
	nmap_group.POST("/batch", m.insertNmapScans())

	return nil
}
