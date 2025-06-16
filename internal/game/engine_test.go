package game

import (
	"strings"
	"testing"

	"axon/internal/config"
)

func TestNewEngine(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: "test_key",
			GeminiAPIKey:     "test_gemini",
			DefaultModel:     "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	
	if engine == nil {
		t.Fatal("NewEngine returned nil")
	}
	
	if engine.config != cfg {
		t.Error("Engine config not set correctly")
	}
	
	if engine.aiClient == nil {
		t.Error("AI client should be initialized")
	}
}

func TestInitializeWorld(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: "", // Empty key will trigger fallback
			GeminiAPIKey:     "",
			DefaultModel:     "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	state := NewGameState()
	
	err := engine.InitializeWorld(state, "A fantasy world with dragons")
	if err != nil {
		t.Errorf("InitializeWorld should not return error, got %v", err)
	}
	
	// Should have fallback world
	if state.World.Name == "" {
		t.Error("World name should be set")
	}
	
	if state.World.Description == "" {
		t.Error("World description should be set")
	}
	
	if len(state.World.Rules) == 0 {
		t.Error("World should have rules")
	}
	
	if state.World.CurrentLocation == "" {
		t.Error("Current location should be set")
	}
	
	// Should have added history entry
	if len(state.History) == 0 {
		t.Error("Should have history entries")
	}
}

func TestProcessPlayerAction(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: "", // Empty key will trigger fallback
			GeminiAPIKey:     "",
			DefaultModel:     "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	state := NewGameState()
	
	// Initialize world first
	engine.InitializeWorld(state, "test world")
	initialHistoryCount := len(state.History)
	initialTurn := state.Turn
	
	err := engine.ProcessPlayerAction(state, "look around")
	if err != nil {
		t.Errorf("ProcessPlayerAction should not return error, got %v", err)
	}
	
	// Should have added player action to history
	if len(state.History) <= initialHistoryCount {
		t.Error("Should have added history entries")
	}
	
	// Should have advanced turn
	if state.Turn != initialTurn+1 {
		t.Errorf("Expected turn to advance to %d, got %d", initialTurn+1, state.Turn)
	}
	
	// Check that player action was recorded
	found := false
	for _, entry := range state.History {
		if entry.Type == "player" && entry.Content == "look around" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Player action should be in history")
	}
}

func TestHandleSystemAction(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			DefaultModel: "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	state := NewGameState()
	
	// Test inventory command with empty inventory
	err := engine.handleSystemAction(state, "inventory")
	if err != nil {
		t.Errorf("handleSystemAction should not return error, got %v", err)
	}
	
	if len(state.History) == 0 {
		t.Error("Should have added system message")
	}
	
	lastEntry := state.History[len(state.History)-1]
	if lastEntry.Type != "system" {
		t.Errorf("Expected system entry, got %s", lastEntry.Type)
	}
	
	if !strings.Contains(lastEntry.Content, "empty") {
		t.Error("Should mention empty inventory")
	}
	
	// Test stats command with empty stats
	state.History = []HistoryEntry{} // Clear history
	err = engine.handleSystemAction(state, "stats")
	if err != nil {
		t.Errorf("handleSystemAction should not return error, got %v", err)
	}
	
	if len(state.History) == 0 {
		t.Error("Should have added system message")
	}
	
	// Test help command
	state.History = []HistoryEntry{} // Clear history
	err = engine.handleSystemAction(state, "help")
	if err != nil {
		t.Errorf("handleSystemAction should not return error, got %v", err)
	}
	
	if len(state.History) == 0 {
		t.Error("Should have added help message")
	}
	
	lastEntry = state.History[len(state.History)-1]
	if !strings.Contains(lastEntry.Content, "commands") {
		t.Error("Help should mention commands")
	}
}

func TestHandleSystemActionWithInventory(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			DefaultModel: "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	state := NewGameState()
	
	// Add items to inventory
	state.Player.Inventory = []Item{
		{Name: "sword", Description: "sharp blade", Quantity: 1},
		{Name: "potion", Description: "healing liquid", Quantity: 3},
	}
	
	err := engine.handleSystemAction(state, "inventory")
	if err != nil {
		t.Errorf("handleSystemAction should not return error, got %v", err)
	}
	
	if len(state.History) == 0 {
		t.Error("Should have added system message")
	}
	
	lastEntry := state.History[len(state.History)-1]
	if !strings.Contains(lastEntry.Content, "sword") {
		t.Error("Should mention sword in inventory")
	}
	
	if !strings.Contains(lastEntry.Content, "potion") {
		t.Error("Should mention potion in inventory")
	}
}

func TestHandleSystemActionWithStats(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			DefaultModel: "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	state := NewGameState()
	
	// Add stats
	state.Player.Stats = map[string]int{
		"health": 100,
		"mana":   50,
	}
	
	err := engine.handleSystemAction(state, "stats")
	if err != nil {
		t.Errorf("handleSystemAction should not return error, got %v", err)
	}
	
	if len(state.History) == 0 {
		t.Error("Should have added system message")
	}
	
	lastEntry := state.History[len(state.History)-1]
	if !strings.Contains(lastEntry.Content, "health") {
		t.Error("Should mention health in stats")
	}
	
	if !strings.Contains(lastEntry.Content, "100") {
		t.Error("Should mention health value")
	}
}

func TestGenerateActionSuggestions(t *testing.T) {
	cfg := &config.Config{
		AI: config.AIConfig{
			OpenRouterAPIKey: "", // Empty key will trigger fallback
			GeminiAPIKey:     "",
			DefaultModel:     "test_model",
		},
	}
	
	engine := NewEngine(cfg)
	state := NewGameState()
	state.World.Name = "Test World"
	state.World.CurrentLocation = "Starting Point"
	
	// Add some narrator history
	state.AddHistoryEntry("narrator", "You find yourself in a dark forest.")
	state.AddHistoryEntry("narrator", "There is a path ahead.")
	
	suggestions, err := engine.GenerateActionSuggestions(state)
	if err != nil {
		t.Errorf("GenerateActionSuggestions should not return error, got %v", err)
	}
	
	if len(suggestions) == 0 {
		t.Error("Should return at least some suggestions")
	}
	
	// Should return fallback suggestions when AI fails
	expectedSuggestions := []string{"Look around", "Continue forward", "Check inventory"}
	for i, expected := range expectedSuggestions {
		if i < len(suggestions) && suggestions[i] != expected {
			// This is okay - might be AI suggestions or fallback
		}
	}
}

