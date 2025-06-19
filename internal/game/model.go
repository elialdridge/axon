package game

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"axon/internal/config"
	"axon/internal/logger"
	"axon/internal/storage"
	"axon/internal/ui"
)

// GameMode represents the current mode of the game
type GameMode int

const (
	ModeMainMenu GameMode = iota
	ModeWorldSetup
	ModePlaying
	ModeSettings
	ModeSaveLoad
)

// Model represents the main game model for Bubble Tea
type Model struct {
	// Configuration
	config *config.Config
	// UI styles
	styles *ui.Styles
	// Game engine
	engine *Engine
	// Storage
	storage *storage.Storage
	// Current game state
	gameState *GameState
	// UI state
	mode         GameMode
	inputValue   string
	scrollOffset int
	width        int
	height       int
	// Action suggestions
	suggestions []string
	// Error message
	errorMessage string
	// Loading state
	isLoading bool
}

// NewModel creates a new game model
func NewModel(cfg *config.Config) *Model {
	return &Model{
		config:    cfg,
		styles:    ui.NewStyles(),
		engine:    NewEngine(cfg),
		storage:   storage.NewStorage(cfg.Game.SaveDir),
		gameState: NewGameState(),
		mode:      ModeMainMenu,
		width:     cfg.Terminal.Width,
		height:    cfg.Terminal.Height,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

// handleKeyPress handles key press events
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		if m.mode == ModePlaying {
			// Allow quitting from game with confirmation
			m.mode = ModeMainMenu
			return m, nil
		}
		return m, tea.Quit

	case "enter":
		return m.handleEnter()

	case "backspace":
		if len(m.inputValue) > 0 {
			m.inputValue = m.inputValue[:len(m.inputValue)-1]
		}
		return m, nil

	case "up":
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		return m, nil

	case "down":
		m.scrollOffset++
		return m, nil

	default:
		// Add character to input
		if len(msg.String()) == 1 {
			m.inputValue += msg.String()
			// Clear error message when user starts typing
			if m.errorMessage != "" {
				m.errorMessage = ""
			}
		}
		return m, nil
	}
}

// handleEnter handles enter key press based on current mode
func (m Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeMainMenu:
		return m.handleMainMenuSelection()
	case ModeWorldSetup:
		return m.handleWorldSetup()
	case ModePlaying:
		return m.handleGameAction()
	default:
		return m, nil
	}
}

// handleMainMenuSelection handles main menu selections
func (m Model) handleMainMenuSelection() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(strings.ToLower(m.inputValue))
	m.inputValue = ""

	switch input {
	case "1", "new", "new game":
		m.mode = ModeWorldSetup
		m.gameState = NewGameState()
	case "2", "load", "load game":
		m.mode = ModeSaveLoad
	case "3", "settings":
		m.mode = ModeSettings
	case "4", "quit", "exit":
		return m, tea.Quit
	default:
		m.errorMessage = "Invalid selection. Choose 1-4."
	}

	return m, nil
}

// handleWorldSetup handles world setup input
func (m Model) handleWorldSetup() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.inputValue)
	logger.Info("World setup input received: %s", input)
	m.inputValue = ""

	if input == "" {
		logger.Debug("Empty world description provided")
		m.errorMessage = "Please enter a world description."
		return m, nil
	}

	// Initialize world with AI
	logger.Info("Starting world initialization process")
	m.isLoading = true
	err := m.engine.InitializeWorld(m.gameState, input)
	m.isLoading = false
	logger.Info("World initialization process completed")

	if err != nil {
		logger.Error("World initialization failed: %v", err)
		m.errorMessage = fmt.Sprintf("Error creating world: %v", err)
		return m, nil
	}

	// Generate initial action suggestions
	logger.Info("Generating initial action suggestions")
	suggestions, _ := m.engine.GenerateActionSuggestions(m.gameState)
	m.suggestions = suggestions
	logger.Debug("Generated suggestions: %v", suggestions)

	m.mode = ModePlaying
	m.scrollOffset = 0
	logger.Info("Switched to playing mode")
	return m, nil
}

// handleGameAction handles game actions during play
func (m Model) handleGameAction() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.inputValue)
	logger.Info("Game action received: %s", input)
	m.inputValue = ""

	if input == "" {
		logger.Debug("Empty action received")
		return m, nil
	}

	// Handle special commands
	if strings.HasPrefix(strings.ToLower(input), "save") {
		parts := strings.SplitN(input, " ", 2)
		saveName := ""
		if len(parts) > 1 {
			saveName = parts[1]
		}
		err := m.storage.SaveGame(saveName, m.gameState)
		if err != nil {
			m.errorMessage = fmt.Sprintf("Error saving game: %v", err)
		} else {
			m.gameState.AddHistoryEntry("system", "Game saved successfully.")
		}
		return m, nil
	}

	if strings.HasPrefix(strings.ToLower(input), "load") {
		parts := strings.SplitN(input, " ", 2)
		if len(parts) < 2 {
			m.errorMessage = "Please specify a save name."
			return m, nil
		}
		saveName := parts[1]
		var loadedState GameState
		err := m.storage.LoadGame(saveName, &loadedState)
		if err != nil {
			m.errorMessage = fmt.Sprintf("Error loading game: %v", err)
		} else {
			m.gameState = &loadedState
			m.gameState.AddHistoryEntry("system", "Game loaded successfully.")
		}
		return m, nil
	}

	// Process normal game action
	logger.Info("Processing normal game action")
	m.isLoading = true
	err := m.engine.ProcessPlayerAction(m.gameState, input)
	m.isLoading = false
	logger.Info("Game action processing completed")

	if err != nil {
		logger.Error("Game action processing failed: %v", err)
		m.errorMessage = fmt.Sprintf("Error processing action: %v", err)
		return m, nil
	}

	// Generate new action suggestions
	logger.Debug("Generating new action suggestions")
	suggestions, _ := m.engine.GenerateActionSuggestions(m.gameState)
	m.suggestions = suggestions
	logger.Debug("New suggestions: %v", suggestions)

	// Auto-scroll to show latest entries
	m.scrollOffset = -1 // Use -1 to indicate we want to show the latest

	return m, nil
}

// View renders the current view
func (m Model) View() string {
	switch m.mode {
	case ModeMainMenu:
		return m.renderMainMenu()
	case ModeWorldSetup:
		return m.renderWorldSetup()
	case ModePlaying:
		return m.renderGame()
	case ModeSettings:
		return m.renderSettings()
	case ModeSaveLoad:
		return m.renderSaveLoad()
	default:
		return "Unknown mode"
	}
}

// renderMainMenu renders the main menu
func (m Model) renderMainMenu() string {
	menu := fmt.Sprintf(`AXON - AI-Driven Adventure Game

1. New Game
2. Load Game
3. Settings
4. Quit

Enter your choice: %s`, m.inputValue)

	if m.errorMessage != "" {
		menu += "\n\nError: " + m.errorMessage
	}

	return menu
}

// renderWorldSetup renders the world setup screen
func (m Model) renderWorldSetup() string {
	setup := fmt.Sprintf(`WORLD SETUP

Describe the world you want to explore.
Be creative! Examples:
- A cyberpunk city in 2077
- A medieval fantasy kingdom
- A space station on the edge of known space
- A post-apocalyptic wasteland

Your world: %s`, m.inputValue)

	if m.isLoading {
		setup += "\n\nCreating world..."
	}

	if m.errorMessage != "" {
		setup += "\n\nError: " + m.errorMessage
	}

	return setup
}

// renderGame renders the main game interface
func (m Model) renderGame() string {
	// Calculate available space
	inputHeight := 4 // Space for input panel
	historyHeight := m.height - inputHeight

	// Render history panel
	historyContent := m.renderHistory(historyHeight)

	// Render input panel
	inputContent := m.renderInput()

	// Simple ASCII separator line
	separator := strings.Repeat("-", m.width)

	return historyContent + "\n" + separator + "\n" + inputContent
}

// renderHistory renders the game history
func (m Model) renderHistory(height int) string {
	if len(m.gameState.History) == 0 {
		return m.wrapText("Welcome to Axon! Your adventure begins...")
	}

	// Calculate which entries to show based on scroll offset
	maxLines := height - 2 // Leave space for borders
	history := m.gameState.History

	// Build display lines with simple ASCII formatting and text wrapping
	var lines []string
	for _, entry := range history {
		var formattedContent string
		switch entry.Type {
		case "player":
			formattedContent = "> " + entry.Content
		case "narrator":
			formattedContent = entry.Content
		case "system":
			formattedContent = "[System] " + entry.Content
		}
		// Wrap text to terminal width
		wrappedLines := m.wrapTextToLines(formattedContent)
		lines = append(lines, wrappedLines...)
	}

	// Handle scrolling - show most recent entries by default
	startLine := 0
	if len(lines) > maxLines {
		// Show most recent entries by default
		startLine = len(lines) - maxLines
		// Allow manual scrolling
		if m.scrollOffset >= 0 && m.scrollOffset <= len(lines)-maxLines {
			startLine = m.scrollOffset
		}
	}

	endLine := startLine + maxLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Ensure we don't go out of bounds
	if startLine < 0 {
		startLine = 0
	}
	if startLine >= len(lines) {
		startLine = len(lines) - 1
		if startLine < 0 {
			startLine = 0
		}
	}

	visibleLines := lines[startLine:endLine]
	return strings.Join(visibleLines, "\n")
}

// renderInput renders the input panel
func (m Model) renderInput() string {
	var content strings.Builder
	content.WriteString("> " + m.inputValue)

	if m.isLoading {
		content.WriteString("\nProcessing...")
	}

	if m.errorMessage != "" {
		content.WriteString("\nError: " + m.wrapText(m.errorMessage))
	}

	// Show action suggestions with text wrapping
	if len(m.suggestions) > 0 && !m.isLoading {
		suggestionText := "Suggestions: " + strings.Join(m.suggestions, ", ")
		content.WriteString("\n" + m.wrapText(suggestionText))
	}

	return content.String()
}

// renderSettings renders the settings screen
func (m Model) renderSettings() string {
	return "Settings screen - Press 'q' to return to main menu"
}

// renderSaveLoad renders the save/load screen
func (m Model) renderSaveLoad() string {
	saves, _ := m.storage.ListSaves()
	content := "Available saves:\n\n"
	for _, save := range saves {
		content += "- " + save + "\n"
	}
	content += "\nPress 'q' to return to main menu"
	return m.wrapText(content)
}

// wrapText wraps text to fit terminal width
func (m Model) wrapText(text string) string {
	if m.width <= 0 {
		return text
	}
	return m.wrapTextToWidth(text, m.width-2) // Leave 2 chars for padding
}

// wrapTextToLines wraps text and returns as slice of lines
func (m Model) wrapTextToLines(text string) []string {
	if m.width <= 0 {
		return []string{text}
	}
	return m.wrapTextToLinesWidth(text, m.width-2) // Leave 2 chars for padding
}

// wrapTextToWidth wraps text to specified width
func (m Model) wrapTextToWidth(text string, width int) string {
	if width <= 0 {
		return text
	}

	lines := m.wrapTextToLinesWidth(text, width)
	return strings.Join(lines, "\n")
}

// wrapTextToLinesWidth wraps text to specified width and returns lines
func (m Model) wrapTextToLinesWidth(text string, width int) []string {
	if width <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		// Check if adding this word would exceed width
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			// Start a new line
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}
		
		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	// Return at least one line
	if len(lines) == 0 {
		return []string{""}
	}

	return lines
}
