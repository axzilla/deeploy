package components

import "github.com/charmbracelet/lipgloss"

func Card(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder())
}
