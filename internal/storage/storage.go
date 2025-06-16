package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Storage handles game save/load operations
type Storage struct {
	saveDir string
}

// NewStorage creates a new storage instance
func NewStorage(saveDir string) *Storage {
	return &Storage{
		saveDir: saveDir,
	}
}

// SaveGame saves the game state to disk
func (s *Storage) SaveGame(name string, state interface{}) error {
	// Ensure save directory exists
	if err := os.MkdirAll(s.saveDir, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// Generate filename with timestamp if name is empty
	if name == "" {
		name = fmt.Sprintf("save_%s", time.Now().Format("20060102_150405"))
	}

	filePath := filepath.Join(s.saveDir, name+".json")

	// Marshal game state to JSON
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal game state: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	return nil
}

// LoadGame loads a game state from disk into the provided state interface
func (s *Storage) LoadGame(name string, state interface{}) error {
	filePath := filepath.Join(s.saveDir, name+".json")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("save file not found: %s", name)
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read save file: %w", err)
	}

	// Unmarshal JSON into provided state
	if err := json.Unmarshal(data, state); err != nil {
		return fmt.Errorf("failed to unmarshal game state: %w", err)
	}

	return nil
}

// ListSaves returns a list of available save files
func (s *Storage) ListSaves() ([]string, error) {
	files, err := os.ReadDir(s.saveDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read save directory: %w", err)
	}

	var saves []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			name := file.Name()[:len(file.Name())-5] // Remove .json extension
			saves = append(saves, name)
		}
	}

	return saves, nil
}

// DeleteSave deletes a save file
func (s *Storage) DeleteSave(name string) error {
	filePath := filepath.Join(s.saveDir, name+".json")
	return os.Remove(filePath)
}
