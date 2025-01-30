package pages

import (
	"encoding/json"
	"net/http"

	"github.com/axzilla/deeploy/internal/data"
	"github.com/axzilla/deeploy/internal/tui/config"
	"github.com/axzilla/deeploy/internal/tui/ui/components"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type DashboardPage struct {
	width    int
	height   int
	message  string
	projects []data.ProjectDTO
	list     list.Model
	err      error
}

type errMsg error
type projectsMsg []data.ProjectDTO

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewDashboard() DashboardPage {
	return DashboardPage{
		list: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p DashboardPage) Init() tea.Cmd {
	return getInitData
}

func (p DashboardPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		p.list.SetSize(msg.Width/2, msg.Height/2)
		return p, nil
	case errMsg:
		p.err = msg
		return p, nil
	case projectsMsg:
		p.projects = msg
		//
		// return p, nil

		items := make([]list.Item, len(p.projects))
		for i, project := range p.projects {
			items[i] = item(project.Title)
		}

		// p.list.Title = "Choose your project"

		// p.list.Items.(true)
		p.list.SetItems(items)
		return p, nil
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
	// return p, nil

}

type model struct {
	list list.Model
}

type item string

// func (i item) Title() string       { return i.title }
// func (i item) Description() string { return "" }
func (i item) FilterValue() string { return "" }

// func (i item) SelectedItem() item {
// 	index := m.Index()
//
// 	items := m.VisibleItems()
// 	if i < 0 || len(items) == 0 || len(items) <= i {
// 		return nil
// 	}
//
// 	return items[i]
// }

func (p DashboardPage) View() string {
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
	projectsListCard := components.Card(40).Render(p.list.View())
	projectsView := lipgloss.JoinVertical(0.5, cards...)
	footer := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("\n[ctrl+c] exit")

	view := lipgloss.JoinVertical(0.5, logo, projectsView, projectsListCard, footer)
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
