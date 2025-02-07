package pages

import (
	"strings"

	"github.com/axzilla/deeploy/internal/tui/config"
	"github.com/axzilla/deeploy/internal/tui/messages"
	"github.com/axzilla/deeploy/internal/tui/ui/components"
	"github.com/axzilla/deeploy/internal/tui/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ActiveArea string

const (
	menu     ActiveArea = "menu"
	projects ActiveArea = "projects"
	settings ActiveArea = "settings"
	logs     ActiveArea = "logs"
)

type Viewstack struct {
	stack []tea.Model
}

type HasInputView interface {
	HasFocusedInput() bool
}

type App struct {
	stack  []tea.Model
	width  int
	height int
}

// /////////////////////////////////////////////////////////////////////////////
// Constructors
// /////////////////////////////////////////////////////////////////////////////

func NewApp() App {
	return App{
		stack: make([]tea.Model, 0),
	}

}

// We wait for window size before creating pages
func (a App) Init() tea.Cmd {
	return nil
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}

		currentPage := a.stack[len(a.stack)-1]
		if page, ok := currentPage.(HasInputView); ok && page.HasFocusedInput() {
			// this disable "q"
		} else if msg.String() == "q" {
			return a, tea.Quit
		}

		if msg.Type == tea.KeyEsc {
			return a, func() tea.Msg {
				return messages.PopPageMsg{}
			}
		}

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// If no pages yet, create first one
		if len(a.stack) == 0 {
			config, err := config.LoadConfig()
			var page tea.Model

			// No config = show login, has config = show dashboard
			if err != nil || config.Server == "" || config.Token == "" {
				page = NewConnectPage()
			} else {
				page = NewDashboard()
			}

			// Add first page to stack
			a.stack = append(a.stack, page)

			// Update page with window size and initialize it
			updatedPage, cmd := page.Update(msg)
			a.stack[len(a.stack)-1] = updatedPage
			return a, tea.Batch(cmd, updatedPage.Init())
		}

		// Update current page's window size
		currentPage := a.stack[len(a.stack)-1]
		updatedPage, cmd := currentPage.Update(msg)
		a.stack[len(a.stack)-1] = updatedPage
		return a, cmd

	case messages.PushPageMsg:
		newPage := msg.Page

		a.stack = append(a.stack, newPage)

		// Batch window size and init commands together
		// This prevents double rendering by ensuring both happen in sequence
		return a, tea.Batch(
			func() tea.Msg {
				return tea.WindowSizeMsg{
					Width:  a.width,
					Height: a.height,
				}
			},
			newPage.Init(),
		)

	case messages.PopPageMsg:
		if len(a.stack) > 1 {
			a.stack = a.stack[:len(a.stack)-1]
			return a, nil
		}

	// All other messages go to current page
	default:
		if len(a.stack) == 0 {
			return a, nil
		}
		currentPage := a.stack[len(a.stack)-1]
		updatedPage, cmd := currentPage.Update(msg)
		a.stack[len(a.stack)-1] = updatedPage
		return a, cmd
	}

	return a, nil
}

type FooterMenuItem struct {
	Key  string
	Desc string
}

func (a App) View() string {
	if len(a.stack) == 0 {
		return "Loading..."
	}

	main := a.stack[len(a.stack)-1].View()

	footerMenuItems := []FooterMenuItem{
		{Key: ":", Desc: "menu"},
		{Key: "esc", Desc: "back"},
		{Key: "q", Desc: "quit"},
	}

	var footer strings.Builder

	for i, v := range footerMenuItems {
		footer.WriteString(styles.FocusedStyle.Render(v.Key))
		footer.WriteString(" ")
		footer.WriteString(v.Desc)
		if len(footerMenuItems)-1 != i {
			footer.WriteString(" • ")
		}
	}

	footerCard := components.Card(components.CardProps{
		Width:   a.width,
		Padding: []int{0, 1},
	}).Render(footer.String())

	horizontal := lipgloss.JoinHorizontal(0.5, main)
	view := lipgloss.JoinVertical(lipgloss.Bottom, horizontal, footerCard)

	return view
}
