package config

import (
	"os"
	"testing"
)

func TestLoad_ValidFile(t *testing.T) {
	// Create a temporary config file
	configContent := `{
  "chores": [
    {"Name": "Kitchen", "Difficulty": 6, "Earned": 5},
    {"Name": "Bathroom", "Difficulty": 5, "Earned": 4}
  ],
  "people": [
    {"Name": "Alice", "EffortCapacity": 0},
    {"Name": "Bob", "EffortCapacity": 15}
  ]
}`

	tmpfile, err := os.CreateTemp("", "test_config_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the config
	config, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify chores
	if len(config.Chores) != 2 {
		t.Errorf("Expected 2 chores, got %d", len(config.Chores))
	}

	if config.Chores[0].Name != "Kitchen" {
		t.Errorf("Expected first chore to be 'Kitchen', got '%s'", config.Chores[0].Name)
	}

	// Verify people
	if len(config.People) != 2 {
		t.Errorf("Expected 2 people, got %d", len(config.People))
	}

	if config.People[0].Name != "Alice" {
		t.Errorf("Expected first person to be 'Alice', got '%s'", config.People[0].Name)
	}

	if config.People[1].EffortCapacity != 15 {
		t.Errorf("Expected Bob's capacity to be 15, got %d", config.People[1].EffortCapacity)
	}

	// Verify chores slices are initialized
	for i, person := range config.People {
		if person.Chores == nil {
			t.Errorf("Person %d (%s) has nil Chores slice", i, person.Name)
		}
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("nonexistent_file.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tmpfile, err := os.CreateTemp("", "test_invalid_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	invalidJSON := `{"chores": [invalid json}`
	if _, err := tmpfile.Write([]byte(invalidJSON)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Try to load the invalid config
	_, err = Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error when loading invalid JSON, got nil")
	}
}
