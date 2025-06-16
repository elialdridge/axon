package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()

	if cfg == nil {
		t.Fatal("defaultConfig returned nil")
	}

	// Test terminal config defaults
	if cfg.Terminal.Width != 80 {
		t.Errorf("Expected terminal width 80, got %d", cfg.Terminal.Width)
	}

	if cfg.Terminal.Height != 24 {
		t.Errorf("Expected terminal height 24, got %d", cfg.Terminal.Height)
	}

	if cfg.Terminal.ColorEnabled {
		t.Error("Expected color to be disabled by default")
	}

	// Test AI config
	if cfg.AI.DefaultModel != "openai/gpt-4o-mini" {
		t.Errorf("Expected default model 'openai/gpt-4o-mini', got %s", cfg.AI.DefaultModel)
	}

	// Test game config
	if cfg.Game.HistoryLimit != 1000 {
		t.Errorf("Expected history limit 1000, got %d", cfg.Game.HistoryLimit)
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "axon_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test config
	cfg := &Config{
		Terminal: TerminalConfig{
			Width:        100,
			Height:       30,
			ColorEnabled: true,
		},
		AI: AIConfig{
			OpenRouterAPIKey: "test_key",
			GeminiAPIKey:     "test_gemini",
			DefaultModel:     "test_model",
		},
		Game: GameConfig{
			HistoryLimit: 500,
			SaveDir:      tempDir,
		},
	}

	// Test JSON marshaling
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	// Test JSON unmarshaling
	var loadedConfig Config
	err = json.Unmarshal(data, &loadedConfig)
	if err != nil {
		t.Fatal(err)
	}

	// Verify loaded config
	if loadedConfig.Terminal.Width != 100 {
		t.Errorf("Expected width 100, got %d", loadedConfig.Terminal.Width)
	}

	if loadedConfig.AI.OpenRouterAPIKey != "test_key" {
		t.Errorf("Expected OpenRouter key 'test_key', got %s", loadedConfig.AI.OpenRouterAPIKey)
	}
}

func TestLoad(t *testing.T) {
	// Test loading default config when no file exists
	cfg := Load()
	if cfg == nil {
		t.Fatal("Load returned nil")
	}

	// Should return default values
	if cfg.Terminal.Width != 80 {
		t.Errorf("Expected default width 80, got %d", cfg.Terminal.Width)
	}
}

func TestConfigPaths(t *testing.T) {
	configPath := getConfigPath()
	if configPath == "" {
		t.Error("getConfigPath returned empty string")
	}

	// Should contain .axon directory
	if !filepath.IsAbs(configPath) {
		t.Error("Config path should be absolute")
	}

	expectedSuffix := filepath.Join(".axon", "config.json")
	if !strings.HasSuffix(configPath, expectedSuffix) {
		t.Errorf("Config path should end with %s, got %s", expectedSuffix, configPath)
	}
}
