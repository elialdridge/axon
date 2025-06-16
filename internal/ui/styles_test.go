package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewStyles(t *testing.T) {
	styles := NewStyles()

	if styles == nil {
		t.Fatal("NewStyles returned nil")
	}

	// Test that styles can render (indicating they are initialized)
	if styles.Base.Render("test") == "" {
		t.Error("Base style should be able to render")
	}

	if styles.HistoryPanel.Render("test") == "" {
		t.Error("HistoryPanel style should be able to render")
	}

	if styles.InputPanel.Render("test") == "" {
		t.Error("InputPanel style should be able to render")
	}

	if styles.PlayerText.Render("test") == "" {
		t.Error("PlayerText style should be able to render")
	}

	if styles.NarratorText.Render("test") == "" {
		t.Error("NarratorText style should be able to render")
	}

	if styles.SystemText.Render("test") == "" {
		t.Error("SystemText style should be able to render")
	}

	if styles.InventoryText.Render("test") == "" {
		t.Error("InventoryText style should be able to render")
	}

	if styles.Border.Render("test") == "" {
		t.Error("Border style should be able to render")
	}

	if styles.Prompt.Render("test") == "" {
		t.Error("Prompt style should be able to render")
	}

	if styles.Scrollbar.Render("test") == "" {
		t.Error("Scrollbar style should be able to render")
	}
}

func TestStylesMonochrome(t *testing.T) {
	styles := NewStyles()

	// Test monochrome color scheme
	white := lipgloss.Color("#FFFFFF")
	black := lipgloss.Color("#000000")

	// Test base style colors
	baseStyle := styles.Base
	if baseStyle.GetForeground() != white {
		t.Error("Base style should have white foreground")
	}

	if baseStyle.GetBackground() != black {
		t.Error("Base style should have black background")
	}

	// Test history panel colors
	historyStyle := styles.HistoryPanel
	if historyStyle.GetForeground() != white {
		t.Error("HistoryPanel style should have white foreground")
	}

	if historyStyle.GetBackground() != black {
		t.Error("HistoryPanel style should have black background")
	}

	// Test input panel colors
	inputStyle := styles.InputPanel
	if inputStyle.GetForeground() != white {
		t.Error("InputPanel style should have white foreground")
	}

	if inputStyle.GetBackground() != black {
		t.Error("InputPanel style should have black background")
	}
}

func TestStylesFormatting(t *testing.T) {
	styles := NewStyles()

	// Test player text formatting
	playerStyle := styles.PlayerText
	if !playerStyle.GetBold() {
		t.Error("PlayerText should be bold")
	}

	// Test system text formatting
	systemStyle := styles.SystemText
	if !systemStyle.GetItalic() {
		t.Error("SystemText should be italic")
	}

	// Test inventory text formatting
	inventoryStyle := styles.InventoryText
	if !inventoryStyle.GetFaint() {
		t.Error("InventoryText should be faint")
	}

	// Test prompt formatting
	promptStyle := styles.Prompt
	if !promptStyle.GetBold() {
		t.Error("Prompt should be bold")
	}
}

func TestStylesBorders(t *testing.T) {
	styles := NewStyles()

	// Test input panel border
	inputStyle := styles.InputPanel
	if !inputStyle.GetBorderTop() {
		t.Error("InputPanel should have top border")
	}

	if inputStyle.GetBorderBottom() || inputStyle.GetBorderLeft() || inputStyle.GetBorderRight() {
		t.Error("InputPanel should only have top border")
	}

	// Test border style
	borderStyle := styles.Border
	if !borderStyle.GetBorderTop() || !borderStyle.GetBorderBottom() ||
		!borderStyle.GetBorderLeft() || !borderStyle.GetBorderRight() {
		t.Error("Border style should have all borders")
	}
}

func TestStylesPadding(t *testing.T) {
	styles := NewStyles()

	// Test history panel padding
	historyStyle := styles.HistoryPanel
	top, right, bottom, left := historyStyle.GetPadding()
	if top != 0 || right != 1 || bottom != 0 || left != 1 {
		t.Errorf("HistoryPanel padding should be (0,1,0,1), got (%d,%d,%d,%d)", top, right, bottom, left)
	}

	// Test input panel padding
	inputStyle := styles.InputPanel
	top, right, bottom, left = inputStyle.GetPadding()
	if top != 0 || right != 1 || bottom != 0 || left != 1 {
		t.Errorf("InputPanel padding should be (0,1,0,1), got (%d,%d,%d,%d)", top, right, bottom, left)
	}
}

func TestStylesRendering(t *testing.T) {
	styles := NewStyles()

	// Test that styles can render text
	playerText := styles.PlayerText.Render("Player says hello")
	if playerText == "" {
		t.Error("PlayerText should render non-empty string")
	}

	narratorText := styles.NarratorText.Render("The narrator describes the scene")
	if narratorText == "" {
		t.Error("NarratorText should render non-empty string")
	}

	systemText := styles.SystemText.Render("System message")
	if systemText == "" {
		t.Error("SystemText should render non-empty string")
	}

	inventoryText := styles.InventoryText.Render("Inventory item")
	if inventoryText == "" {
		t.Error("InventoryText should render non-empty string")
	}

	promptText := styles.Prompt.Render("> ")
	if promptText == "" {
		t.Error("Prompt should render non-empty string")
	}
}

func TestStylesWidth(t *testing.T) {
	styles := NewStyles()

	// Test setting width on styles
	width := 80
	baseWithWidth := styles.Base.Width(width)
	if baseWithWidth.GetWidth() != width {
		t.Errorf("Expected width %d, got %d", width, baseWithWidth.GetWidth())
	}

	historyWithWidth := styles.HistoryPanel.Width(width)
	if historyWithWidth.GetWidth() != width {
		t.Errorf("Expected width %d, got %d", width, historyWithWidth.GetWidth())
	}
}

func TestStylesHeight(t *testing.T) {
	styles := NewStyles()

	// Test setting height on styles
	height := 24
	baseWithHeight := styles.Base.Height(height)
	if baseWithHeight.GetHeight() != height {
		t.Errorf("Expected height %d, got %d", height, baseWithHeight.GetHeight())
	}

	historyWithHeight := styles.HistoryPanel.Height(height)
	if historyWithHeight.GetHeight() != height {
		t.Errorf("Expected height %d, got %d", height, historyWithHeight.GetHeight())
	}
}
