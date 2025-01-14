package ui

import (
	"strings"

	"github.com/axzilla/deeploy/internal/app/utils"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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
}

func NewRegister() RegisterModel {
	m := RegisterModel{
		fields: make([]RegisterField, 3),
		errs:   make(map[int]string),
	}
	for i := range m.fields {
		t := textinput.New()
		switch i {
		case RegisterFormEmail:
			// m.form[i].label = styles.LabelStyle.Render("email")
			t.Placeholder = "email"
			t.Focus()
			t.PromptStyle = styles.FocusedStyle
			t.TextStyle = styles.FocusedStyle
		case RegisterFormPassword:
			// m.form[i].label = styles.LabelStyle.Render("password")
			t.Placeholder = "password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case RegisterFormPasswordConfirm:
			// m.form[i].label = styles.LabelStyle.Render("confirm password")
			t.Placeholder = "confirm password"
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
			if m.focusIndex == RegisterFormSubmit { // Submit
				m.validate()
				if len(m.errs) == 0 {
					return m, func() tea.Msg { return viewtypes.Dashboard }
				}
			} else if m.focusIndex == RegisterFormLoginLink { // Login Link
				return m, func() tea.Msg { return viewtypes.Login }
			} else {
				m.nextInput()
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
		m.errs[RegisterFormEmail] = "not a valid email"
	}

	if m.fields[RegisterFormEmail].input.Value() == "" {
		m.errs[RegisterFormEmail] = "email is required"
	}

	if m.fields[RegisterFormPassword].input.Value() == "" {
		m.errs[RegisterFormPassword] = "password is required"
	}

	if m.fields[RegisterFormPassword].input.Value() != m.fields[RegisterFormPasswordConfirm].input.Value() {
		m.errs[RegisterFormPassword] = "passwords do not match"
	}
}

func (m *RegisterModel) resetErrs() {
	m.errs = make(map[int]string)
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString("\nREGISTER\n\n")

	for _, field := range m.fields {
		// b.WriteString(field.label + "\n")
		b.WriteString(field.input.View() + "\n")
	}

	if len(m.errs) > 0 {
		b.WriteString("\n")
	}
	if err, ok := m.errs[RegisterFormEmail]; ok {
		b.WriteString(styles.ErrorStyle.Render("* "+err) + "\n")
	}
	if err, ok := m.errs[RegisterFormPassword]; ok {
		b.WriteString(styles.ErrorStyle.Render("* "+err) + "\n")
	}
	if err, ok := m.errs[RegisterFormPasswordConfirm]; ok {
		b.WriteString(styles.ErrorStyle.Render("* "+err) + "\n")
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

	return b.String()
}
