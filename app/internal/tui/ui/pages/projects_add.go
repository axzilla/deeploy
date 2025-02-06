package pages

import (
	"encoding/json"
	"log"

	"github.com/axzilla/deeploy/internal/data"
	"github.com/axzilla/deeploy/internal/tui/ui/components"
	"github.com/axzilla/deeploy/internal/tui/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// /////////////////////////////////////////////////////////////////////////////
// Types & Messages
// /////////////////////////////////////////////////////////////////////////////

type ProjectAddPage struct {
	titleInput textinput.Model
	width      int
	height     int
}

type projectAddedMsg data.ProjectDTO

///////////////////////////////////////////////////////////////////////////////
// Constructors
///////////////////////////////////////////////////////////////////////////////

func NewProjectAddPage() ProjectAddPage {
	titleInput := textinput.New()
	titleInput.Focus()
	titleInput.Placeholder = "Title"

	return ProjectAddPage{
		titleInput: titleInput,
	}
}

// /////////////////////////////////////////////////////////////////////////////
// Bubbletea Interface
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectAddPage) Init() tea.Cmd {
	return textinput.Blink
}

func (p ProjectAddPage) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if len(p.titleInput.Value()) > 0 {
				return p, tea.Batch(
					p.AddProject,
					func() tea.Msg { return PopPageMsg{} },
				)
			}

		}

	case tea.WindowSizeMsg:
		p.width = msg.Width
		p.height = msg.Height
		return p, nil
	}

	p.titleInput, cmd = p.titleInput.Update(msg)
	return p, cmd
}

func (p ProjectAddPage) View() string {
	logo := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")

	card := components.Card(components.CardProps{
		Width:   p.width / 2,
		Padding: []int{2, 3},
	}).Render(p.titleInput.View())
	view := lipgloss.JoinVertical(0.5, logo, card)

	layout := lipgloss.Place(p.width, p.height, lipgloss.Center, lipgloss.Center, view)

	return layout
}

// /////////////////////////////////////////////////////////////////////////////
// Helper Methods
// /////////////////////////////////////////////////////////////////////////////

func (p ProjectAddPage) HasFocusedInput() bool {
	return p.titleInput.Focused()
}

func (p ProjectAddPage) AddProject() tea.Msg {
	postData := struct {
		Title string
	}{
		Title: p.titleInput.Value(),
	}

	res, err := utils.Request(utils.RequestProps{
		Method: "POST",
		URL:    "/projects",
		Data:   postData,
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	defer res.Body.Close()

	var project data.ProjectDTO
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		log.Println("xxx: ", err)
		return nil
	}

	return projectAddedMsg(project)
}
