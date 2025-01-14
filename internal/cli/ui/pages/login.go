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
	LoginFormLoginLink
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
			// m.form[i].label = styles.LabelStyle.Render("email")
			t.Placeholder = "email"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case LoginFormPassword:
			// m.form[i].label = styles.LabelStyle.Render("password")
			t.Placeholder = "password"
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
			return m, nil

		case tea.KeyShiftTab, tea.KeyUp:
			m.prevInput()
			return m, nil

		case tea.KeyEnter:
			m.resetErrs()
			if m.focusIndex == LoginFormSubmit { // Submit
				m.validate()
				if len(m.errs) == 0 {
					return m, func() tea.Msg { return viewtypes.Dashboard }
				}
			} else if m.focusIndex == LoginFormLoginLink { // Login Link
				return m, func() tea.Msg { return viewtypes.Register }
			} else {
				m.nextInput()
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
	if m.focusIndex > LoginFormLoginLink {
		m.focusIndex = 0
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *LoginModel) prevInput() {
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = LoginFormLoginLink
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
		m.errs[LoginFormEmail] = "not a valid email"
	}

	if m.fields[LoginFormEmail].input.Value() == "" {
		m.errs[LoginFormEmail] = "email is required"
	}

	if m.fields[LoginFormPassword].input.Value() == "" {
		m.errs[LoginFormPassword] = "password is required"
	}
}

func (m *LoginModel) resetErrs() {
	m.errs = make(map[int]string)
}

func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString("\nLOGIN\n\n")

	for _, field := range m.fields {
		// b.WriteString(field.label + "\n")
		b.WriteString(field.input.View() + "\n")
	}

	if len(m.errs) > 0 {
		b.WriteString("\n")
	}
	if err, ok := m.errs[LoginFormEmail]; ok {
		b.WriteString(styles.ErrorStyle.Render("* "+err) + "\n")
	}
	if err, ok := m.errs[LoginFormPassword]; ok {
		b.WriteString(styles.ErrorStyle.Render("* "+err) + "\n")
	}

	// Submit Button
	button := styles.BlurredButton
	if m.focusIndex == LoginFormSubmit {
		button = styles.FocusedButton
	}
	b.WriteString("\n" + button + "\n")

	// Login Link
	loginText := "Don't have an account yet?"
	if m.focusIndex == LoginFormLoginLink {
		loginText = styles.FocusedStyle.Render(loginText)
	}
	b.WriteString("\n" + loginText + "\n")

	card := styles.AuthCard.Render(b.String())
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
}
