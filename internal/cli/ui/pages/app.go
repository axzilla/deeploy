package ui

import (
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	currentView viewtypes.View
	register    RegisterModel
	login       LoginModel
}

func NewApp() *AppModel {
	return &AppModel{
		login:    NewLogin(),
		register: NewRegister(),
	}
}

func (m *AppModel) Init() tea.Cmd {
	m.currentView = viewtypes.Register
	return m.register.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewtypes.View:
		switch msg {
		case viewtypes.Login:
			m.login = NewLogin()
			m.currentView = viewtypes.Login
			return m, m.login.Init()
		case viewtypes.Register:
			m.register = NewRegister()
			m.currentView = viewtypes.Register
			return m, m.register.Init()
		case viewtypes.Dashboard:
			m.currentView = viewtypes.Dashboard
			return m, nil
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
		if register, ok := model.(RegisterModel); ok {
			m.register = register
		}
		return m, cmd
	}

	if m.currentView == viewtypes.Login {
		model, cmd := m.login.Update(msg)
		if login, ok := model.(LoginModel); ok {
			m.login = login
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
		return m.login.View()
	case viewtypes.Dashboard:
		return "Dashboard View"
	}
	return ""
}
