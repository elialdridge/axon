package game

import (
	"fmt"
	"strings"

	"axon/internal/ai"
	"axon/internal/config"
)

// Engine represents the game engine
type Engine struct {
	aiClient *ai.Client
	config   *config.Config
}

// NewEngine creates a new game engine
func NewEngine(cfg *config.Config) *Engine {
	aiClient := ai.NewClient(cfg.AI.OpenRouterAPIKey, cfg.AI.GeminiAPIKey)
	return &Engine{
		aiClient: aiClient,
		config:   cfg,
	}
}

// InitializeWorld creates the initial game world based on a seed prompt
func (e *Engine) InitializeWorld(state *GameState, seedPrompt string) error {
	// Use the best model for world building
	model := e.aiClient.GetBestModel("world_building")

	// Create context for world generation
	context := []string{
		"You are a master world builder for a text-based adventure game.",
		"Create a detailed, immersive world based on the user's prompt.",
		"Respond with JSON in this format:",
		`{"name": "World Name", "description": "Detailed description", "setting": "Genre/Setting", "rules": ["rule1", "rule2"], "starting_location": "location name", "starting_location_desc": "description"}`,
	}

	prompt := fmt.Sprintf("Create a world for this scenario: %s", seedPrompt)

	req := ai.Request{
		Prompt:    prompt,
		Model:     model,
		MaxTokens: 1000,
		Context:   context,
	}

	resp, err := e.aiClient.Generate(req)
	if err != nil {
		return fmt.Errorf("failed to generate world: %w", err)
	}

	if resp.Error != nil {
		// Fallback to basic world creation
		state.World.Name = "Unknown Realm"
		state.World.Description = "A mysterious world awaits your exploration."
		state.World.Setting = "Fantasy"
		state.World.Rules = []string{"Anything is possible", "Actions have consequences"}
		state.World.CurrentLocation = "Starting Point"
		state.World.Locations["Starting Point"] = "You find yourself in an unknown place, ready to begin your adventure."
		state.AddHistoryEntry("system", "World creation failed, using default world.")
	} else {
		// Parse AI response and populate world (simplified - in practice you'd parse JSON)
		state.World.Name = "AI Generated World"
		state.World.Description = resp.Text
		state.World.Setting = "AI Generated"
		state.World.Rules = []string{"AI-driven narrative", "Player choices matter"}
		state.World.CurrentLocation = "Starting Point"
		state.World.Locations["Starting Point"] = resp.Text
		state.AddHistoryEntry("narrator", resp.Text)
	}

	return nil
}

// ProcessPlayerAction processes a player action and generates response
func (e *Engine) ProcessPlayerAction(state *GameState, action string) error {
	// Add player action to history
	state.AddHistoryEntry("player", action)

	// Get recent history for context
	recentHistory := state.GetRecentHistory(10)
	contextLines := make([]string, 0)

	// Build context from recent history
	for _, entry := range recentHistory {
		contextLines = append(contextLines, fmt.Sprintf("%s: %s", entry.Type, entry.Content))
	}

	// Choose appropriate model based on action type
	var model string

	actionLower := strings.ToLower(action)
	switch {
	case strings.Contains(actionLower, "say") || strings.Contains(actionLower, "talk") || strings.Contains(actionLower, "speak"):
		model = e.aiClient.GetBestModel("dialog")
	case strings.Contains(actionLower, "inventory") || strings.Contains(actionLower, "stats"):
		return e.handleSystemAction(state, action)
	default:
		model = e.aiClient.GetBestModel("storytelling")
	}

	// Create AI request
	context := []string{
		"You are the Game Master for a text-based adventure game.",
		fmt.Sprintf("World: %s - %s", state.World.Name, state.World.Description),
		fmt.Sprintf("Current Location: %s", state.World.CurrentLocation),
		fmt.Sprintf("Player: %s - %s", state.Player.Name, state.Player.Description),
		"Recent game history:",
		strings.Join(contextLines, "\n"),
		"Respond to the player's action with narrative description. Keep responses concise but engaging.",
	}

	prompt := fmt.Sprintf("Player action: %s", action)

	req := ai.Request{
		Prompt:    prompt,
		Model:     model,
		MaxTokens: 500,
		Context:   context,
	}

	resp, err := e.aiClient.Generate(req)
	if err != nil {
		return fmt.Errorf("failed to generate response: %w", err)
	}

	if resp.Error != nil {
		// Fallback response
		state.AddHistoryEntry("narrator", "Something happens in response to your action.")
		state.AddHistoryEntry("system", fmt.Sprintf("AI Error: %v", resp.Error))
	} else {
		state.AddHistoryEntry("narrator", resp.Text)
	}

	// Advance turn
	state.NextTurn()

	return nil
}

// handleSystemAction handles system actions like inventory, stats, etc.
func (e *Engine) handleSystemAction(state *GameState, action string) error {
	actionLower := strings.ToLower(action)

	switch {
	case strings.Contains(actionLower, "inventory"):
		if len(state.Player.Inventory) == 0 {
			state.AddHistoryEntry("system", "Your inventory is empty.")
		} else {
			inventoryList := "Inventory:\n"
			for _, item := range state.Player.Inventory {
				inventoryList += fmt.Sprintf("- %s (x%d): %s\n", item.Name, item.Quantity, item.Description)
			}
			state.AddHistoryEntry("system", inventoryList)
		}

	case strings.Contains(actionLower, "stats"):
		if len(state.Player.Stats) == 0 {
			state.AddHistoryEntry("system", "No stats to display.")
		} else {
			statsList := "Character Stats:\n"
			for stat, value := range state.Player.Stats {
				statsList += fmt.Sprintf("- %s: %d\n", stat, value)
			}
			state.AddHistoryEntry("system", statsList)
		}

	case strings.Contains(actionLower, "help"):
		helpText := `Available commands:
- Type any action to interact with the world
- 'inventory' or 'inv' to check your items
- 'stats' to view character statistics
- 'save [name]' to save your game
- 'load [name]' to load a saved game
- 'quit' to exit the game`
		state.AddHistoryEntry("system", helpText)

	default:
		state.AddHistoryEntry("system", "Unknown system command. Type 'help' for available commands.")
	}

	return nil
}

// GenerateActionSuggestions generates suggested actions for the player
func (e *Engine) GenerateActionSuggestions(state *GameState) ([]string, error) {
	model := e.aiClient.GetBestModel("rule_setting")

	context := []string{
		"Generate 3-4 brief action suggestions for the player in this situation.",
		fmt.Sprintf("World: %s", state.World.Name),
		fmt.Sprintf("Location: %s", state.World.CurrentLocation),
		"Provide only the action suggestions, one per line, without numbers or bullets.",
	}

	// Get recent context
	recentHistory := state.GetRecentHistory(3)
	contextLines := make([]string, 0)
	for _, entry := range recentHistory {
		if entry.Type == "narrator" {
			contextLines = append(contextLines, entry.Content)
		}
	}

	prompt := "Current situation: " + strings.Join(contextLines, " ")

	req := ai.Request{
		Prompt:    prompt,
		Model:     model,
		MaxTokens: 200,
		Context:   context,
	}

	resp, err := e.aiClient.Generate(req)
	if err != nil {
		return []string{"Look around", "Continue forward", "Check inventory"}, nil
	}

	if resp.Error != nil {
		return []string{"Look around", "Continue forward", "Check inventory"}, nil
	}

	// Parse suggestions from response
	suggestions := strings.Split(strings.TrimSpace(resp.Text), "\n")
	if len(suggestions) == 0 {
		return []string{"Look around", "Continue forward", "Check inventory"}, nil
	}

	// Clean up suggestions
	cleanSuggestions := make([]string, 0)
	for _, suggestion := range suggestions {
		clean := strings.TrimSpace(suggestion)
		if clean != "" {
			cleanSuggestions = append(cleanSuggestions, clean)
		}
	}

	return cleanSuggestions, nil
}
