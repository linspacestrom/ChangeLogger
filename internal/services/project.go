package services

import (
	"context"
	"log"

	"github.com/linspacestrom/ChangeLogger/internal/domain"
	"github.com/linspacestrom/ChangeLogger/internal/repositories"
)

type ProjectService struct {
	repo repositories.ProjectRepository
}

func NewProjectService(repo repositories.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, projectCreate domain.ProjectCreateOrUpdate) (domain.Project, error) {
	project, err := s.repo.Create(ctx, projectCreate)
	if err != nil {
		log.Printf("ProjectService failed: %s", err)
	}
	return project, err

}

func (s *ProjectService) ListProjects(ctx context.Context) ([]domain.Project, error) {
	projects, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("[ProjectService] failed to list projects: %v\n", err)
	}
	return projects, err
}
