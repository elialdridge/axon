//go:build debug
// +build debug

package main

import (
	"fmt"
	"os"

	"axon/internal/ai"
	"axon/internal/config"
	"axon/internal/game"
	"axon/internal/logger"
)

func main() {
	// Initialize logger
	if err := logger.Init(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Info("Starting world creation debug test")

	// Load configuration
	cfg := config.Load()
	fmt.Printf("API Key configured: %t\n", cfg.AI.OpenRouterAPIKey != "")
	logger.Debug("Config: %+v", cfg)

	// Create AI client
	client := ai.NewClient(cfg.AI.OpenRouterAPIKey, cfg.AI.GeminiAPIKey)
	fmt.Printf("AI Client created\n")

	// Test simple AI request first
	fmt.Printf("Testing simple AI request...\n")
	req := ai.Request{
		Prompt:    "Say hello",
		Model:     "mistralai/mistral-7b-instruct:free",
		MaxTokens: 50,
		Context:   []string{"You are a helpful assistant."},
	}

	resp, err := client.Generate(req)
	if err != nil {
		fmt.Printf("Simple AI request failed: %v\n", err)
	} else if resp.Error != nil {
		fmt.Printf("Simple AI request error: %v\n", resp.Error)
	} else {
		fmt.Printf("Simple AI request successful: %s\n", resp.Text)
	}

	// Test world creation
	fmt.Printf("\nTesting world creation...\n")
	engine := game.NewEngine(cfg)
	gameState := game.NewGameState()

	err = engine.InitializeWorld(gameState, "A cyberpunk city in 2077")
	if err != nil {
		fmt.Printf("World creation failed: %v\n", err)
	} else {
		fmt.Printf("World creation successful!\n")
		fmt.Printf("World Name: %s\n", gameState.World.Name)
		fmt.Printf("World Description: %s\n", gameState.World.Description)
		fmt.Printf("History entries: %d\n", len(gameState.History))
		for i, entry := range gameState.History {
			fmt.Printf("  %d. [%s] %s\n", i+1, entry.Type, entry.Content)
		}
	}

	fmt.Printf("\nCheck axon_debug.log for detailed logs\n")
}
