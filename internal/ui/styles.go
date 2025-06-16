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

// NewStyles creates new monochrome styles
func NewStyles() *Styles {
	return &Styles{
		Base: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")),

		HistoryPanel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Padding(0, 1),

		InputPanel: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1),

		PlayerText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Bold(true),

		NarratorText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")),

		SystemText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Italic(true),

		InventoryText: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Faint(true),

		Border: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#FFFFFF")),

		Prompt: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Bold(true),

		Scrollbar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")),
	}
}
