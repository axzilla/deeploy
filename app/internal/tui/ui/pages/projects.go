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
	width         int
	height        int
	message       string
	projects      []data.ProjectDTO
	selectedIndex int
	err           error
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
		switch msg.String() {
		case "n":
			return p, func() tea.Msg {
				return PushPageMsg{Page: NewProjectFormPage(nil)}
			}
		case "e":
			return p, func() tea.Msg {
				return PushPageMsg{Page: NewProjectFormPage(&p.projects[p.selectedIndex])}
			}
		case "down", "j":
			if p.selectedIndex == len(p.projects)-1 {
				p.selectedIndex = 0
			} else {
				p.selectedIndex++
			}
		case "up", "k":
			if p.selectedIndex == 0 {
				p.selectedIndex = len(p.projects) - 1
			} else {
				p.selectedIndex--
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
	case projectCreatedMsg:
		p.projects = append(p.projects, data.ProjectDTO(msg))
		return p, nil
	case projectUpdatedMsg:
		project := msg
		for i, v := range p.projects {
			if v.ID == project.ID {
				p.projects[i] = data.ProjectDTO(project)
			}
		}
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
		for i, project := range p.projects {
			props := components.CardProps{
				Width:   30,
				Padding: []int{0, 1},
			}
			if p.selectedIndex == i {
				props.BorderForeground = styles.ColorPrimary
			}
			cards = append(cards, components.Card(props).Render(project.Title))
		}
	}

	projectsView := lipgloss.JoinVertical(0.5, cards...)

	if len(cards) == 0 {
		projectsView = components.Card(components.CardProps{Width: 30}).Align(lipgloss.Position(0.5)).Render(styles.FocusedStyle.Render("No projects yet"))
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
