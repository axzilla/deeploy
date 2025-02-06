package components

import (
	"github.com/axzilla/deeploy/internal/tui/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

type CardProps struct {
	Width            int
	Padding          []int
	BorderForeground lipgloss.TerminalColor
}

func Card(p CardProps) lipgloss.Style {
	baseStyle := lipgloss.NewStyle().
		Width(p.Width).
		Border(lipgloss.RoundedBorder())

	actualWidth := p.Width - baseStyle.GetHorizontalBorderSize()

	return lipgloss.NewStyle().
		Width(actualWidth).
		Padding(p.Padding...).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(p.BorderForeground)
}

func ErrorCard(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(styles.ColorError).
		Width(width).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
}
