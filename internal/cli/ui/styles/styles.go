package styles

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	FocusedStyle                       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle                       = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	CursorStyle                        = FocusedStyle
	NoStyle                            = lipgloss.NewStyle()
	HelpStyle                          = BlurredStyle
	CursorModeHelpStyle                = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	ErrorStyle          lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))
)
