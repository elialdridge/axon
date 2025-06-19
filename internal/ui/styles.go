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

// NewStyles creates new monochrome styles following project design
func NewStyles() *Styles {
	// Define monochrome color scheme
	white := lipgloss.Color("#FFFFFF")
	black := lipgloss.Color("#000000")

	return &Styles{
		Base: lipgloss.NewStyle().
			Foreground(white).
			Background(black),

		HistoryPanel: lipgloss.NewStyle().
			Foreground(white).
			Background(black).
			Padding(0, 1),

		InputPanel: lipgloss.NewStyle().
			Foreground(white).
			Background(black).
			Padding(0, 1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()),

		PlayerText: lipgloss.NewStyle().
			Foreground(white).
			Bold(true),

		NarratorText: lipgloss.NewStyle().
			Foreground(white),

		SystemText: lipgloss.NewStyle().
			Foreground(white).
			Italic(true),

		InventoryText: lipgloss.NewStyle().
			Foreground(white).
			Faint(true),

		Border: lipgloss.NewStyle().
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()),

		Prompt: lipgloss.NewStyle().
			Foreground(white).
			Bold(true),

		Scrollbar: lipgloss.NewStyle().
			Foreground(white),
	}
}

// NewMinimalStyles creates styles suitable for minimal terminals
func NewMinimalStyles() *Styles {
	return &Styles{
		Base:          lipgloss.NewStyle(),
		HistoryPanel:  lipgloss.NewStyle(),
		InputPanel:    lipgloss.NewStyle(),
		PlayerText:    lipgloss.NewStyle(),
		NarratorText:  lipgloss.NewStyle(),
		SystemText:    lipgloss.NewStyle(),
		InventoryText: lipgloss.NewStyle(),
		Border:        lipgloss.NewStyle(),
		Prompt:        lipgloss.NewStyle(),
		Scrollbar:     lipgloss.NewStyle(),
	}
}

// NewSystemVStyles creates styles suitable for UNIX System V terminals
func NewSystemVStyles() *Styles {
	return &Styles{
		Base:          lipgloss.NewStyle(),
		HistoryPanel:  lipgloss.NewStyle(),
		InputPanel:    lipgloss.NewStyle(),
		PlayerText:    lipgloss.NewStyle(),
		NarratorText:  lipgloss.NewStyle(),
		SystemText:    lipgloss.NewStyle(),
		InventoryText: lipgloss.NewStyle(),
		Border:        lipgloss.NewStyle(),
		Prompt:        lipgloss.NewStyle(),
		Scrollbar:     lipgloss.NewStyle(),
	}
}
