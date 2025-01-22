package services

import (
	"context"

	"github.com/axzilla/deeploy/internal/app/auth"
	"github.com/axzilla/deeploy/internal/shared/forms"
	"github.com/axzilla/deeploy/internal/shared/models"
	"github.com/axzilla/deeploy/internal/shared/repos"
	"github.com/google/uuid"
)

type ProjectServiceInterface interface {
	Create(ctx context.Context, project forms.ProjectForm) (*models.ProjectApp, error)
	Project(id string) (*models.ProjectApp, error)
	ProjectsByUser(id string) ([]models.ProjectApp, error)
	Update(project models.ProjectDB) error
	Delete(id string) error
}

type ProjectService struct {
	repo repos.ProjectRepoInterface
}

func NewProjectService(repo *repos.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, form forms.ProjectForm) (*models.ProjectApp, error) {
	project := models.ProjectDB{
		ID:          uuid.New().String(),
		UserID:      auth.GetUser(ctx).ID,
		Title:       form.Title,
		Description: form.Description,
	}
	err := s.repo.Create(&project)
	if err != nil {
		return nil, err
	}
	return project.ToProjectApp(), nil
}

func (s *ProjectService) Project(id string) (*models.ProjectApp, error) {
	project, err := s.repo.Project(id)
	if err != nil {
		return nil, err
	}
	return project.ToProjectApp(), nil
}

func (s *ProjectService) ProjectsByUser(id string) ([]models.ProjectApp, error) {
	projects, err := s.repo.ProjectsByUser(id)
	if err != nil {
		return nil, err
	}
	// TODO: Rename x to ?
	projectsApp := []models.ProjectApp{}

	for _, project := range projects {
		projectsApp = append(projectsApp, *project.ToProjectApp())
	}
	return projectsApp, nil
}

func (s *ProjectService) Update(project models.ProjectDB) error {
	err := s.repo.Update(project)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
