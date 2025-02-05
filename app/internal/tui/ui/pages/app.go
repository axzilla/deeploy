package pages

import (
	"strings"

	"github.com/axzilla/deeploy/internal/tui/config"
	"github.com/axzilla/deeploy/internal/tui/ui/components"
	"github.com/axzilla/deeploy/internal/tui/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

// App is like a stack of papers (pages). The top page is what you see.
// Think of it like browser tabs or a deck of cards.
type App struct {
	stack  []tea.Model // All our pages, last one is visible
	width  int         // Terminal width
	height int         // Terminal height
}

// Message to add a new page on top of the stack
// Example: Going from login to dashboard, or opening settings
type PushPageMsg struct {
	Page tea.Model
}

// Message to remove the top page and go back to previous
// Example: Closing settings to go back to dashboard
type PopPageMsg struct{}

// /////////////////////////////////////////////////////////////////////////////
// Constructors
// /////////////////////////////////////////////////////////////////////////////

// Create new empty app
func NewApp() App {
	return App{
		stack: make([]tea.Model, 0),
	}

}

// Called when app starts - we wait for window size before creating pages
func (a App) Init() tea.Cmd {
	return nil
}

// Handles all events/messages in the app
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit (ctrl+c, esc)
	if msg, ok := msg.(tea.KeyMsg); ok {
		if msg.String() == "q" {
			return a, tea.Quit
		}
	}

	switch msg := msg.(type) {
	// When terminal size changes
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

		// When adding a new page (push)
		// Use push when:
		// - Opening a new view on top (settings, details, etc)
		// - Want to keep previous page in history
		// - Need "back" functionality
	case PushPageMsg:
		newPage := msg.Page

		// Add to stack
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

	// When removing top page (pop)
	// Use pop when:
	// - Closing a sub-view
	// - Going "back" to previous view
	// - Cancelling an operation
	case PopPageMsg:
		// Only pop if we have more than one page
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

// Shows current (top) page
func (a App) View() string {
	if len(a.stack) == 0 {
		return "Loading..."
	}

	main := a.stack[len(a.stack)-1].View()

	footerMenuItems := []FooterMenuItem{
		{Key: "q", Desc: "quit"},
		{Key: "esc", Desc: "back"},
	}

	var footer strings.Builder

	for i, v := range footerMenuItems {
		footer.WriteString(v.Key)
		footer.WriteString(" ")
		footer.WriteString(styles.BlurredStyle.Render(v.Desc))
		if len(footerMenuItems)-1 != i {
			footer.WriteString(" • ")
		}
	}

	footerCard := components.Card(components.CardProps{
		Width:   a.width,
		Padding: []int{0, 1},
	}).Render(footer.String())

	view := lipgloss.JoinVertical(lipgloss.Center, main, footerCard)

	return view
}

// Examples of when to use Push/Pop:
//
// Push examples:
// - Login -> Dashboard (push dashboard)
// - Dashboard -> Settings (push settings)
// - Dashboard -> User Details (push details)
// - Dashboard -> Help Page (push help)
//
// Pop examples:
// - Settings -> Back to Dashboard (pop settings)
// - Help -> Back to previous view (pop help)
// - Details -> Back to list (pop details)
//
// Navigation flow example:
// 1. [Login]
// 2. Push Dashboard -> [Login, Dashboard]
// 3. Push Settings -> [Login, Dashboard, Settings]
// 4. Pop Settings -> [Login, Dashboard]
// 5. Push Help -> [Login, Dashboard, Help]
// 6. Pop Help -> [Login, Dashboard]
