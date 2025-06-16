package game

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"axon/internal/config"
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
	m.inputValue = ""

	if input == "" {
		m.errorMessage = "Please enter a world description."
		return m, nil
	}

	// Initialize world with AI
	m.isLoading = true
	err := m.engine.InitializeWorld(m.gameState, input)
	m.isLoading = false

	if err != nil {
		m.errorMessage = fmt.Sprintf("Error creating world: %v", err)
		return m, nil
	}

	// Generate initial action suggestions
	suggestions, _ := m.engine.GenerateActionSuggestions(m.gameState)
	m.suggestions = suggestions

	m.mode = ModePlaying
	m.scrollOffset = 0
	return m, nil
}

// handleGameAction handles game actions during play
func (m Model) handleGameAction() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.inputValue)
	m.inputValue = ""

	if input == "" {
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
	m.isLoading = true
	err := m.engine.ProcessPlayerAction(m.gameState, input)
	m.isLoading = false

	if err != nil {
		m.errorMessage = fmt.Sprintf("Error processing action: %v", err)
		return m, nil
	}

	// Generate new action suggestions
	suggestions, _ := m.engine.GenerateActionSuggestions(m.gameState)
	m.suggestions = suggestions

	// Auto-scroll to bottom
	m.scrollOffset = len(m.gameState.History)

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
	menu := fmt.Sprintf(`%s

AXON - AI-Driven Adventure Game

1. New Game
2. Load Game
3. Settings
4. Quit

Enter your choice: %s`,
		m.styles.Base.Render(""),
		m.styles.PlayerText.Render(m.inputValue))

	if m.errorMessage != "" {
		menu += "\n\n" + m.styles.SystemText.Render("Error: "+m.errorMessage)
		m.errorMessage = ""
	}

	return m.styles.Base.Width(m.width).Height(m.height).Render(menu)
}

// renderWorldSetup renders the world setup screen
func (m Model) renderWorldSetup() string {
	setup := fmt.Sprintf(`%s

WORLD SETUP

Describe the world you want to explore.
Be creative! Examples:
- A cyberpunk city in 2077
- A medieval fantasy kingdom
- A space station on the edge of known space
- A post-apocalyptic wasteland

Your world: %s`,
		m.styles.Base.Render(""),
		m.styles.PlayerText.Render(m.inputValue))

	if m.isLoading {
		setup += "\n\n" + m.styles.SystemText.Render("Creating world...")
	}

	if m.errorMessage != "" {
		setup += "\n\n" + m.styles.SystemText.Render("Error: "+m.errorMessage)
		m.errorMessage = ""
	}

	return m.styles.Base.Width(m.width).Height(m.height).Render(setup)
}

// renderGame renders the main game interface
func (m Model) renderGame() string {
	// Calculate available space
	inputHeight := 4 // Space for input panel
	historyHeight := m.height - inputHeight

	// Render history panel
	historyContent := m.renderHistory(historyHeight)
	historyPanel := m.styles.HistoryPanel.Width(m.width).Height(historyHeight).Render(historyContent)

	// Render input panel
	inputContent := m.renderInput()
	inputPanel := m.styles.InputPanel.Width(m.width).Height(inputHeight).Render(inputContent)

	return lipgloss.JoinVertical(lipgloss.Left, historyPanel, inputPanel)
}

// renderHistory renders the game history
func (m Model) renderHistory(height int) string {
	if len(m.gameState.History) == 0 {
		return "Welcome to Axon! Your adventure begins..."
	}

	// Calculate which entries to show based on scroll offset
	maxLines := height - 2 // Leave space for borders
	history := m.gameState.History

	// Build display lines
	var lines []string
	for _, entry := range history {
		var style lipgloss.Style
		switch entry.Type {
		case "player":
			style = m.styles.PlayerText
			lines = append(lines, style.Render("> "+entry.Content))
		case "narrator":
			style = m.styles.NarratorText
			lines = append(lines, style.Render(entry.Content))
		case "system":
			style = m.styles.SystemText
			lines = append(lines, style.Render("[System] "+entry.Content))
		}
	}

	// Handle scrolling
	startLine := 0
	if len(lines) > maxLines {
		startLine = len(lines) - maxLines
		if m.scrollOffset < len(lines)-maxLines {
			startLine = m.scrollOffset
		}
	}

	endLine := startLine + maxLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	visibleLines := lines[startLine:endLine]
	return strings.Join(visibleLines, "\n")
}

// renderInput renders the input panel
func (m Model) renderInput() string {
	prompt := m.styles.Prompt.Render("> ")
	input := m.styles.PlayerText.Render(m.inputValue)

	var content strings.Builder
	content.WriteString(prompt + input)

	if m.isLoading {
		content.WriteString("\n" + m.styles.SystemText.Render("Processing..."))
	}

	if m.errorMessage != "" {
		content.WriteString("\n" + m.styles.SystemText.Render("Error: "+m.errorMessage))
		m.errorMessage = ""
	}

	// Show action suggestions
	if len(m.suggestions) > 0 && !m.isLoading {
		content.WriteString("\n" + m.styles.SystemText.Render("Suggestions: "+strings.Join(m.suggestions, ", ")))
	}

	return content.String()
}

// renderSettings renders the settings screen
func (m Model) renderSettings() string {
	return m.styles.Base.Width(m.width).Height(m.height).Render("Settings screen - Press 'q' to return to main menu")
}

// renderSaveLoad renders the save/load screen
func (m Model) renderSaveLoad() string {
	saves, _ := m.storage.ListSaves()
	content := "Available saves:\n\n"
	for _, save := range saves {
		content += "- " + save + "\n"
	}
	content += "\nPress 'q' to return to main menu"
	return m.styles.Base.Width(m.width).Height(m.height).Render(content)
}

