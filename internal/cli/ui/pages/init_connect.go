package ui

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/axzilla/deeploy/internal/cli/config"
	"github.com/axzilla/deeploy/internal/cli/ui/components"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type authCallback struct {
	token string
	err   error
}

type AuthErrorMsg struct {
	err error
}

type AuthSuccessMsg struct {
	token string
}

// Starts a local server for ayth callback
func startLocalAuthServer() (int, chan authCallback) {
	callback := make(chan authCallback)

	mux := http.NewServeMux()

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, _ := io.ReadAll(r.Body)
		callback <- authCallback{token: string(token)}
		w.Write([]byte("OK"))
	})

	// Get a free random port
	listener, _ := net.Listen("tcp", "localhost:0")
	port := listener.Addr().(*net.TCPAddr).Port

	go http.Serve(listener, mux)

	return port, callback
}

func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "bsd", etc.
		cmd = "xdg-open"
	}

	return exec.Command(cmd, append(args, url)...).Start()
}

func (m InitConnectModel) startBrowserAuth() tea.Cmd {
	return func() tea.Msg {
		port, callback := startLocalAuthServer()
		log.Println(port)

		// Open browser
		authURL := fmt.Sprintf(
			"http://%s/cli-auth?port=%d",
			m.serverInput.Value(),
			port,
		)
		openBrowser(authURL)

		// Waiting for token
		result := <-callback
		if result.err != nil {
			return AuthErrorMsg{err: result.err}
		}

		// Save config
		cfg := config.Config{
			Server: m.serverInput.Value(),
			Token:  result.token,
		}
		if err := config.SaveConfig(&cfg); err != nil {
			return AuthErrorMsg{err: err}
		}

		return AuthSuccessMsg{token: result.token}
	}
}

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
				return m, m.startBrowserAuth()
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
