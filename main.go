package main

import (
	"fmt" // Package for formatting I/O
	// Package for file system interfaces
	"os" // Package for OS functions

	"github.com/HolyStarGazer/windows-tui-go/ui" // Package for styling terminal output
	tea "github.com/charmbracelet/bubbletea"     // Package for building terminal user interfaces
)

func main() {
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
