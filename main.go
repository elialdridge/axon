package main

import (
	"fmt"
	"os"

	"axon/internal/config"
	"axon/internal/game"
	"axon/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize logger
	if err := logger.Init(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Starting Axon game")

	// Load configuration
	cfg := config.Load()
	logger.Debug("Configuration loaded: %+v", cfg)

	// Initialize game model
	model := game.NewModel(cfg)
	logger.Info("Game model initialized")

	// Create and start the Bubble Tea program with proper terminal detection
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	logger.Info("Starting Bubble Tea program")

	if _, err := p.Run(); err != nil {
		logger.Error("Error running game: %v", err)
		fmt.Printf("Error running game: %v\n", err)
		os.Exit(1)
	}

	logger.Info("Game ended normally")
}
