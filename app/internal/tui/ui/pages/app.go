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

type Viewstack struct {
	stack []tea.Model
}

type HasInputView interface {
	HasFocusedInput() bool
}

type App struct {
	stack       []tea.Model
	currentPage tea.Model
	width       int
	height      int
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

		currentPage := a.currentPage
		if page, ok := currentPage.(HasInputView); ok && page.HasFocusedInput() {
			// this disable "q"
		} else if msg.String() == "q" {
			return a, tea.Quit
		}

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		// If no pages yet, create first one
		if a.currentPage == nil {
			config, err := config.LoadConfig()
			var page tea.Model

			// No config = show login, has config = show dashboard
			if err != nil || config.Server == "" || config.Token == "" {
				page = NewConnectPage()
			} else {
				page = NewDashboard()
			}

			// Add first page to stack
			a.currentPage = page

			// Update page with window size and initialize it
			updatedPage, cmd := page.Update(msg)
			a.currentPage = updatedPage
			return a, tea.Batch(cmd, updatedPage.Init())
		}

		// Update current page's window size
		currentPage := a.currentPage
		updatedPage, cmd := currentPage.Update(msg)
		a.currentPage = updatedPage
		return a, cmd

	case messages.ChangePageMsg:
		newPage := msg.Page

		a.currentPage = newPage

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

	// All other messages go to current page
	default:
		if a.currentPage == nil {
			return a, nil
		}
		currentPage := a.currentPage
		updatedPage, cmd := currentPage.Update(msg)
		a.currentPage = updatedPage
		return a, cmd
	}

	return a, nil
}

type FooterMenuItem struct {
	Key  string
	Desc string
}

func (a App) View() string {
	if a.currentPage == nil {
		return "Loading..."
	}

	main := a.currentPage.View()

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
