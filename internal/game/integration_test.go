package game

import (
	"os"
	"strings"
	"testing"
	"time"

	"axon/internal/config"
)

// TestWorldGenerationIntegration tests world generation with real API
func TestWorldGenerationIntegration(t *testing.T) {
	// Skip if no API key is available
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		t.Skip("OPENROUTER_API_KEY not set, skipping integration test")
	}

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
	}

	engine := NewEngine(cfg)

	// Test world generation with different prompts
	testCases := []struct {
		name   string
		prompt string
	}{
		{"Fantasy World", "A medieval fantasy kingdom with magic and dragons"},
		{"Sci-Fi World", "A space station orbiting a distant planet"},
		{"Modern World", "A bustling modern city with hidden secrets"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create fresh state for each test
			testState := NewGameState()

			// Test world initialization
			err := engine.InitializeWorld(testState, tc.prompt)
			if err != nil {
				t.Errorf("InitializeWorld failed: %v", err)
				return
			}

			// Verify world was created
			if testState.World.Name == "" {
				t.Error("World name should not be empty")
			}

			if testState.World.Description == "" {
				t.Error("World description should not be empty")
			}

			if len(testState.World.Rules) == 0 {
				t.Error("World should have rules")
			}

			if testState.World.CurrentLocation == "" {
				t.Error("Current location should not be empty")
			}

			// Verify history was populated
			if len(testState.History) == 0 {
				t.Error("Should have history entries after world creation")
			}

			// Check that we have a narrator entry describing the world
			foundNarratorEntry := false
			for _, entry := range testState.History {
				if entry.Type == "narrator" && len(entry.Content) > 0 {
					foundNarratorEntry = true
					break
				}
			}
			if !foundNarratorEntry {
				t.Error("Should have narrator entry describing the world")
			}

			t.Logf("Successfully created world: %s", testState.World.Name)
			t.Logf("Description: %s", testState.World.Description)
			t.Logf("Current Location: %s", testState.World.CurrentLocation)
		})
	}
}

// TestPlayerActionsIntegration tests player actions with real API
func TestPlayerActionsIntegration(t *testing.T) {
	// Skip if no API key is available
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		t.Skip("OPENROUTER_API_KEY not set, skipping integration test")
	}

	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: apiKey,
			GeminiAPIKey:     "",
			DefaultModel:     "openai/gpt-4o-mini",
		},
	}

	engine := NewEngine(cfg)
	state := NewGameState()

	// Initialize world first
	err := engine.InitializeWorld(state, "A simple tavern in a fantasy town")
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	initialHistoryCount := len(state.History)
	initialTurn := state.Turn

	// Test various player actions
	actions := []string{
		"look around",
		"examine the room",
		"talk to the bartender",
		"order a drink",
	}

	for i, action := range actions {
		t.Run("Action_"+strings.ReplaceAll(action, " ", "_"), func(t *testing.T) {
			prevHistoryCount := len(state.History)
			prevTurn := state.Turn

			// Add some delay to avoid rate limiting
			if i > 0 {
				time.Sleep(1 * time.Second)
			}

			err := engine.ProcessPlayerAction(state, action)
			if err != nil {
				t.Errorf("ProcessPlayerAction failed for '%s': %v", action, err)
				return
			}

			// Verify history was updated
			if len(state.History) <= prevHistoryCount {
				t.Errorf("History should have been updated after action '%s'", action)
			}

			// Verify turn advanced
			if state.Turn != prevTurn+1 {
				t.Errorf("Turn should advance after action '%s', expected %d, got %d", action, prevTurn+1, state.Turn)
			}

			// Check that player action was recorded
			foundPlayerAction := false
			foundNarratorResponse := false
			for j := prevHistoryCount; j < len(state.History); j++ {
				entry := state.History[j]
				if entry.Type == "player" && entry.Content == action {
					foundPlayerAction = true
				}
				if entry.Type == "narrator" && len(entry.Content) > 0 {
					foundNarratorResponse = true
				}
			}

			if !foundPlayerAction {
				t.Errorf("Player action '%s' should be recorded in history", action)
			}

			if !foundNarratorResponse {
				t.Errorf("Should have narrator response to action '%s'", action)
			}

			t.Logf("Action '%s' processed successfully, turn %d", action, state.Turn)
		})
	}

	// Verify overall game state
	if len(state.History) <= initialHistoryCount {
		t.Error("Game history should have grown")
	}

	if state.Turn <= initialTurn {
		t.Error("Game turn should have advanced")
	}

	t.Logf("Integration test completed. Final turn: %d, History entries: %d", state.Turn, len(state.History))
}

// TestActionSuggestionsIntegration tests action suggestions with real API
func TestActionSuggestionsIntegration(t *testing.T) {
	// Skip if no API key is available
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		t.Skip("OPENROUTER_API_KEY not set, skipping integration test")
	}

	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: apiKey,
			GeminiAPIKey:     "",
			DefaultModel:     "openai/gpt-4o-mini",
		},
	}

	engine := NewEngine(cfg)
	state := NewGameState()

	// Initialize world
	err := engine.InitializeWorld(state, "A mysterious forest clearing with an ancient stone circle")
	if err != nil {
		t.Fatalf("Failed to initialize world: %v", err)
	}

	// Test action suggestions
	suggestions, err := engine.GenerateActionSuggestions(state)
	if err != nil {
		t.Errorf("GenerateActionSuggestions failed: %v", err)
		return
	}

	if len(suggestions) == 0 {
		t.Error("Should generate at least some action suggestions")
		return
	}

	t.Logf("Generated %d action suggestions:", len(suggestions))
	for i, suggestion := range suggestions {
		t.Logf("  %d. %s", i+1, suggestion)
		
		// Verify suggestions are not empty
		if strings.TrimSpace(suggestion) == "" {
			t.Errorf("Suggestion %d should not be empty", i+1)
		}
	}

	// Test that suggestions change after player actions
	err = engine.ProcessPlayerAction(state, "examine the stone circle")
	if err != nil {
		t.Errorf("Failed to process player action: %v", err)
		return
	}

	// Add delay to avoid rate limiting
	time.Sleep(1 * time.Second)

	newSuggestions, err := engine.GenerateActionSuggestions(state)
	if err != nil {
		t.Errorf("GenerateActionSuggestions failed after action: %v", err)
		return
	}

	if len(newSuggestions) == 0 {
		t.Error("Should generate suggestions after player action")
		return
	}

	t.Logf("Generated %d new action suggestions after player action:", len(newSuggestions))
	for i, suggestion := range newSuggestions {
		t.Logf("  %d. %s", i+1, suggestion)
	}
}

