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
	RegisterFormEmail = iota
	RegisterFormPassword
	RegisterFormPasswordConfirm
	RegisterFormSubmit
	RegisterFormLoginLink
)

type RegisterField struct {
	input textinput.Model
	label string
}

type RegisterModel struct {
	focusIndex int
	cursorMode cursor.Mode
	fields     []RegisterField
	errs       map[int]string
	width      int
	height     int
}

func NewRegister(w, h int) RegisterModel {
	m := RegisterModel{
		fields: make([]RegisterField, 3),
		errs:   make(map[int]string),
		width:  w,
		height: h,
	}
	for i := range m.fields {
		t := textinput.New()
		switch i {
		case RegisterFormEmail:
			m.fields[i].label = "Email"
			t.Placeholder = "your@email.com"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case RegisterFormPassword:
			m.fields[i].label = "Password"
			t.Placeholder = "min. 8 chars"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case RegisterFormPasswordConfirm:
			m.fields[i].label = "Confirm"
			t.Placeholder = "repeat password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		m.fields[i].input = t
	}
	return m
}

func (m RegisterModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focusIndex == RegisterFormSubmit { // Submit
				m.validate()
				if len(m.errs) == 0 {
					return m, func() tea.Msg { return viewtypes.Dashboard }
				}
			} else if m.focusIndex == RegisterFormLoginLink { // Login Link
				return m, func() tea.Msg { return viewtypes.Login }
			} else {
				m.nextInput()
				return m, tea.Batch(textinput.Blink)
			}
		}
	}

	// Update the currently focused input field
	var cmd tea.Cmd
	switch m.focusIndex {
	case RegisterFormEmail:
		m.fields[RegisterFormEmail].input, cmd = m.fields[RegisterFormEmail].input.Update(msg)
	case RegisterFormPassword:
		m.fields[RegisterFormPassword].input, cmd = m.fields[RegisterFormPassword].input.Update(msg)
	case RegisterFormPasswordConfirm:
		m.fields[RegisterFormPasswordConfirm].input, cmd = m.fields[RegisterFormPasswordConfirm].input.Update(msg)
	}

	return m, cmd
}

func (m *RegisterModel) nextInput() {
	m.focusIndex++
	if m.focusIndex > RegisterFormLoginLink {
		m.focusIndex = 0
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *RegisterModel) prevInput() {
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = RegisterFormLoginLink
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *RegisterModel) updateFocus() {
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

func (m *RegisterModel) validate() {
	if !utils.IsEmailValid(m.fields[RegisterFormEmail].input.Value()) {
		m.errs[RegisterFormEmail] = "Invalid email"
	}

	if m.fields[RegisterFormEmail].input.Value() == "" {
		m.errs[RegisterFormEmail] = "Email required"
	}

	if len(m.fields[RegisterFormPassword].input.Value()) < 8 {
		m.errs[RegisterFormPassword] = "Min. 8 chars"
	}

	if m.fields[RegisterFormPassword].input.Value() != m.fields[RegisterFormPasswordConfirm].input.Value() {
		m.errs[RegisterFormPasswordConfirm] = "Passwords don't match"
	}
}

func (m *RegisterModel) resetErrs() {
	m.errs = make(map[int]string)
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString("\nREGISTER\n\n")

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
	if m.focusIndex == RegisterFormSubmit {
		button = styles.FocusedButton
	}
	b.WriteString("\n" + button + "\n")

	// Login Link
	loginText := "Already have an account?"
	if m.focusIndex == RegisterFormLoginLink {
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
