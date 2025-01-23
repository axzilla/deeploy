package services

import (
	"context"

	"github.com/axzilla/deeploy/internal/auth"
	"github.com/axzilla/deeploy/internal/data"
	"github.com/axzilla/deeploy/internal/forms"
	"github.com/google/uuid"
)

type ProjectServiceInterface interface {
	Create(ctx context.Context, project forms.ProjectForm) (*data.Project, error)
	Project(id string) (*data.Project, error)
	ProjectsByUser(id string) ([]data.Project, error)
	Update(project data.Project) error
	Delete(id string) error
}

type ProjectService struct {
	repo data.ProjectRepoInterface
}

func NewProjectService(repo *data.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) Create(ctx context.Context, form forms.ProjectForm) (*data.Project, error) {
	project := data.Project{
		ID:          uuid.New().String(),
		UserID:      auth.GetUser(ctx).ID,
		Title:       form.Title,
		Description: form.Description,
	}
	err := s.repo.Create(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *ProjectService) Project(id string) (*data.Project, error) {
	project, err := s.repo.Project(id)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) ProjectsByUser(id string) ([]data.Project, error) {
	projects, err := s.repo.ProjectsByUser(id)
	if err != nil {
		return nil, err
	}
	// TODO: Rename x to ?
	projectsApp := []data.Project{}

	for _, project := range projects {
		projectsApp = append(projectsApp, project)
	}
	return projectsApp, nil
}

func (s *ProjectService) Update(project data.Project) error {
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
