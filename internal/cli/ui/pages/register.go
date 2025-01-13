package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/axzilla/deeploy/internal/cli/forms"
	"github.com/axzilla/deeploy/internal/cli/ui/styles"
	"github.com/axzilla/deeploy/internal/cli/viewtypes"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	registersteps = 3
)

type RegisterModel struct {
	focusIndex int
	cursorMode cursor.Mode
	form       forms.RegisterForm
	errs       forms.RegisterErrors
}

func NewRegister() RegisterModel {
	m := RegisterModel{
		form: forms.RegisterForm{
			Email:           textinput.New(),
			Password:        textinput.New(),
			PasswordConfirm: textinput.New(),
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

	// Confirm password input
	m.form.PasswordConfirm.Placeholder = "Confirm Password"
	m.form.PasswordConfirm.EchoMode = textinput.EchoPassword
	m.form.PasswordConfirm.EchoCharacter = '•'

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
			m.resetErrs()

			// Move to the next input
			m.focusIndex++
			if m.focusIndex > registersteps {
				m.focusIndex = -1
			}
			m.updateFocus()
			return m, nil

		case tea.KeyShiftTab, tea.KeyUp:
			m.resetErrs()

			// Move to the previous input
			m.focusIndex--
			if m.focusIndex < 0 {
				m.focusIndex = registersteps
			}
			m.updateFocus()
			return m, nil

		case tea.KeyEnter:
			m.resetErrs()

			if m.focusIndex == registersteps { // Submit button
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
				if m.focusIndex > registersteps {
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
	case 2:
		m.form.PasswordConfirm, cmd = m.form.PasswordConfirm.Update(msg)
	}

	return m, cmd
}

func (m *RegisterModel) updateFocus() {
	// Reset focus and styles
	m.form.Email.Blur()
	m.form.Email.PromptStyle = styles.NoStyle
	m.form.Email.TextStyle = styles.NoStyle

	m.form.Password.Blur()
	m.form.Password.PromptStyle = styles.NoStyle
	m.form.Password.TextStyle = styles.NoStyle

	m.form.PasswordConfirm.Blur()
	m.form.PasswordConfirm.PromptStyle = styles.NoStyle
	m.form.PasswordConfirm.TextStyle = styles.NoStyle

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
	case 2:
		m.form.PasswordConfirm.Focus()
		m.form.PasswordConfirm.PromptStyle = styles.FocusedStyle
		m.form.PasswordConfirm.TextStyle = styles.FocusedStyle
	}

}

func (m *RegisterModel) resetErrs() {
	m.errs = forms.RegisterErrors{}
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString("\n" + strconv.Itoa(m.focusIndex) + "\n")

	b.WriteString("\nREGISTER\n\n")

	// Render Email input
	b.WriteString(m.form.Email.View() + "\n")

	// Render Password input
	b.WriteString(m.form.Password.View() + "\n")

	// Render Confirm Password input
	b.WriteString(m.form.PasswordConfirm.View() + "\n")
	if m.errs.Email != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.errs.Email)) + "\n")
	}
	if m.errs.Password != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.errs.Password)) + "\n")
	}
	if m.errs.PasswordConfirm != "" {
		b.WriteString(styles.ErrorStyle.Render(fmt.Sprintf("Error: %s", m.errs.PasswordConfirm)) + "\n")
	}

	// Render Submit button
	button := styles.BlurredButton
	if m.focusIndex == registersteps {
		button = styles.FocusedButton
	}
	b.WriteString("\n" + button + "\n")

	return b.String()
}
