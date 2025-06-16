package game

import (
	"fmt"
	"testing"
	"time"
)

func TestNewGameState(t *testing.T) {
	state := NewGameState()
	
	if state == nil {
		t.Fatal("NewGameState returned nil")
	}
	
	// Test initial values
	if state.Turn != 0 {
		t.Errorf("Expected initial turn 0, got %d", state.Turn)
	}
	
	if len(state.History) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(state.History))
	}
	
	if state.World == nil {
		t.Fatal("World should not be nil")
	}
	
	if state.Player == nil {
		t.Fatal("Player should not be nil")
	}
	
	// Test world initialization
	if state.World.Locations == nil {
		t.Error("World locations should be initialized")
	}
	
	// Test player initialization
	if state.Player.Inventory == nil {
		t.Error("Player inventory should be initialized")
	}
	
	if state.Player.Stats == nil {
		t.Error("Player stats should be initialized")
	}
	
	// Test timestamps
	if state.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	
	if state.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestAddHistoryEntry(t *testing.T) {
	state := NewGameState()
	initialTime := state.UpdatedAt
	
	// Add history entry
	time.Sleep(1 * time.Millisecond) // Ensure time difference
	state.AddHistoryEntry("player", "test action")
	
	// Check history was added
	if len(state.History) != 1 {
		t.Errorf("Expected 1 history entry, got %d", len(state.History))
	}
	
	entry := state.History[0]
	if entry.Type != "player" {
		t.Errorf("Expected type 'player', got %s", entry.Type)
	}
	
	if entry.Content != "test action" {
		t.Errorf("Expected content 'test action', got %s", entry.Content)
	}
	
	if entry.Turn != 0 {
		t.Errorf("Expected turn 0, got %d", entry.Turn)
	}
	
	// Check timestamp was set
	if entry.Timestamp.IsZero() {
		t.Error("Entry timestamp should be set")
	}
	
	// Check updated time changed
	if !state.UpdatedAt.After(initialTime) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestGetRecentHistory(t *testing.T) {
	state := NewGameState()
	
	// Add multiple history entries
	for i := 0; i < 5; i++ {
		state.AddHistoryEntry("system", fmt.Sprintf("entry %d", i))
	}
	
	// Test getting more entries than exist
	recent := state.GetRecentHistory(10)
	if len(recent) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(recent))
	}
	
	// Test getting fewer entries
	recent = state.GetRecentHistory(3)
	if len(recent) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(recent))
	}
	
	// Check it returns the last entries
	if recent[0].Content != "entry 2" {
		t.Errorf("Expected first recent entry to be 'entry 2', got %s", recent[0].Content)
	}
	
	if recent[2].Content != "entry 4" {
		t.Errorf("Expected last recent entry to be 'entry 4', got %s", recent[2].Content)
	}
}

func TestNextTurn(t *testing.T) {
	state := NewGameState()
	initialTime := state.UpdatedAt
	
	time.Sleep(1 * time.Millisecond) // Ensure time difference
	state.NextTurn()
	
	if state.Turn != 1 {
		t.Errorf("Expected turn 1, got %d", state.Turn)
	}
	
	if !state.UpdatedAt.After(initialTime) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestWorldStruct(t *testing.T) {
	world := &World{
		Name:        "Test World",
		Description: "A test world",
		Setting:     "Fantasy",
		Rules:       []string{"rule1", "rule2"},
		Locations:   map[string]string{"start": "Starting location"},
		CurrentLocation: "start",
	}
	
	if world.Name != "Test World" {
		t.Errorf("Expected name 'Test World', got %s", world.Name)
	}
	
	if len(world.Rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(world.Rules))
	}
	
	if len(world.Locations) != 1 {
		t.Errorf("Expected 1 location, got %d", len(world.Locations))
	}
}

func TestPlayerStruct(t *testing.T) {
	player := &Player{
		Name:        "Test Player",
		Description: "A test player",
		Inventory:   []Item{{Name: "sword", Description: "sharp", Quantity: 1}},
		Stats:       map[string]int{"health": 100, "mana": 50},
		Status:      "healthy",
	}
	
	if player.Name != "Test Player" {
		t.Errorf("Expected name 'Test Player', got %s", player.Name)
	}
	
	if len(player.Inventory) != 1 {
		t.Errorf("Expected 1 item, got %d", len(player.Inventory))
	}
	
	if player.Inventory[0].Name != "sword" {
		t.Errorf("Expected item name 'sword', got %s", player.Inventory[0].Name)
	}
	
	if len(player.Stats) != 2 {
		t.Errorf("Expected 2 stats, got %d", len(player.Stats))
	}
	
	if player.Stats["health"] != 100 {
		t.Errorf("Expected health 100, got %d", player.Stats["health"])
	}
}

func TestItemStruct(t *testing.T) {
	item := Item{
		Name:        "Health Potion",
		Description: "Restores health",
		Quantity:    3,
	}
	
	if item.Name != "Health Potion" {
		t.Errorf("Expected name 'Health Potion', got %s", item.Name)
	}
	
	if item.Quantity != 3 {
		t.Errorf("Expected quantity 3, got %d", item.Quantity)
	}
}

func TestHistoryEntryStruct(t *testing.T) {
	now := time.Now()
	entry := HistoryEntry{
		Type:      "narrator",
		Content:   "You see a door",
		Timestamp: now,
		Turn:      5,
	}
	
	if entry.Type != "narrator" {
		t.Errorf("Expected type 'narrator', got %s", entry.Type)
	}
	
	if entry.Content != "You see a door" {
		t.Errorf("Expected content 'You see a door', got %s", entry.Content)
	}
	
	if entry.Turn != 5 {
		t.Errorf("Expected turn 5, got %d", entry.Turn)
	}
	
	if !entry.Timestamp.Equal(now) {
		t.Error("Timestamp should match")
	}
}

