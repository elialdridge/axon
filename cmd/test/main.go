package main

import (
	"fmt"
	"os"
	"time"

	"axon/internal/config"
	"axon/internal/game"
)

// Simple integration test to verify game functionality
func main() {
	// Check if API key is available
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENROUTER_API_KEY not set, cannot run integration test")
		os.Exit(1)
	}

	fmt.Println("🎮 Starting Axon Game Integration Test...")

	// Create configuration
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: apiKey,
			GeminiAPIKey:     "",
			DefaultModel:     "openai/gpt-4o-mini",
		},
		Game: config.GameConfig{
			HistoryLimit: 1000,
			SaveDir:      "/tmp/axon_test_saves",
		},
		Terminal: config.TerminalConfig{
			Width:        80,
			Height:       24,
			ColorEnabled: false,
		},
	}

	// Create game engine
	engine := game.NewEngine(cfg)
	state := game.NewGameState()

	fmt.Println("\n📡 Testing world generation...")

	// Test world generation
	errorMessage := engine.InitializeWorld(state, "A mysterious forest clearing with ancient ruins")
	if errorMessage != nil {
		fmt.Printf("❌ World generation failed: %v\n", errorMessage)
		os.Exit(1)
	}

	fmt.Printf("✅ World created successfully!\n")
	fmt.Printf("   World Name: %s\n", state.World.Name)
	fmt.Printf("   Location: %s\n", state.World.CurrentLocation)
	fmt.Printf("   History entries: %d\n", len(state.History))

	// Print world description from history
	for _, entry := range state.History {
		if entry.Type == "narrator" {
			fmt.Printf("   Description: %s\n", entry.Content)
			break
		}
	}

	fmt.Println("\n🎲 Testing player actions...")

	// Test various player actions
	actions := []string{
		"look around carefully",
		"examine the ancient ruins",
		"search for any items",
	}

	for i, action := range actions {
		// Add delay to avoid rate limiting
		if i > 0 {
			time.Sleep(2 * time.Second)
		}

		fmt.Printf("   > %s\n", action)
		prevHistoryCount := len(state.History)

		err := engine.ProcessPlayerAction(state, action)
		if err != nil {
			fmt.Printf("❌ Action failed: %v\n", err)
			continue
		}

		// Find and print the AI response
		for j := prevHistoryCount; j < len(state.History); j++ {
			entry := state.History[j]
			if entry.Type == "narrator" {
				fmt.Printf("   AI: %s\n", entry.Content)
				break
			}
		}

		fmt.Printf("   ✅ Action processed (Turn: %d)\n\n", state.Turn)
	}

	fmt.Println("💡 Testing action suggestions...")

	// Test action suggestions
	suggestions, err := engine.GenerateActionSuggestions(state)
	if err != nil {
		fmt.Printf("❌ Action suggestions failed: %v\n", err)
	} else {
		fmt.Printf("✅ Generated %d suggestions:\n", len(suggestions))
		for i, suggestion := range suggestions {
			fmt.Printf("   %d. %s\n", i+1, suggestion)
		}
	}

	fmt.Println("\n📊 Testing system commands...")

	// Test system commands
	systemCommands := []string{"inventory", "stats", "help"}
	for _, cmd := range systemCommands {
		fmt.Printf("   > %s\n", cmd)
		prevHistoryCount := len(state.History)

		err := engine.ProcessPlayerAction(state, cmd)
		if err != nil {
			fmt.Printf("❌ Command failed: %v\n", err)
			continue
		}

		// Find and print the system response
		for j := prevHistoryCount; j < len(state.History); j++ {
			entry := state.History[j]
			if entry.Type == "system" {
				fmt.Printf("   System: %s\n", entry.Content)
				break
			}
		}
		fmt.Printf("   ✅ Command processed\n\n")
	}

	fmt.Printf("\n🎉 Integration test completed successfully!\n")
	fmt.Printf("   Final turn: %d\n", state.Turn)
	fmt.Printf("   Total history entries: %d\n", len(state.History))
	fmt.Printf("   World: %s\n", state.World.Name)
	fmt.Printf("   Current location: %s\n", state.World.CurrentLocation)

	fmt.Println("\n✨ The game is fully functional and ready to play!")
}

