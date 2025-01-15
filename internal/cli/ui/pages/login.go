package ui

import (
	"strings"

	"github.com/axzilla/deeploy/internal/app/utils"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	LoginFormEmail = iota
	LoginFormPassword
	LoginFormSubmit
	LoginFormRegisterLink
)

type LoginField struct {
	input textinput.Model
	label string
}

type LoginModel struct {
	focusIndex int
	cursorMode cursor.Mode
	fields     []LoginField
	errs       map[int]string
	width      int
	height     int
}

func NewLogin(w, h int) LoginModel {
	m := LoginModel{
		fields: make([]LoginField, 2),
		errs:   make(map[int]string),
		width:  w,
		height: h,
	}
	for i := range m.fields {
		t := textinput.New()
		switch i {
		case LoginFormEmail:
			m.fields[i].label = "Email"
			t.Placeholder = "your@email.com"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case LoginFormPassword:
			m.fields[i].label = "Password"
			t.Placeholder = "enter password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.fields[i].input = t
	}
	return m
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyTab, tea.KeyDown:
			m.nextInput()
			return m, tea.Batch(textinput.Blink)

		case tea.KeyShiftTab, tea.KeyUp:
			m.prevInput()
			return m, tea.Batch(textinput.Blink)

		case tea.KeyEnter:
			m.resetErrs()
			if m.focusIndex == LoginFormSubmit { // Submit
				m.validate()
				if len(m.errs) == 0 {
					return m, func() tea.Msg { return viewtypes.Dashboard }
				}
			} else if m.focusIndex == LoginFormRegisterLink { // Login Link
				return m, func() tea.Msg { return viewtypes.Register }
			} else {
				m.nextInput()
				return m, tea.Batch(textinput.Blink)
			}
		}
	}

	// Update the currently focused input field
	var cmd tea.Cmd
	switch m.focusIndex {
	case LoginFormEmail:
		m.fields[LoginFormEmail].input, cmd = m.fields[LoginFormEmail].input.Update(msg)
	case LoginFormPassword:
		m.fields[LoginFormPassword].input, cmd = m.fields[LoginFormPassword].input.Update(msg)
	}

	return m, cmd
}

func (m *LoginModel) nextInput() {
	m.focusIndex++
	if m.focusIndex > LoginFormRegisterLink {
		m.focusIndex = 0
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *LoginModel) prevInput() {
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = LoginFormRegisterLink
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *LoginModel) updateFocus() {
	for i := range m.fields {
		if i == m.focusIndex {
			m.fields[i].input.Focus()
			m.fields[i].input.PromptStyle = styles.FocusedStyle
			m.fields[i].input.TextStyle = styles.FocusedStyle

		} else {
			m.fields[i].input.Blur()
			m.fields[i].input.PromptStyle = styles.NoStyle
			m.fields[i].input.TextStyle = styles.NoStyle
		}
	}
}

func (m *LoginModel) validate() {
	if !utils.IsEmailValid(m.fields[LoginFormEmail].input.Value()) {
		m.errs[LoginFormEmail] = "Invalid email"
	}

	if m.fields[LoginFormEmail].input.Value() == "" {
		m.errs[LoginFormEmail] = "Email required"
	}

	if m.fields[LoginFormPassword].input.Value() == "" {
		m.errs[LoginFormPassword] = "Password required"
	}
}

func (m *LoginModel) resetErrs() {
	m.errs = make(map[int]string)
}

func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString("\nLOGIN\n\n")

	for i, field := range m.fields {
		// Label
		if field.label != "" {
			if m.focusIndex == i {
				b.WriteString(styles.FocusedStyle.Render(field.label) + "\n")
			} else {
				b.WriteString(field.label + "\n")
			}
		}

		// Input
		b.WriteString(field.input.View() + "\n")

		// Error
		if err, ok := m.errs[i]; ok {
			b.WriteString(styles.ErrorStyle.Render("* "+err) + "\n\n")
		} else {
			b.WriteString("\n")
		}

	}

	// Submit Button
	button := styles.BlurredButton
	if m.focusIndex == LoginFormSubmit {
		button = styles.FocusedButton
	}
	b.WriteString("\n" + button + "\n")

	//  Register Link
	loginText := "Don't have an account yet?"
	if m.focusIndex == LoginFormRegisterLink {
		loginText = styles.FocusedStyle.Render(loginText)
	}
	b.WriteString("\n" + loginText + "\n")

	card := styles.AuthCard.Render(b.String())
	logo := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("🔥deeploy.sh\n")

	view := lipgloss.JoinVertical(0.5, logo, card)
	layout := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, view)
	return layout
}
