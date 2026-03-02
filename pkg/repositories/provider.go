package repositories

import (
	"github.com/Robin-Van-de-Merghel/Shiryoku/pkg/utils"
)

const (
	NMAP_REPOSITORY      = "nmap"
	DASHBOARD_REPOSITORY = "dashboard"
)

// RepositoryProvider allows access to repositories and custom extensions
type RepositoryProvider interface {
	GetRepository(name string) utils.Checkable
	RegisterRepository(name string, repo utils.Checkable) error
}

// DefaultRepositoryProvider implements RepositoryProvider
type DefaultRepositoryProvider struct {
	repos map[string]utils.Checkable
}

func NewRepositoryProvider() *DefaultRepositoryProvider {
	return &DefaultRepositoryProvider{
		repos: make(map[string]utils.Checkable),
	}
}

func (p *DefaultRepositoryProvider) GetRepository(name string) utils.Checkable {
	return p.repos[name]
}

func (p *DefaultRepositoryProvider) RegisterRepository(name string, repo utils.Checkable) error {
	p.repos[name] = repo
	return nil
}
