package components

import (
	"github.com/axzilla/deeploy/internal/tui/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

func Card(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
}

func ErrorCard(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(styles.ColorError).
		Width(width).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
}
