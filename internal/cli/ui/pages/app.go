package ui

import (
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	tea "github.com/charmbracelet/bubbletea"
)

type AppModel struct {
	currentView viewtypes.View
	register    RegisterModel
	initConnect InitConnectModel
	login       LoginModel
	width       int
	height      int
}

func NewApp() *AppModel {
	return &AppModel{
		login:       NewLogin(0, 0),
		register:    NewRegister(0, 0),
		initConnect: NewInitConnect(0, 0),
	}
}

func (m *AppModel) Init() tea.Cmd {
	// TODO: Implement auth
	// config, err := config.LoadConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	m.currentView = viewtypes.InitConnect
	return m.initConnect.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case viewtypes.View:
		switch msg {
		case viewtypes.InitConnect:
			m.initConnect = NewInitConnect(m.width, m.height)
			m.currentView = viewtypes.InitConnect
			return m, m.initConnect.Init()
		case viewtypes.Login:
			m.login = NewLogin(m.width, m.height)
			m.currentView = viewtypes.Login
			return m, m.login.Init()
		case viewtypes.Register:
			m.register = NewRegister(m.width, m.height)
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

	if m.currentView == viewtypes.InitConnect {
		model, cmd := m.initConnect.Update(msg)
		if initConnect, ok := model.(InitConnectModel); ok {
			m.initConnect = initConnect
		}
		return m, cmd
	}
	return m, nil
}

func (m *AppModel) View() string {
	switch m.currentView {
	case viewtypes.InitConnect:
		return m.initConnect.View()
	case viewtypes.Register:
		return m.register.View()
	case viewtypes.Login:
		return m.login.View()
	case viewtypes.Dashboard:
		return "Dashboard View"
	}
	return ""
}
