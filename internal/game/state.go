package game

import (
	"time"
)

// GameState represents the current state of the game
type GameState struct {
	// World information
	World *World `json:"world"`
	// Player information
	Player *Player `json:"player"`
	// Game history
	History []HistoryEntry `json:"history"`
	// Current turn
	Turn int `json:"turn"`
	// Game metadata
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// World represents the game world
type World struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Setting         string            `json:"setting"`
	Rules           []string          `json:"rules"`
	Locations       map[string]string `json:"locations"`
	CurrentLocation string            `json:"current_location"`
}

// Player represents the player character
type Player struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Inventory   []Item         `json:"inventory"`
	Stats       map[string]int `json:"stats"`
	Status      string         `json:"status"`
}

// Item represents an inventory item
type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}

// HistoryEntry represents an entry in the game history
type HistoryEntry struct {
	Type      string    `json:"type"` // "player", "narrator", "system"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
	Turn      int       `json:"turn"`
}

// NewGameState creates a new game state
func NewGameState() *GameState {
	now := time.Now()
	return &GameState{
		World: &World{
			Locations: make(map[string]string),
		},
		Player: &Player{
			Inventory: make([]Item, 0),
			Stats:     make(map[string]int),
		},
		History:   make([]HistoryEntry, 0),
		Turn:      0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddHistoryEntry adds an entry to the game history
func (gs *GameState) AddHistoryEntry(entryType, content string) {
	entry := HistoryEntry{
		Type:      entryType,
		Content:   content,
		Timestamp: time.Now(),
		Turn:      gs.Turn,
	}
	gs.History = append(gs.History, entry)
	gs.UpdatedAt = time.Now()
}

// GetRecentHistory returns the last N history entries
func (gs *GameState) GetRecentHistory(n int) []HistoryEntry {
	if len(gs.History) <= n {
		return gs.History
	}
	return gs.History[len(gs.History)-n:]
}

// NextTurn advances to the next turn
func (gs *GameState) NextTurn() {
	gs.Turn++
	gs.UpdatedAt = time.Now()
}
