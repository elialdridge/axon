package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all game configuration
type Config struct {
	// Terminal settings
	Terminal TerminalConfig `json:"terminal"`
	// AI settings
	AI AIConfig `json:"ai"`
	// Game settings
	Game GameConfig `json:"game"`
}

// TerminalConfig contains terminal-specific settings
type TerminalConfig struct {
	Width            int  `json:"width"`
	Height           int  `json:"height"`
	ColorEnabled     bool `json:"color_enabled"`
	ForceMinimal     bool `json:"force_minimal"`
	ForceSystemV     bool `json:"force_systemv"`
	AutoDetect       bool `json:"auto_detect"`
	MouseEnabled     bool `json:"mouse_enabled"`
	AltScreenEnabled bool `json:"alt_screen_enabled"`
}

// AIConfig contains AI model settings
type AIConfig struct {
	OpenRouterAPIKey string `json:"openrouter_api_key"`
	GeminiAPIKey     string `json:"gemini_api_key"`
	DefaultModel     string `json:"default_model"`
}

// GameConfig contains game-specific settings
type GameConfig struct {
	HistoryLimit int    `json:"history_limit"`
	SaveDir      string `json:"save_dir"`
}

// Load loads configuration from file or creates default
func Load() *Config {
	cfg := defaultConfig()

	// Try to load from config file
	configPath := getConfigPath()
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, cfg); err == nil {
			return cfg
		}
	}

	return cfg
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0o644)
}

// defaultConfig returns default configuration
func defaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	saveDir := filepath.Join(homeDir, ".axon", "saves")

	return &Config{
		Terminal: TerminalConfig{
			Width:            80,
			Height:           24,
			ColorEnabled:     false, // Monochrome by default
			ForceMinimal:     false,
			ForceSystemV:     false,
			AutoDetect:       true, // Enable auto-detection by default
			MouseEnabled:     true, // Will be disabled for incompatible terminals
			AltScreenEnabled: true, // Will be disabled for incompatible terminals
		},
		AI: AIConfig{
			OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
			GeminiAPIKey:     os.Getenv("GEMINI_API_KEY"),
			DefaultModel:     "openai/gpt-4o-mini",
		},
		Game: GameConfig{
			HistoryLimit: 1000,
			SaveDir:      saveDir,
		},
	}
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".axon", "config.json")
}
