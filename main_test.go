package main

import (
	"os"
	"testing"
)

func TestDistributeChores_BasicDistribution(t *testing.T) {
	chores := []Chore{
		{Name: "Kitchen", Difficulty: 6, Earned: 5},
		{Name: "Bathroom", Difficulty: 5, Earned: 4},
		{Name: "Living room", Difficulty: 4, Earned: 3},
		{Name: "Bedroom", Difficulty: 3, Earned: 2},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// Check all chores were assigned
	totalChoresAssigned := 0
	for _, person := range result {
		totalChoresAssigned += len(person.Chores)
	}
	if totalChoresAssigned != len(chores) {
		t.Errorf("Expected %d chores assigned, got %d", len(chores), totalChoresAssigned)
	}

	// Check earnings are relatively balanced (within 2 dollars)
	if len(result) >= 2 {
		diff := abs(result[0].TotalEarned - result[1].TotalEarned)
		if diff > 2 {
			t.Errorf("Earnings not balanced: Alice=$%d, Bob=$%d (diff=$%d)",
				result[0].TotalEarned, result[1].TotalEarned, diff)
		}
	}
}

func TestDistributeChores_WithCapacity(t *testing.T) {
	chores := []Chore{
		{Name: "Chore1", Difficulty: 10, Earned: 5},
		{Name: "Chore2", Difficulty: 5, Earned: 3},
		{Name: "Chore3", Difficulty: 5, Earned: 3},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 10, Chores: []Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// Alice should only get Chore1 (difficulty 10, at capacity)
	if result[0].TotalDifficulty > result[0].EffortCapacity {
		t.Errorf("Alice exceeded capacity: %d > %d",
			result[0].TotalDifficulty, result[0].EffortCapacity)
	}

	// Bob should get the remaining chores
	if result[1].TotalDifficulty != 10 {
		t.Errorf("Bob should have difficulty 10, got %d", result[1].TotalDifficulty)
	}
}

func TestDistributeChores_NoCapacity(t *testing.T) {
	chores := []Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 5, Earned: 4},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// With no capacity limits and same earned value, should distribute evenly
	if len(result[0].Chores) != 1 || len(result[1].Chores) != 1 {
		t.Errorf("Expected 1 chore each, got Alice=%d, Bob=%d",
			len(result[0].Chores), len(result[1].Chores))
	}
}

func TestDistributeChores_InsufficientCapacity(t *testing.T) {
	chores := []Chore{
		{Name: "BigChore", Difficulty: 20, Earned: 10},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 10, Chores: []Chore{}},
		{Name: "Bob", EffortCapacity: 10, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// Neither person can take the chore, so no one should have any chores
	totalAssigned := 0
	for _, person := range result {
		totalAssigned += len(person.Chores)
	}
	if totalAssigned != 0 {
		t.Errorf("Expected 0 chores assigned when all exceed capacity, got %d", totalAssigned)
	}
}

func TestDistributeChores_SinglePerson(t *testing.T) {
	chores := []Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 3, Earned: 2},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// Single person should get all chores
	if len(result[0].Chores) != len(chores) {
		t.Errorf("Expected %d chores, got %d", len(chores), len(result[0].Chores))
	}

	expectedEarned := 6
	if result[0].TotalEarned != expectedEarned {
		t.Errorf("Expected total earned $%d, got $%d", expectedEarned, result[0].TotalEarned)
	}

	expectedDifficulty := 8
	if result[0].TotalDifficulty != expectedDifficulty {
		t.Errorf("Expected total difficulty %d, got %d", expectedDifficulty, result[0].TotalDifficulty)
	}
}

func TestDistributeChores_EmptyChores(t *testing.T) {
	chores := []Chore{}

	people := []Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// No chores to assign
	for _, person := range result {
		if len(person.Chores) != 0 {
			t.Errorf("Expected no chores for %s, got %d", person.Name, len(person.Chores))
		}
	}
}

func TestDistributeChores_TotalsCalculatedCorrectly(t *testing.T) {
	chores := []Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 3, Earned: 2},
		{Name: "Chore3", Difficulty: 4, Earned: 3},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// Verify totals are calculated correctly
	expectedDifficulty := 12
	expectedEarned := 9

	if result[0].TotalDifficulty != expectedDifficulty {
		t.Errorf("Expected total difficulty %d, got %d", expectedDifficulty, result[0].TotalDifficulty)
	}

	if result[0].TotalEarned != expectedEarned {
		t.Errorf("Expected total earned %d, got %d", expectedEarned, result[0].TotalEarned)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
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
	config, err := loadConfig(tmpfile.Name())
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

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := loadConfig("nonexistent_file.json")
	if err == nil {
		t.Error("Expected error when loading nonexistent file, got nil")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
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
	_, err = loadConfig(tmpfile.Name())
	if err == nil {
		t.Error("Expected error when loading invalid JSON, got nil")
	}
}

func TestDistributeChores_CapacityEdgeCase(t *testing.T) {
	chores := []Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 5, Earned: 4},
		{Name: "Chore3", Difficulty: 1, Earned: 1},
	}

	people := []Person{
		{Name: "Alice", EffortCapacity: 10, Chores: []Chore{}},
		{Name: "Bob", EffortCapacity: 10, Chores: []Chore{}},
	}

	result := distributeChores(chores, people)

	// Each person should be at or under capacity
	for _, person := range result {
		if person.EffortCapacity > 0 && person.TotalDifficulty > person.EffortCapacity {
			t.Errorf("%s exceeded capacity: %d > %d",
				person.Name, person.TotalDifficulty, person.EffortCapacity)
		}
	}

	// All chores should be assigned
	totalAssigned := 0
	for _, person := range result {
		totalAssigned += len(person.Chores)
	}
	if totalAssigned != len(chores) {
		t.Errorf("Expected %d chores assigned, got %d", len(chores), totalAssigned)
	}
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
