package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all UI styles for the game
type Styles struct {
	Base          lipgloss.Style
	HistoryPanel  lipgloss.Style
	InputPanel    lipgloss.Style
	PlayerText    lipgloss.Style
	NarratorText  lipgloss.Style
	SystemText    lipgloss.Style
	InventoryText lipgloss.Style
	Border        lipgloss.Style
	Prompt        lipgloss.Style
	Scrollbar     lipgloss.Style
}

// NewStyles creates new simple ASCII-only styles
func NewStyles() *Styles {
	return &Styles{
		Base: lipgloss.NewStyle(),

		HistoryPanel: lipgloss.NewStyle().
			Padding(0, 1),

		InputPanel: lipgloss.NewStyle().
			Padding(0, 1),

		PlayerText: lipgloss.NewStyle(),

		NarratorText: lipgloss.NewStyle(),

		SystemText: lipgloss.NewStyle(),

		InventoryText: lipgloss.NewStyle(),

		Border: lipgloss.NewStyle(),

		Prompt: lipgloss.NewStyle(),

		Scrollbar: lipgloss.NewStyle(),
	}
}
