package main

import (
	"fmt"
	"os"

	"axon/internal/config"
	"axon/internal/game"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize game model
	model := game.NewModel(cfg)

	// Create and start the Bubble Tea program
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}
}
