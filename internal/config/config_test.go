package config

import (
	"os"
	"testing"
)

func TestLoad_ValidFile(t *testing.T) {
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

	config, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Chores) != 2 {
		t.Errorf("Expected 2 chores, got %d", len(config.Chores))
	}

	if config.Chores[0].Name != "Kitchen" {
		t.Errorf("Expected first chore to be 'Kitchen', got '%s'", config.Chores[0].Name)
	}

	if len(config.People) != 2 {
		t.Errorf("Expected 2 people, got %d", len(config.People))
	}

	if config.People[0].Name != "Alice" {
		t.Errorf("Expected first person to be 'Alice', got '%s'", config.People[0].Name)
	}

	if config.People[1].EffortCapacity != 15 {
		t.Errorf("Expected Bob's capacity to be 15, got %d", config.People[1].EffortCapacity)
	}

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

	_, err = Load(tmpfile.Name())
	if err == nil {
		t.Error("Expected error when loading invalid JSON, got nil")
	}
}

func TestLoad_WithPreAssignedChores(t *testing.T) {
	configContent := `{
  "chores": [
    {"Name": "Kitchen", "Difficulty": 6, "Earned": 5}
  ],
  "people": [
    {
      "Name": "Tommy",
      "EffortCapacity": 10,
      "PreAssignedChores": [
        {
          "Name": "Clean Bedroom",
          "Difficulty": 2,
          "Earned": 2,
          "Description": "Clean and organize personal bedroom"
        },
        {
          "Name": "Feed Pets",
          "Difficulty": 1,
          "Earned": 0
        }
      ]
    },
    {"Name": "Alice", "EffortCapacity": 0}
  ]
}`

	tmpfile, err := os.CreateTemp("", "test_preassigned_*.json")
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

	config, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.People[0].PreAssignedChores) != 2 {
		t.Errorf("Expected 2 pre-assigned chores for Tommy, got %d", len(config.People[0].PreAssignedChores))
	}

	if config.People[0].PreAssignedChores[0].Name != "Clean Bedroom" {
		t.Errorf("Expected first pre-assigned chore to be 'Clean Bedroom', got '%s'", config.People[0].PreAssignedChores[0].Name)
	}

	if config.People[0].PreAssignedChores[0].Description != "Clean and organize personal bedroom" {
		t.Errorf("Expected description for first pre-assigned chore, got '%s'", config.People[0].PreAssignedChores[0].Description)
	}

	if config.People[0].TotalDifficulty != 3 {
		t.Errorf("Expected TotalDifficulty to be 3 (2+1), got %d", config.People[0].TotalDifficulty)
	}

	if config.People[0].TotalEarned != 2 {
		t.Errorf("Expected TotalEarned to be 2 (2+0), got %d", config.People[0].TotalEarned)
	}

	if config.People[1].TotalDifficulty != 0 {
		t.Errorf("Expected Alice's TotalDifficulty to be 0, got %d", config.People[1].TotalDifficulty)
	}

	if config.People[1].TotalEarned != 0 {
		t.Errorf("Expected Alice's TotalEarned to be 0, got %d", config.People[1].TotalEarned)
	}
}

func TestLoad_ZeroValuePreAssignedChores(t *testing.T) {
	configContent := `{
  "chores": [],
  "people": [
    {
      "Name": "Alice",
      "EffortCapacity": 5,
      "PreAssignedChores": [
        {
          "Name": "Test Chore",
          "Difficulty": 0,
          "Earned": 0
        }
      ]
    }
  ]
}`

	tmpfile, err := os.CreateTemp("", "test_zero_*.json")
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

	config, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.People[0].PreAssignedChores) != 1 {
		t.Errorf("Expected 1 pre-assigned chore, got %d", len(config.People[0].PreAssignedChores))
	}

	if config.People[0].TotalDifficulty != 0 {
		t.Errorf("Expected TotalDifficulty to be 0, got %d", config.People[0].TotalDifficulty)
	}

	if config.People[0].TotalEarned != 0 {
		t.Errorf("Expected TotalEarned to be 0, got %d", config.People[0].TotalEarned)
	}
}
