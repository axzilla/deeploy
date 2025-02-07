package messages

import (
	"github.com/axzilla/deeploy/internal/data"
	tea "github.com/charmbracelet/bubbletea"
)

// Navigation Messages
type PushPageMsg struct {
	Page tea.Model
}
type PopPageMsg struct{}

// Auth Messages
type AuthErrorMsg struct {
	Err error
}
type AuthSuccessMsg struct {
	token string
}

// Project Messages
type ProjectCreatedMsg data.ProjectDTO
type ProjectUpdatedMsg data.ProjectDTO
type ProjectDeleteMsg *data.ProjectDTO
type ProjectErrMsg error
type ProjectsInitDataMsg []data.ProjectDTO
