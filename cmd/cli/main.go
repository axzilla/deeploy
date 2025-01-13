package main

import (
	"fmt"
	"os"

	"github.com/axzilla/deeploy/internal/cli/ui/pages"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := ui.NewApp()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
