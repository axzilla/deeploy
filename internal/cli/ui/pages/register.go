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
)

type RegisterField struct {
	input textinput.Model
	label string
}

type RegisterModel struct {
	focusIndex int
	cursorMode cursor.Mode
	form       []RegisterField
	errs       map[int]string
}

func NewRegister() RegisterModel {
	m := RegisterModel{
		form: make([]RegisterField, 3),
		errs: make(map[int]string),
	}
	for i := range m.form {
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
		m.form[i].input = t
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

			if m.focusIndex == len(m.form) { // Submit button
				m.validate()
				if len(m.errs) == 0 {
					cmd := func() tea.Msg { return viewtypes.Dashboard }
					return m, cmd
					// return m, tea.Quit
				}
			} else {
				m.nextInput()
				return m, nil
			}
		}
	}

	// Update the currently focused input field
	var cmd tea.Cmd
	switch m.focusIndex {
	case RegisterFormEmail:
		m.form[RegisterFormEmail].input, cmd = m.form[RegisterFormEmail].input.Update(msg)
	case RegisterFormPassword:
		m.form[RegisterFormPassword].input, cmd = m.form[RegisterFormPassword].input.Update(msg)
	case RegisterFormPasswordConfirm:
		m.form[RegisterFormPasswordConfirm].input, cmd = m.form[RegisterFormPasswordConfirm].input.Update(msg)
	}

	return m, cmd
}

func (m *RegisterModel) nextInput() {
	m.focusIndex++
	if m.focusIndex > len(m.form) {
		m.focusIndex = 0
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *RegisterModel) prevInput() {
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = len(m.form)
	}
	m.updateFocus()
	m.resetErrs()
}

func (m *RegisterModel) updateFocus() {
	for i := range m.form {
		if i == m.focusIndex {
			m.form[i].input.Focus()
			m.form[i].input.PromptStyle = styles.FocusedStyle
			m.form[i].input.TextStyle = styles.FocusedStyle

		} else {
			m.form[i].input.Blur()
			m.form[i].input.PromptStyle = styles.NoStyle
			m.form[i].input.TextStyle = styles.NoStyle
		}
	}
}

func (m *RegisterModel) validate() {
	if !utils.IsEmailValid(m.form[RegisterFormEmail].input.Value()) {
		m.errs[RegisterFormEmail] = "not a valid email"
	}

	if m.form[RegisterFormEmail].input.Value() == "" {
		m.errs[RegisterFormEmail] = "email is required"
	}

	if m.form[RegisterFormPassword].input.Value() == "" {
		m.errs[RegisterFormPassword] = "password is required"
	}

	if m.form[RegisterFormPassword].input.Value() != m.form[RegisterFormPasswordConfirm].input.Value() {
		m.errs[RegisterFormPassword] = "passwords do not match"
	}
}

func (m *RegisterModel) resetErrs() {
	m.errs = make(map[int]string)
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString("\nREGISTER\n\n")

	for _, field := range m.form {
		// b.WriteString(field.label + "\n")
		b.WriteString(field.input.View() + "\n")
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

	button := styles.BlurredButton
	if m.focusIndex == len(m.form) {
		button = styles.FocusedButton
	}
	b.WriteString("\n" + button + "\n")

	return b.String()
}
