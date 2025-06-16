package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestGameState represents a test game state for storage tests
type TestGameState struct {
	Name      string    `json:"name"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

func TestNewStorage(t *testing.T) {
	saveDir := "/tmp/test_saves"
	storage := NewStorage(saveDir)
	
	if storage == nil {
		t.Fatal("NewStorage returned nil")
	}
	
	if storage.saveDir != saveDir {
		t.Errorf("Expected save dir %s, got %s", saveDir, storage.saveDir)
	}
}

func TestSaveGame(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Create test game state
	state := &TestGameState{
		Name:      "Test Game",
		Level:     5,
		CreatedAt: time.Now(),
	}
	
	// Save with specific name
	err = storage.SaveGame("test_save", state)
	if err != nil {
		t.Errorf("SaveGame failed: %v", err)
	}
	
	// Check file exists
	filePath := filepath.Join(tempDir, "test_save.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Save file was not created")
	}
	
	// Verify file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	
	var loadedState TestGameState
	err = json.Unmarshal(data, &loadedState)
	if err != nil {
		t.Fatal(err)
	}
	
	if loadedState.Name != state.Name {
		t.Errorf("Expected name %s, got %s", state.Name, loadedState.Name)
	}
	
	if loadedState.Level != state.Level {
		t.Errorf("Expected level %d, got %d", state.Level, loadedState.Level)
	}
}

func TestSaveGameAutoName(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Create test game state
	state := &TestGameState{
		Name:      "Auto Save Game",
		Level:     10,
		CreatedAt: time.Now(),
	}
	
	// Save with empty name (should auto-generate)
	err = storage.SaveGame("", state)
	if err != nil {
		t.Errorf("SaveGame failed: %v", err)
	}
	
	// Check that files were created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}
	
	if len(files) == 0 {
		t.Error("No save files were created")
	}
	
	// Check filename format
	fileName := files[0].Name()
	if !strings.HasSuffix(fileName, ".json") {
		t.Error("Save file should have .json extension")
	}
	
	if !strings.HasPrefix(fileName, "save_") {
		t.Error("Auto-generated filename should start with 'save_'")
	}
}

func TestLoadGame(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Create and save test game state
	originalState := &TestGameState{
		Name:      "Load Test Game",
		Level:     7,
		CreatedAt: time.Now(),
	}
	
	err = storage.SaveGame("load_test", originalState)
	if err != nil {
		t.Fatal(err)
	}
	
	// Load the game state
	var loadedState TestGameState
	err = storage.LoadGame("load_test", &loadedState)
	if err != nil {
		t.Errorf("LoadGame failed: %v", err)
	}
	
	// Verify loaded data
	if loadedState.Name != originalState.Name {
		t.Errorf("Expected name %s, got %s", originalState.Name, loadedState.Name)
	}
	
	if loadedState.Level != originalState.Level {
		t.Errorf("Expected level %d, got %d", originalState.Level, loadedState.Level)
	}
}

func TestLoadGameNotFound(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Try to load non-existent game
	var state TestGameState
	err = storage.LoadGame("nonexistent", &state)
	if err == nil {
		t.Error("LoadGame should fail for non-existent save")
	}
	
	if !os.IsNotExist(err) && err.Error() != "save file not found: nonexistent" {
		t.Errorf("Expected 'not found' error, got: %v", err)
	}
}

func TestListSaves(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Test empty directory
	saves, err := storage.ListSaves()
	if err != nil {
		t.Errorf("ListSaves failed: %v", err)
	}
	
	if len(saves) != 0 {
		t.Errorf("Expected 0 saves, got %d", len(saves))
	}
	
	// Create some save files
	state := &TestGameState{Name: "Test", Level: 1, CreatedAt: time.Now()}
	storage.SaveGame("save1", state)
	storage.SaveGame("save2", state)
	storage.SaveGame("save3", state)
	
	// Create a non-save file (should be ignored)
	os.WriteFile(filepath.Join(tempDir, "not_a_save.txt"), []byte("test"), 0644)
	
	// List saves
	saves, err = storage.ListSaves()
	if err != nil {
		t.Errorf("ListSaves failed: %v", err)
	}
	
	if len(saves) != 3 {
		t.Errorf("Expected 3 saves, got %d", len(saves))
	}
	
	// Check save names (order might vary)
	expectedSaves := map[string]bool{"save1": true, "save2": true, "save3": true}
	for _, save := range saves {
		if !expectedSaves[save] {
			t.Errorf("Unexpected save name: %s", save)
		}
		delete(expectedSaves, save)
	}
	
	if len(expectedSaves) > 0 {
		t.Error("Some expected saves were not found")
	}
}

func TestListSavesNonExistentDir(t *testing.T) {
	storage := NewStorage("/path/that/does/not/exist")
	
	saves, err := storage.ListSaves()
	if err != nil {
		t.Errorf("ListSaves should not fail for non-existent directory, got: %v", err)
	}
	
	if len(saves) != 0 {
		t.Errorf("Expected 0 saves for non-existent directory, got %d", len(saves))
	}
}

func TestDeleteSave(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Create a save file
	state := &TestGameState{Name: "To Delete", Level: 1, CreatedAt: time.Now()}
	err = storage.SaveGame("delete_test", state)
	if err != nil {
		t.Fatal(err)
	}
	
	// Verify file exists
	filePath := filepath.Join(tempDir, "delete_test.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("Save file was not created")
	}
	
	// Delete the save
	err = storage.DeleteSave("delete_test")
	if err != nil {
		t.Errorf("DeleteSave failed: %v", err)
	}
	
	// Verify file no longer exists
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Save file was not deleted")
	}
}

func TestDeleteSaveNotFound(t *testing.T) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "axon_storage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	storage := NewStorage(tempDir)
	
	// Try to delete non-existent save
	err = storage.DeleteSave("nonexistent")
	if err == nil {
		t.Error("DeleteSave should fail for non-existent save")
	}
}

