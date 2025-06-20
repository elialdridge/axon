package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"axon/internal/config"
	"axon/internal/game"
	"axon/internal/logger"
	"axon/internal/terminal"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Close()

	logger.Info("Starting Axon game")

	// Detect terminal capabilities
	termInfo := terminal.DetectTerminal()
	logger.Info("Terminal detection completed")

	// Load configuration
	cfg := config.Load()
	logger.Debug("Configuration loaded: %+v", cfg)

	// Apply terminal detection to configuration if auto-detect is enabled
	if cfg.Terminal.AutoDetect {
		applyTerminalDetection(cfg, termInfo)
		logger.Info("Applied terminal detection to configuration")
	}

	// Initialize game model
	model := game.NewModel(cfg, termInfo)
	logger.Info("Game model initialized")

	// Create Bubble Tea program options based on terminal capabilities
	options := buildTeaOptions(cfg, termInfo)
	p := tea.NewProgram(model, options...)
	logger.Info("Starting Bubble Tea program")

	if _, err := p.Run(); err != nil {
		logger.Error("Error running game: %v", err)
		fmt.Printf("Error running game: %v\n", err)
		return
	}

	logger.Info("Game ended normally")
}

// applyTerminalDetection applies detected terminal capabilities to configuration
func applyTerminalDetection(cfg *config.Config, termInfo *terminal.TerminalInfo) {
	// Override size with detected values if they seem reasonable
	width, height := termInfo.GetSafeSize()
	if width > 0 && height > 0 {
		cfg.Terminal.Width = width
		cfg.Terminal.Height = height
	}

	// Disable features for minimal or System V terminals
	if termInfo.IsMinimal || termInfo.IsSystemV {
		cfg.Terminal.MouseEnabled = false
		cfg.Terminal.AltScreenEnabled = false
		cfg.Terminal.ColorEnabled = false
	} else {
		// Enable features based on terminal capabilities
		cfg.Terminal.MouseEnabled = termInfo.MouseSupport
		cfg.Terminal.AltScreenEnabled = termInfo.AltScreenSupport
		cfg.Terminal.ColorEnabled = termInfo.ColorSupport
	}

	// Force minimal mode for very constrained terminals
	if termInfo.Width < 50 || termInfo.Height < 15 {
		cfg.Terminal.ForceMinimal = true
	}

	logger.Debug("Applied terminal settings - Mouse: %v, AltScreen: %v, Color: %v",
		cfg.Terminal.MouseEnabled, cfg.Terminal.AltScreenEnabled, cfg.Terminal.ColorEnabled)
}

// buildTeaOptions builds Bubble Tea program options based on configuration and terminal info
func buildTeaOptions(cfg *config.Config, termInfo *terminal.TerminalInfo) []tea.ProgramOption {
	var options []tea.ProgramOption

	// Add alt screen support if enabled and supported
	if cfg.Terminal.AltScreenEnabled && termInfo.AltScreenSupport {
		options = append(options, tea.WithAltScreen())
		logger.Debug("Enabled alt screen support")
	}

	// Add mouse support if enabled and supported
	if cfg.Terminal.MouseEnabled && termInfo.MouseSupport {
		options = append(options, tea.WithMouseCellMotion())
		logger.Debug("Enabled mouse support")
	}

	// For minimal terminals, use minimal rendering
	if cfg.Terminal.ForceMinimal || termInfo.IsMinimal {
		// Minimal terminals get basic functionality only
		logger.Debug("Using minimal terminal mode")
	}

	return options
}
