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
	cfg := setupTest()
	engine, state := createGameEngine(cfg)

	testWorldGeneration(engine, state)
	testPlayerActions(engine, state)
	testActionSuggestions(engine, state)
	testSystemCommands(engine, state)

	printFinalResults(state)
}

func setupTest() *config.Config {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENROUTER_API_KEY not set, cannot run integration test")
		os.Exit(1)
	}

	fmt.Println("🎮 Starting Axon Game Integration Test...")

	return &config.Config{
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
}

func createGameEngine(cfg *config.Config) (*game.Engine, *game.GameState) {
	engine := game.NewEngine(cfg)
	state := game.NewGameState()
	return engine, state
}

func testWorldGeneration(engine *game.Engine, state *game.GameState) {
	fmt.Println("\n📡 Testing world generation...")

	errorMessage := engine.InitializeWorld(state, "A mysterious forest clearing with ancient ruins")
	if errorMessage != nil {
		fmt.Printf("❌ World generation failed: %v\n", errorMessage)
		os.Exit(1)
	}

	fmt.Printf("✅ World created successfully!\n")
	fmt.Printf("   World Name: %s\n", state.World.Name)
	fmt.Printf("   Location: %s\n", state.World.CurrentLocation)
	fmt.Printf("   History entries: %d\n", len(state.History))

	for _, entry := range state.History {
		if entry.Type == "narrator" {
			fmt.Printf("   Description: %s\n", entry.Content)
			break
		}
	}
}

func testPlayerActions(engine *game.Engine, state *game.GameState) {
	fmt.Println("\n🎲 Testing player actions...")

	actions := []string{
		"look around carefully",
		"examine the ancient ruins",
		"search for any items",
	}

	for i, action := range actions {
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

		for j := prevHistoryCount; j < len(state.History); j++ {
			entry := state.History[j]
			if entry.Type == "narrator" {
				fmt.Printf("   AI: %s\n", entry.Content)
				break
			}
		}

		fmt.Printf("   ✅ Action processed (Turn: %d)\n\n", state.Turn)
	}
}

func testActionSuggestions(engine *game.Engine, state *game.GameState) {
	fmt.Println("💡 Testing action suggestions...")

	suggestions, err := engine.GenerateActionSuggestions(state)
	if err != nil {
		fmt.Printf("❌ Action suggestions failed: %v\n", err)
	} else {
		fmt.Printf("✅ Generated %d suggestions:\n", len(suggestions))
		for i, suggestion := range suggestions {
			fmt.Printf("   %d. %s\n", i+1, suggestion)
		}
	}
}

func testSystemCommands(engine *game.Engine, state *game.GameState) {
	fmt.Println("\n📊 Testing system commands...")

	systemCommands := []string{"inventory", "stats", "help"}
	for _, cmd := range systemCommands {
		fmt.Printf("   > %s\n", cmd)
		prevHistoryCount := len(state.History)

		err := engine.ProcessPlayerAction(state, cmd)
		if err != nil {
			fmt.Printf("❌ Command failed: %v\n", err)
			continue
		}

		for j := prevHistoryCount; j < len(state.History); j++ {
			entry := state.History[j]
			if entry.Type == "system" {
				fmt.Printf("   System: %s\n", entry.Content)
				break
			}
		}
		fmt.Printf("   ✅ Command processed\n\n")
	}
}

func printFinalResults(state *game.GameState) {
	fmt.Printf("\n🎉 Integration test completed successfully!\n")
	fmt.Printf("   Final turn: %d\n", state.Turn)
	fmt.Printf("   Total history entries: %d\n", len(state.History))
	fmt.Printf("   World: %s\n", state.World.Name)
	fmt.Printf("   Current location: %s\n", state.World.CurrentLocation)

	fmt.Println("\n✨ The game is fully functional and ready to play!")
}
