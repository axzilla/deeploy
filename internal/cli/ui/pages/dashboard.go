package pages

import (
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/axzilla/deeploy/internal/cli/config"
	"github.com/axzilla/deeploy/internal/cli/ui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type DashboardPage struct {
	width   int
	height  int
	message string
}

type welcomeMessage string

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewDashboard() DashboardPage {
	return DashboardPage{}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p DashboardPage) Init() tea.Cmd {
	return getWelcomeMessage
}

func (p DashboardPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	case welcomeMessage:
		p.message = string(msg) // Nur die Message updaten
		return p, nil           // Kein zusätzliches Command
	}
	return p, nil
}

func (p DashboardPage) View() string {
	var b strings.Builder

	b.WriteString(strconv.Itoa(p.width))

	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(0).Render(p.message)
	footer := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("\n[ctrl+c] exit")

	view := lipgloss.JoinVertical(0.5, logo, card, footer)
	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func getWelcomeMessage() tea.Msg {
	config, err := config.LoadConfig()
	if err != nil {
		return PushPageMsg{Page: NewConnectPage()}
	}

	r, err := http.NewRequest("POST", "http://"+config.Server+"/dashboard", nil)
	if err != nil {
		return PushPageMsg{Page: NewConnectPage()}
	}
	r.Header.Set("Authorization", "Bearer "+config.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return PushPageMsg{Page: NewConnectPage()}
	}
	if res.StatusCode == http.StatusUnauthorized {
		return PushPageMsg{Page: NewConnectPage()}
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return PushPageMsg{Page: NewConnectPage()}
	}

	return welcomeMessage(result)
}
