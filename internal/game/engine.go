package game

import (
	"fmt"
	"strings"

	"axon/internal/ai"
	"axon/internal/config"
	"axon/internal/logger"
)

const (
	// String constants for entry types
	entryTypePlayer   = "player"
	entryTypeNarrator = "narrator"
	entryTypeSystem   = "system"
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
	logger.Info("Starting world initialization with prompt: %s", seedPrompt)
	logger.LogWorldCreation("start", seedPrompt)

	// Use the best model for world building
	model := e.aiClient.GetBestModel("world_building")
	logger.Debug("Selected model for world building: %s", model)

	// Create simplified context for world generation
	context := []string{
		"You are creating a world for a text-based adventure game.",
		"Create a brief, engaging world description based on the user's prompt.",
		"Keep your response concise and immersive.",
	}

	prompt := fmt.Sprintf("Create a world: %s. Describe the setting in 2-3 sentences.", seedPrompt)
	logger.Debug("World creation prompt: %s", prompt)
	logger.LogWorldCreation("context", context)

	req := ai.Request{
		Prompt:    prompt,
		Model:     model,
		MaxTokens: 200, // Reduced for faster response
		Context:   context,
	}

	logger.LogWorldCreation("request", req)

	logger.Info("Sending world creation request to AI")
	resp, err := e.aiClient.Generate(req)
	if err != nil {
		logger.Error("AI world generation failed: %v", err)
		return fmt.Errorf("failed to generate world: %w", err)
	}

	logger.LogWorldCreation("response", resp)

	if resp.Error != nil {
		logger.Error("AI response contains error: %v", resp.Error)
		logger.LogWorldCreation("fallback", "using themed world based on prompt")
		// Create themed fallback world based on the seed prompt
		themeWorld := e.createThemedWorld(seedPrompt)
		state.World.Name = themeWorld.Name
		state.World.Description = themeWorld.Description
		state.World.Setting = themeWorld.Setting
		state.World.Rules = themeWorld.Rules
		state.World.CurrentLocation = themeWorld.CurrentLocation
		state.World.Locations[themeWorld.CurrentLocation] = themeWorld.Description
		state.AddHistoryEntry(entryTypeNarrator, themeWorld.Description)
		state.AddHistoryEntry(
			entryTypeNarrator,
			"The mists of reality shimmer and coalesce, drawing upon ancient memories and forgotten tales to weave this world into existence...",
		)
	} else {
		logger.Info("AI world creation successful, parsing response")
		logger.LogWorldCreation("ai_success", resp.Text)
		// Parse AI response and populate world
		state.World.Name = "Generated World"
		state.World.Description = resp.Text
		state.World.Setting = "AI Generated"
		state.World.Rules = []string{"AI-driven narrative", "Player choices matter"}
		state.World.CurrentLocation = "Starting Point"
		state.World.Locations["Starting Point"] = resp.Text
		state.AddHistoryEntry(entryTypeNarrator, resp.Text)
	}

	logger.LogGameState(state)
	logger.Info("World initialization completed successfully")
	return nil
}

// ProcessPlayerAction processes a player action and generates response
func (e *Engine) ProcessPlayerAction(state *GameState, action string) error {
	logger.Info("Processing player action: %s", action)
	// Add player action to history
	state.AddHistoryEntry(entryTypePlayer, action)

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

	logger.Info("Sending action request to AI")
	resp, err := e.aiClient.Generate(req)
	if err != nil {
		logger.Error("AI action processing failed: %v", err)
		return fmt.Errorf("failed to generate response: %w", err)
	}

	if resp.Error != nil {
		logger.Error("AI response contains error: %v", resp.Error)
		// Immersive fallback response based on action type
		fallbackResponse := e.generateFallbackResponse(action, state)
		state.AddHistoryEntry(entryTypeNarrator, fallbackResponse)
	} else {
		logger.Info("AI action processing successful")
		logger.Debug("AI response: %s", resp.Text)
		state.AddHistoryEntry(entryTypeNarrator, resp.Text)
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
			state.AddHistoryEntry(entryTypeSystem, "Your inventory is empty.")
		} else {
			inventoryList := "Inventory:\n"
			for _, item := range state.Player.Inventory {
				inventoryList += fmt.Sprintf("- %s (x%d): %s\n", item.Name, item.Quantity, item.Description)
			}
			state.AddHistoryEntry(entryTypeSystem, inventoryList)
		}

	case strings.Contains(actionLower, "stats"):
		if len(state.Player.Stats) == 0 {
			state.AddHistoryEntry(entryTypeSystem, "No stats to display.")
		} else {
			statsList := "Character Stats:\n"
			for stat, value := range state.Player.Stats {
				statsList += fmt.Sprintf("- %s: %d\n", stat, value)
			}
			state.AddHistoryEntry(entryTypeSystem, statsList)
		}

	case strings.Contains(actionLower, "help"):
		helpText := `Available commands:
- Type any action to interact with the world
- 'inventory' or 'inv' to check your items
- 'stats' to view character statistics
- 'save [name]' to save your game
- 'load [name]' to load a saved game
- 'quit' to exit the game`
		state.AddHistoryEntry(entryTypeSystem, helpText)

	default:
		state.AddHistoryEntry(entryTypeSystem, "Unknown system command. Type 'help' for available commands.")
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
		if entry.Type == entryTypeNarrator {
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

// createThemedWorld creates a themed world based on user input when AI is unavailable
func (e *Engine) createThemedWorld(seedPrompt string) *World {
	lower := strings.ToLower(seedPrompt)

	// Use helper functions to detect themes and reduce complexity
	if e.isCyberpunkTheme(lower) {
		return e.createCyberpunkWorld()
	}
	if e.isFantasyTheme(lower) {
		return e.createFantasyWorld()
	}
	if e.isSpaceTheme(lower) {
		return e.createSpaceWorld()
	}
	if e.isApocalypticTheme(lower) {
		return e.createApocalypticWorld()
	}

	// Default modern/mystery theme
	return e.createDefaultWorld()
}

func (e *Engine) isCyberpunkTheme(lower string) bool {
	return strings.Contains(lower, "cyberpunk") || strings.Contains(lower, "2077") || strings.Contains(lower, "cyber")
}

func (e *Engine) isFantasyTheme(lower string) bool {
	return strings.Contains(lower, "fantasy") || strings.Contains(lower, "medieval") ||
		strings.Contains(lower, "kingdom") ||
		strings.Contains(lower, "magic")
}

func (e *Engine) isSpaceTheme(lower string) bool {
	return strings.Contains(lower, "space") || strings.Contains(lower, "station") ||
		strings.Contains(lower, "galaxy") ||
		strings.Contains(lower, "alien")
}

func (e *Engine) isApocalypticTheme(lower string) bool {
	return strings.Contains(lower, "apocalyptic") || strings.Contains(lower, "wasteland") ||
		strings.Contains(lower, "survivor") ||
		strings.Contains(lower, "ruins")
}

func (e *Engine) createCyberpunkWorld() *World {
	return &World{
		Name:        "Neo-Tokyo 2077",
		Description: "Towering neon-lit skyscrapers pierce the smoggy sky above rain-slicked streets. Corporate megastructures cast shadows over bustling markets where cybernetic implants gleam under holographic advertisements. You stand at the edge of the underground district, where rebels and hackers gather in the shadows.",
		Setting:     "Cyberpunk",
		Rules: []string{
			"Technology rules all",
			"Corporate power is absolute",
			"Information is currency",
			"Trust no one",
		},
		CurrentLocation: "Underground District",
		Locations:       make(map[string]string),
	}
}

func (e *Engine) createFantasyWorld() *World {
	return &World{
		Name:        "Realm of Eldoria",
		Description: "Ancient stone towers rise from mist-covered valleys where dragons once soared. Cobblestone paths wind through enchanted forests filled with mysterious creatures. You find yourself at the edge of a village where flickering torches cast dancing shadows on thatched roofs.",
		Setting:     "High Fantasy",
		Rules: []string{
			"Magic flows through all things",
			"Ancient powers stir",
			"Honor above all",
			"Knowledge is power",
		},
		CurrentLocation: "Village Edge",
		Locations:       make(map[string]string),
	}
}

func (e *Engine) createSpaceWorld() *World {
	return &World{
		Name:        "Frontier Station Alpha",
		Description: "The vast expanse of space stretches endlessly beyond reinforced viewports. This research station orbits a mysterious planet where strange energy readings emanate from the surface. Emergency lights flicker in the corridors as you hear the hum of life support systems working overtime.",
		Setting:     "Space Opera",
		Rules: []string{
			"The void is unforgiving",
			"Technology can fail",
			"First contact protocols exist",
			"Survival is paramount",
		},
		CurrentLocation: "Station Corridor",
		Locations:       make(map[string]string),
	}
}

func (e *Engine) createApocalypticWorld() *World {
	return &World{
		Name:            "The Shattered Lands",
		Description:     "Crumbling ruins of civilization stretch across a barren landscape under an eternally grey sky. Rusted vehicles and collapsed buildings tell the story of a world that once was. You emerge from a makeshift shelter, scanning the horizon for signs of other survivors or threats.",
		Setting:         "Post-Apocalyptic",
		Rules:           []string{"Resources are scarce", "Trust is earned", "The past is gone", "Adapt or perish"},
		CurrentLocation: "Wasteland Outpost",
		Locations:       make(map[string]string),
	}
}

func (e *Engine) createDefaultWorld() *World {
	return &World{
		Name:        "The Unknown",
		Description: "You find yourself in a place that defies easy description. Familiar yet strange, ordinary yet filled with hidden possibilities. The air itself seems to whisper of secrets waiting to be discovered and adventures yet to unfold.",
		Setting:     "Modern Mystery",
		Rules: []string{
			"Nothing is as it seems",
			"Every choice matters",
			"Mysteries abound",
			"Reality is flexible",
		},
		CurrentLocation: "Starting Point",
		Locations:       make(map[string]string),
	}
}

// generateFallbackResponse creates immersive fallback responses when AI is unavailable
func (e *Engine) generateFallbackResponse(action string, state *GameState) string {
	actionLower := strings.ToLower(action)

	// Try to match action with predefined responses
	if response := e.tryMatchActionType(actionLower); response != "" {
		return response
	}

	// Return default mystical response
	return e.getDefaultFallbackResponse(state)
}

func (e *Engine) tryMatchActionType(actionLower string) string {
	if e.isMovementAction(actionLower) {
		return "You move through the area, your footsteps echoing softly as you explore your surroundings. The path ahead remains shrouded in mystery, waiting for your next decision."
	}
	if e.isObservationAction(actionLower) {
		return "You take a moment to carefully observe your surroundings. Details emerge from the shadows - subtle signs and hidden clues that might prove important on your journey."
	}
	if e.isSearchAction(actionLower) {
		return "You search methodically, running your hands along surfaces and peering into dark corners. Though nothing immediately reveals itself, you sense that persistence might yet yield results."
	}
	if e.isCombatAction(actionLower) {
		return "Your muscles tense as you prepare for conflict. The air crackles with tension, and you feel the familiar rush of adrenaline coursing through your veins."
	}
	if e.isTakingAction(actionLower) {
		return "You reach out carefully, your fingers closing around the object. A sense of acquisition fills you as you secure this new addition to your belongings."
	}
	if e.isCommunicationAction(actionLower) {
		return "Your words hang in the air, carrying with them the weight of intention. Whether anyone is listening remains to be seen, but you have made your voice heard."
	}
	if e.isOpeningAction(actionLower) {
		return "With determined effort, you work to overcome the obstacle before you. Progress is slow but steady, and you sense that your persistence will eventually pay off."
	}
	if e.isRestingAction(actionLower) {
		return "Time passes quietly as you pause in your journey. The world continues its ancient rhythms around you, and you feel a moment of peace amidst the uncertainty."
	}
	if e.isEscapeAction(actionLower) {
		return "Your heart pounds as you move swiftly away from potential danger. The landscape blurs past you as survival instincts take over, guiding your hurried steps."
	}
	return ""
}

func (e *Engine) isMovementAction(action string) bool {
	return strings.Contains(action, "go") || strings.Contains(action, "walk") || strings.Contains(action, "move") ||
		strings.Contains(action, "head")
}

func (e *Engine) isObservationAction(action string) bool {
	return strings.Contains(action, "look") || strings.Contains(action, "examine") ||
		strings.Contains(action, "observe") ||
		strings.Contains(action, "inspect")
}

func (e *Engine) isSearchAction(action string) bool {
	return strings.Contains(action, "search") || strings.Contains(action, "find") || strings.Contains(action, "seek")
}

func (e *Engine) isCombatAction(action string) bool {
	return strings.Contains(action, "attack") || strings.Contains(action, "fight") ||
		strings.Contains(action, "strike") ||
		strings.Contains(action, "hit")
}

func (e *Engine) isTakingAction(action string) bool {
	return strings.Contains(action, "take") || strings.Contains(action, "grab") || strings.Contains(action, "pick") ||
		strings.Contains(action, "get")
}

func (e *Engine) isCommunicationAction(action string) bool {
	return strings.Contains(action, "say") || strings.Contains(action, "speak") || strings.Contains(action, "talk") ||
		strings.Contains(action, "tell")
}

func (e *Engine) isOpeningAction(action string) bool {
	return strings.Contains(action, "open") || strings.Contains(action, "unlock") || strings.Contains(action, "break")
}

func (e *Engine) isRestingAction(action string) bool {
	return strings.Contains(action, "wait") || strings.Contains(action, "rest") || strings.Contains(action, "pause") ||
		strings.Contains(action, "sit")
}

func (e *Engine) isEscapeAction(action string) bool {
	return strings.Contains(action, "run") || strings.Contains(action, "flee") || strings.Contains(action, "escape")
}

func (e *Engine) getDefaultFallbackResponse(state *GameState) string {
	fallbackResponses := []string{
		"The fabric of reality ripples slightly in response to your action, though the full consequences remain hidden in the mists of time.",
		"Something stirs in the unseen spaces around you. Your action has been noted by forces beyond immediate comprehension.",
		"The world responds to your intent in ways both subtle and profound. Change flows through the environment like water through stone.",
		"Ancient energies swirl around your action, weaving new possibilities into the tapestry of your adventure.",
		"Your deed echoes through the mysterious realm, creating ripples that will shape future moments in ways yet unknown.",
	}

	// Use turn number to pseudo-randomly select response
	responseIndex := state.Turn % len(fallbackResponses)
	return fallbackResponses[responseIndex]
}
