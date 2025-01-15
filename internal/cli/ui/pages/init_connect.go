package ui

import (
	"os/exec"
	"strings"

	"github.com/axzilla/deeploy/internal/cli/ui/components"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InitConnectModel struct {
	serverInput textinput.Model
	status      string
	waiting     bool
	width       int
	height      int
	err         string
}

func NewInitConnect(width, height int) InitConnectModel {
	ti := textinput.New()
	ti.Placeholder = "e.g. 123.45.67.89:8090"
	ti.Focus()

	return InitConnectModel{
		serverInput: ti,
		width:       width,
		height:      height,
	}
}

func (m InitConnectModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InitConnectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		m.resetErr()

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.validate()
			if m.err == "" {
				m.waiting = true
				return m, func() tea.Msg {
					return exec.Command("open", "https://deeploy.sh").Run()
				}
			}
		}
	}

	m.serverInput, cmd = m.serverInput.Update(msg)
	return m, cmd
}

func (m InitConnectModel) View() string {
	var b strings.Builder

	if m.waiting {
		b.WriteString("✨ Browser opened for authentication. Waiting for completion.")
	} else {
		b.WriteString("CONNECT TO SERVER\n\n")
		b.WriteString(styles.FocusedStyle.Render("Server "))
		b.WriteString(m.serverInput.View())
		if m.err != "" {
			b.WriteString(styles.ErrorStyle.Render("\n* " + m.err))
		}
		if m.status != "" {
			b.WriteString(m.status)
		}
	}

	logo := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(50).Render(b.String())
	footer := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("\n[ctrl+c] exit")

	view := lipgloss.JoinVertical(0.5, logo, card, footer)
	layout := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}

func (m *InitConnectModel) validate() {
	if m.serverInput.Value() == "" {
		m.err = "Server required"
	}
	// TODO: Add add. validations like: no valid ip/server whatever and so on
}

func (m *InitConnectModel) resetErr() {
	m.err = ""
}
