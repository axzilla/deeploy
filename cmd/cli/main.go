package main

import (
	"fmt"
	"github.com/axzilla/deeploy/internal/cli/ui/pages"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {
	// Logging Setup
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Start App
	m := ui.NewApp()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
