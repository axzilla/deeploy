package ui

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/axzilla/deeploy/internal/cli/config"
	"github.com/axzilla/deeploy/internal/cli/ui/components"
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DashboardModel struct {
	width   int
	height  int
	message string
}

func NewDashboard(width, height int) DashboardModel {
	return DashboardModel{
		width:  width,
		height: height,
	}
}

type welcomeMessage string

func getWelcomeMessage() tea.Msg {
	config, err := config.LoadConfig()
	if err != nil {
		log.Println("Config error:", err)
		return viewtypes.InitConnect
	}

	r, err := http.NewRequest("POST", "http://"+config.Server+"/dashboard", nil)
	if err != nil {
		log.Println(err)
		return viewtypes.InitConnect
	}
	r.Header.Set("Authorization", "Bearer "+config.Token)

	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Println(err)
		return viewtypes.InitConnect
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("HTTP error:", err)
		return viewtypes.InitConnect
	}
	log.Println("Got result:", string(result))
	return welcomeMessage(result)
}

func (m DashboardModel) Init() tea.Cmd {
	return getWelcomeMessage
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case welcomeMessage:
		return DashboardModel{
			width:   m.width,
			height:  m.height,
			message: string(msg),
		}, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
	var b strings.Builder

	b.WriteString(strconv.Itoa(m.width))

	logo := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(0).Render(m.message)
	footer := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("\n[ctrl+c] exit")

	view := lipgloss.JoinVertical(0.5, logo, card, footer)
	layout := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}
