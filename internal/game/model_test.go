package game

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"axon/internal/config"
)

func TestNewModel(t *testing.T) {
	cfg := &config.Config{
		Terminal: config.TerminalConfig{
			Width:  80,
			Height: 24,
		},
		AI: config.AIConfig{
			DefaultModel: "test_model",
		},
		Game: config.GameConfig{
			HistoryLimit: 1000,
			SaveDir:      "/tmp/test_saves",
		},
	}
	
	model := NewModel(cfg)
	
	if model == nil {
		t.Fatal("NewModel returned nil")
	}
	
	if model.config != cfg {
		t.Error("Model config not set correctly")
	}
	
	if model.engine == nil {
		t.Error("Engine should be initialized")
	}
	
	if model.storage == nil {
		t.Error("Storage should be initialized")
	}
	
	if model.gameState == nil {
		t.Error("Game state should be initialized")
	}
	
	if model.styles == nil {
		t.Error("Styles should be initialized")
	}
	
	if model.mode != ModeMainMenu {
		t.Errorf("Expected initial mode to be ModeMainMenu, got %v", model.mode)
	}
	
	if model.width != 80 {
		t.Errorf("Expected width 80, got %d", model.width)
	}
	
	if model.height != 24 {
		t.Errorf("Expected height 24, got %d", model.height)
	}
}

func TestModelInit(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	cmd := model.Init()
	if cmd != nil {
		t.Error("Init should return nil command")
	}
}

func TestModelWindowResize(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	// Test window resize
	resizeMsg := tea.WindowSizeMsg{
		Width:  100,
		Height: 30,
	}
	
	newModel, cmd := model.Update(resizeMsg)
	if cmd != nil {
		t.Error("Window resize should not return command")
	}
	
	updatedModel := newModel.(Model)
	if updatedModel.width != 100 {
		t.Errorf("Expected width 100, got %d", updatedModel.width)
	}
	
	if updatedModel.height != 30 {
		t.Errorf("Expected height 30, got %d", updatedModel.height)
	}
}

func TestModelKeyHandling(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	// Test character input
	keyMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'h'},
	}
	
	newModel, cmd := model.Update(keyMsg)
	if cmd != nil {
		t.Error("Character input should not return command")
	}
	
	updatedModel := newModel.(Model)
	if updatedModel.inputValue != "h" {
		t.Errorf("Expected input 'h', got %s", updatedModel.inputValue)
	}
	
	// Test backspace
	backspaceMsg := tea.KeyMsg{
		Type: tea.KeyBackspace,
	}
	
	newModel, cmd = updatedModel.Update(backspaceMsg)
	if cmd != nil {
		t.Error("Backspace should not return command")
	}
	
	updatedModel = newModel.(Model)
	if updatedModel.inputValue != "" {
		t.Errorf("Expected empty input after backspace, got %s", updatedModel.inputValue)
	}
}

func TestModelScrolling(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	// Test scroll up
	model.scrollOffset = 5
	upMsg := tea.KeyMsg{
		Type: tea.KeyUp,
	}
	
	newModel, cmd := model.Update(upMsg)
	if cmd != nil {
		t.Error("Scroll up should not return command")
	}
	
	updatedModel := newModel.(Model)
	if updatedModel.scrollOffset != 4 {
		t.Errorf("Expected scroll offset 4, got %d", updatedModel.scrollOffset)
	}
	
	// Test scroll down
	downMsg := tea.KeyMsg{
		Type: tea.KeyDown,
	}
	
	newModel, cmd = updatedModel.Update(downMsg)
	if cmd != nil {
		t.Error("Scroll down should not return command")
	}
	
	updatedModel = newModel.(Model)
	if updatedModel.scrollOffset != 5 {
		t.Errorf("Expected scroll offset 5, got %d", updatedModel.scrollOffset)
	}
}

func TestModelMainMenuNavigation(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	// Test new game selection
	model.inputValue = "1"
	enterMsg := tea.KeyMsg{
		Type: tea.KeyEnter,
	}
	
	newModel, cmd := model.Update(enterMsg)
	if cmd != nil {
		t.Error("Enter should not return command")
	}
	
	updatedModel := newModel.(Model)
	if updatedModel.mode != ModeWorldSetup {
		t.Errorf("Expected mode ModeWorldSetup, got %v", updatedModel.mode)
	}
	
	if updatedModel.inputValue != "" {
		t.Error("Input should be cleared after selection")
	}
}

func TestModelQuitHandling(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	// Test quit key
	quitMsg := tea.KeyMsg{
		Type: tea.KeyCtrlC,
	}
	
	_, cmd := model.Update(quitMsg)
	if cmd == nil {
		t.Error("Quit should return tea.Quit command")
	}
	
	// Test q key
	qMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune{'q'},
	}
	
	_, cmd = model.Update(qMsg)
	if cmd == nil {
		t.Error("Q key should return tea.Quit command")
	}
}

func TestModelViewRendering(t *testing.T) {
	cfg := &config.Config{
		Terminal: config.TerminalConfig{
			Width:  80,
			Height: 24,
		},
	}
	model := NewModel(cfg)
	
	// Test main menu view
	view := model.View()
	if view == "" {
		t.Error("View should not be empty")
	}
	
	if !strings.Contains(view, "AXON") {
		t.Error("Main menu should contain game title")
	}
	
	if !strings.Contains(view, "New Game") {
		t.Error("Main menu should contain New Game option")
	}
	
	// Test world setup view
	model.mode = ModeWorldSetup
	view = model.View()
	if !strings.Contains(view, "WORLD SETUP") {
		t.Error("World setup view should contain title")
	}
	
	// Test playing view
	model.mode = ModePlaying
	view = model.View()
	if view == "" {
		t.Error("Playing view should not be empty")
	}
	
	// Test settings view
	model.mode = ModeSettings
	view = model.View()
	if !strings.Contains(view, "Settings") {
		t.Error("Settings view should contain Settings")
	}
	
	// Test save/load view
	model.mode = ModeSaveLoad
	view = model.View()
	if !strings.Contains(view, "saves") {
		t.Error("Save/load view should mention saves")
	}
}

func TestModelGameModes(t *testing.T) {
	// Test GameMode constants
	if ModeMainMenu != 0 {
		t.Errorf("Expected ModeMainMenu to be 0, got %d", ModeMainMenu)
	}
	
	if ModeWorldSetup != 1 {
		t.Errorf("Expected ModeWorldSetup to be 1, got %d", ModeWorldSetup)
	}
	
	if ModePlaying != 2 {
		t.Errorf("Expected ModePlaying to be 2, got %d", ModePlaying)
	}
	
	if ModeSettings != 3 {
		t.Errorf("Expected ModeSettings to be 3, got %d", ModeSettings)
	}
	
	if ModeSaveLoad != 4 {
		t.Errorf("Expected ModeSaveLoad to be 4, got %d", ModeSaveLoad)
	}
}

func TestModelErrorHandling(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	
	// Test error message display
	model.errorMessage = "Test error"
	view := model.View()
	
	if !strings.Contains(view, "Test error") {
		t.Error("View should contain error message")
	}
	
	// Note: The error message clearing behavior depends on the implementation
	// Since View() is called on a copy in this test, the original model's
	// errorMessage field won't be modified. In the actual UI, the model
	// would be properly updated through the Update method.
}

func TestModelLoadingState(t *testing.T) {
	cfg := &config.Config{}
	model := NewModel(cfg)
	model.mode = ModeWorldSetup
	
	// Test loading state display
	model.isLoading = true
	view := model.View()
	
	if !strings.Contains(view, "Creating world") {
		t.Error("View should show loading message when creating world")
	}
	
	// Test loading state in game mode
	model.mode = ModePlaying
	view = model.View()
	
	if !strings.Contains(view, "Processing") {
		t.Error("View should show processing message when loading in game")
	}
}

