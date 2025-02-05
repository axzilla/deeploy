package pages

import (
	"encoding/json"
	"net/http"

	"github.com/axzilla/deeploy/internal/data"
	"github.com/axzilla/deeploy/internal/tui/config"
	"github.com/axzilla/deeploy/internal/tui/ui/components"
	"github.com/axzilla/deeploy/internal/tui/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectPage struct {
	width    int
	height   int
	message  string
	projects []data.ProjectDTO
	err      error
}

type errMsg error
type projectsMsg []data.ProjectDTO

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectPage() ProjectPage {
	return ProjectPage{}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectPage) Init() tea.Cmd {
	return getInitData
}

func (p ProjectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return p, func() tea.Msg {
				return PopPageMsg{}
			}
		}
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	case errMsg:
		p.err = msg
	case projectsMsg:
		p.projects = msg
		return p, nil
	}
	return p, nil
}

func (p ProjectPage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")

	var cards []string
	if p.err != nil {
		cards = append(cards, components.ErrorCard(30).Render(p.err.Error()))
	} else {
		for _, project := range p.projects {
			cards = append(cards, components.Card(30).Render(project.Title))
		}
	}

	projectsView := lipgloss.JoinVertical(0.5, cards...)

	if len(cards) == 0 {
		projectsView = components.Card(30).Align(lipgloss.Position(0.5)).Render(styles.FocusedStyle.Render("No projects yet"))
	}

	view := lipgloss.JoinVertical(0.5, logo, projectsView)

	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)

	return layout
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func getInitData() tea.Msg {
	config, err := config.LoadConfig()
	if err != nil {
		return PushPageMsg{Page: NewConnectPage()}
	}

	r, err := http.NewRequest("GET", "http://"+config.Server+"/api/projects", nil)
	if err != nil {
		return errMsg(err)
	}
	r.Header.Set("Authorization", "Bearer "+config.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return errMsg(err)
	}
	if res.StatusCode == http.StatusUnauthorized {
		return PushPageMsg{Page: NewConnectPage()}
	}
	defer res.Body.Close()

	var projects []data.ProjectDTO
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return errMsg(err)
	}

	return projectsMsg(projects)
}
