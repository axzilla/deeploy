package ui

import (
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	currentView viewtypes.View
	register    RegisterModel
}

func NewApp() *AppModel {
	return &AppModel{
		register: NewRegister(),
	}
}

func (m *AppModel) Init() tea.Cmd {
	m.currentView = viewtypes.Register

	return tea.Batch(
		m.register.Init(),
	)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewtypes.View:
		switch msg {
		case viewtypes.Login:
			m.currentView = viewtypes.Login
		case viewtypes.Register:
			m.currentView = viewtypes.Register
		case viewtypes.Dashboard:
			m.currentView = viewtypes.Dashboard
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	// Send updates to "children"
	if m.currentView == viewtypes.Register {
		model, cmd := m.register.Update(msg)
		if reg, ok := model.(RegisterModel); ok {
			m.register = reg
		}
		return m, cmd
	}
	return m, nil
}

func (m *AppModel) View() string {
	switch m.currentView {
	case viewtypes.Register:
		return m.register.View()
	case viewtypes.Login:
		return "Login View"
	case viewtypes.Dashboard:
		return "Dashboard View"
	}
	return ""
}
