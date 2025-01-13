package ui

import (
	"fmt"
	"strings"

	"github.com/axzilla/deeploy/internal/cli/forms"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	focusIndex int
	cursorMode cursor.Mode
	form       forms.LoginForm
	errs       forms.LoginErrors
}

func NewLogin() LoginModel {
	m := LoginModel{
		form: forms.LoginForm{
			Email:    textinput.New(),
			Password: textinput.New(),
		},
	}

	// Email input
	m.form.Email.Placeholder = "Email"
	m.form.Email.Focus()
	m.form.Email.PromptStyle = styles.FocusedStyle
	m.form.Email.TextStyle = styles.FocusedStyle

	// Password input
	m.form.Password.Placeholder = "Password"
	m.form.Password.EchoMode = textinput.EchoPassword
	m.form.Password.EchoCharacter = '•'

	return m
}

func (m LoginModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab":
			m.resetErrs()

			// Move to the next input
			m.focusIndex++
			if m.focusIndex > 3 {
				m.focusIndex = -1
			}
			m.updateFocus()
			return m, nil

		case "shift+tab":
			m.resetErrs()

			// Move to the previous input
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = 3
			}
			m.updateFocus()
			return m, nil

		case "up":
			m.resetErrs()

			// Move to the previous input
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = 3
			}
			m.updateFocus()
			return m, nil

		case "down":
			m.resetErrs()

			// Move to the next input
			m.focusIndex++
			if m.focusIndex > 3 {
				m.focusIndex = 0
			}
			m.updateFocus()
			return m, nil

		case "enter":
			m.resetErrs()

			if m.focusIndex == 3 { // Submit button
				m.errs = m.form.Validate()

				if !m.errs.HasErrors() {
					fmt.Println("Form submitted successfully!")
					cmd := func() tea.Msg { return viewtypes.Dashboard }
					return m, cmd
					// return m, tea.Quit
				}
			} else {
				// Move to the next input
				m.focusIndex++
				if m.focusIndex > 3 {
					m.focusIndex = 0
				}
				m.updateFocus()
				return m, nil
			}
		}
	}

	// Update the currently focused input field
	var cmd tea.Cmd
	switch m.focusIndex {
	case 0:
		m.form.Email, cmd = m.form.Email.Update(msg)
	case 1:
		m.form.Password, cmd = m.form.Password.Update(msg)
	}

	return m, cmd
}

func (m *LoginModel) updateFocus() {
	// Reset focus and styles
	m.form.Email.Blur()
	m.form.Email.PromptStyle = styles.NoStyle
	m.form.Email.TextStyle = styles.NoStyle

	m.form.Password.Blur()
	m.form.Password.PromptStyle = styles.NoStyle
	m.form.Password.TextStyle = styles.NoStyle

	// Set focus abd styles bases on current index
	switch m.focusIndex {
	case 0:
		m.form.Email.Focus()
		m.form.Email.PromptStyle = styles.FocusedStyle
		m.form.Email.TextStyle = styles.FocusedStyle
	case 1:
		m.form.Password.Focus()
		m.form.Password.PromptStyle = styles.FocusedStyle
		m.form.Password.TextStyle = styles.FocusedStyle
	}
}

func (m *LoginModel) resetErrs() {
	m.errs = forms.LoginErrors{}
}

func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString("\nREGISTER\n\n")

	// Render Email input
	b.WriteString(m.form.Email.View() + "\n")

	// Render Password input
	b.WriteString(m.form.Password.View() + "\n")

	if m.errs.Email != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.errs.Email)) + "\n")
	}
	if m.errs.Password != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.errs.Password)) + "\n")
	}

	// Render Submit button
	button := styles.BlurredButton
	if m.focusIndex == 3 {
		button = styles.FocusedButton
	}
	b.WriteString("\n" + button + "\n")

	return b.String()
}
