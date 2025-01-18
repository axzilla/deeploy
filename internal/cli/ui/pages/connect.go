package pages

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/axzilla/deeploy/internal/cli/config"
	"github.com/axzilla/deeploy/internal/cli/ui/components"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ConnectPage struct {
	serverInput textinput.Model
	status      string
	waiting     bool
	width       int
	height      int
	err         string
}

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

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewConnectPage() ConnectPage {
	ti := textinput.New()
	ti.Placeholder = "e.g. 123.45.67.89:8090"
	ti.Focus()

	return ConnectPage{
		serverInput: ti,
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ConnectPage) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConnectPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
	case AuthSuccessMsg:
		return m, func() tea.Msg {
			return viewtypes.Dashboard
		}
	}
	m.serverInput, cmd = m.serverInput.Update(msg)
	return m, cmd
}

func (p ConnectPage) View() string {
	var b strings.Builder

	if p.waiting {
		b.WriteString("✨ Browser opened for authentication. Waiting for completion.")
	} else {
		b.WriteString("CONNECT TO SERVER\n\n")
		b.WriteString(styles.FocusedStyle.Render("Server "))
		b.WriteString(p.serverInput.View())
		if p.err != "" {
			b.WriteString(styles.ErrorStyle.Render("\n* " + p.err))
		}
		if p.status != "" {
			b.WriteString(p.status)
		}
	}

	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")
	card := components.Card(50).Render(b.String())
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

func (p *ConnectPage) validate() {
	value := strings.TrimSpace(p.serverInput.Value())

	// 1. Base-Check
	if value == "" {
		p.err = "Server required"
		return
	}

	// 2. Format-Check (Host:Port)
	host, port, err := net.SplitHostPort(value)
	if err != nil {
		p.err = "Invalid format. Use host:port (e.g. server:8090)"
		return
	}

	// 3. Port-Validation
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		p.err = "Invalid port number"
		return
	}

	// 4. Simple Host-Check
	if host == "" {
		p.err = "Invalid host"
		return
	}

	// Connection test
	timeout := time.Second * 3
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		p.err = "Server not reachable"
		p.status = ""
		return
	}
	defer conn.Close()
}

func (p *ConnectPage) resetErr() {
	p.err = ""
}

// Starts a local server for ayth callback
func startLocalAuthServer() (int, chan authCallback) {
	callback := make(chan authCallback)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /callback", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		token, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		callback <- authCallback{token: string(token)}
		w.Write([]byte("OK"))
		return
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

func (m ConnectPage) startBrowserAuth() tea.Cmd {
	return func() tea.Msg {
		port, callback := startLocalAuthServer()

		// Open browser
		authURL := fmt.Sprintf(
			"http://%s?cli=true&port=%d",
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

		return PushPageMsg{Page: NewDashboard()}
	}
}
