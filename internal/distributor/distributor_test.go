package distributor

import (
	"bytes"
	"strings"
	"testing"

	"github.com/faradayfan/chore-distributor/internal/models"
)

func TestDistribute_BasicDistribution(t *testing.T) {
	chores := []models.Chore{
		{Name: "Kitchen", Difficulty: 6, Earned: 5},
		{Name: "Bathroom", Difficulty: 5, Earned: 4},
		{Name: "Living room", Difficulty: 4, Earned: 3},
		{Name: "Bedroom", Difficulty: 3, Earned: 2},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []models.Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	totalChoresAssigned := 0
	for _, person := range result {
		totalChoresAssigned += len(person.Chores)
	}
	if totalChoresAssigned != len(chores) {
		t.Errorf("Expected %d chores assigned, got %d", len(chores), totalChoresAssigned)
	}

	if len(result) >= 2 {
		diff := abs(result[0].TotalEarned - result[1].TotalEarned)
		if diff > 2 {
			t.Errorf("Earnings not balanced: Alice=$%d, Bob=$%d (diff=$%d)",
				result[0].TotalEarned, result[1].TotalEarned, diff)
		}
	}
}

func TestDistribute_WithCapacity(t *testing.T) {
	chores := []models.Chore{
		{Name: "Chore1", Difficulty: 10, Earned: 5},
		{Name: "Chore2", Difficulty: 5, Earned: 3},
		{Name: "Chore3", Difficulty: 5, Earned: 3},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 10, Chores: []models.Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	if result[0].TotalDifficulty > result[0].EffortCapacity {
		t.Errorf("Alice exceeded capacity: %d > %d",
			result[0].TotalDifficulty, result[0].EffortCapacity)
	}

	if result[1].TotalDifficulty != 10 {
		t.Errorf("Bob should have difficulty 10, got %d", result[1].TotalDifficulty)
	}
}

func TestDistribute_NoCapacity(t *testing.T) {
	chores := []models.Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 5, Earned: 4},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []models.Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	if len(result[0].Chores) != 1 || len(result[1].Chores) != 1 {
		t.Errorf("Expected 1 chore each, got Alice=%d, Bob=%d",
			len(result[0].Chores), len(result[1].Chores))
	}
}

func TestDistribute_InsufficientCapacity(t *testing.T) {
	chores := []models.Chore{
		{Name: "BigChore", Difficulty: 20, Earned: 10},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 10, Chores: []models.Chore{}},
		{Name: "Bob", EffortCapacity: 10, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	totalAssigned := 0
	for _, person := range result {
		totalAssigned += len(person.Chores)
	}
	if totalAssigned != 0 {
		t.Errorf("Expected 0 chores assigned when all exceed capacity, got %d", totalAssigned)
	}
}

func TestDistribute_SinglePerson(t *testing.T) {
	chores := []models.Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 3, Earned: 2},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

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

func TestDistribute_EmptyChores(t *testing.T) {
	chores := []models.Chore{}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []models.Chore{}},
		{Name: "Bob", EffortCapacity: 0, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	for _, person := range result {
		if len(person.Chores) != 0 {
			t.Errorf("Expected no chores for %s, got %d", person.Name, len(person.Chores))
		}
	}
}

func TestDistribute_TotalsCalculatedCorrectly(t *testing.T) {
	chores := []models.Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 3, Earned: 2},
		{Name: "Chore3", Difficulty: 4, Earned: 3},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 0, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	expectedDifficulty := 12
	expectedEarned := 9

	if result[0].TotalDifficulty != expectedDifficulty {
		t.Errorf("Expected total difficulty %d, got %d", expectedDifficulty, result[0].TotalDifficulty)
	}

	if result[0].TotalEarned != expectedEarned {
		t.Errorf("Expected total earned %d, got %d", expectedEarned, result[0].TotalEarned)
	}
}

func TestDistribute_CapacityEdgeCase(t *testing.T) {
	chores := []models.Chore{
		{Name: "Chore1", Difficulty: 5, Earned: 4},
		{Name: "Chore2", Difficulty: 5, Earned: 4},
		{Name: "Chore3", Difficulty: 1, Earned: 1},
	}

	people := []models.Person{
		{Name: "Alice", EffortCapacity: 10, Chores: []models.Chore{}},
		{Name: "Bob", EffortCapacity: 10, Chores: []models.Chore{}},
	}

	result := Distribute(chores, people)

	for _, person := range result {
		if person.EffortCapacity > 0 && person.TotalDifficulty > person.EffortCapacity {
			t.Errorf("%s exceeded capacity: %d > %d",
				person.Name, person.TotalDifficulty, person.EffortCapacity)
		}
	}

	totalAssigned := 0
	for _, person := range result {
		totalAssigned += len(person.Chores)
	}
	if totalAssigned != len(chores) {
		t.Errorf("Expected %d chores assigned, got %d", len(chores), totalAssigned)
	}
}

func TestPrintDistribution_DefaultMode(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Alice",
			EffortCapacity:  10,
			Chores:          []models.Chore{{Name: "Kitchen", Difficulty: 6, Earned: 5}},
			TotalDifficulty: 6,
			TotalEarned:     5,
		},
	}

	var buf bytes.Buffer
	PrintDistribution(&buf, people, PrintOptions{Verbose: false})
	output := buf.String()

	if !strings.Contains(output, "Earns: $5") {
		t.Error("Default output should contain earnings")
	}
	if !strings.Contains(output, "Total Earned: $5") {
		t.Error("Default output should contain total earned")
	}

	if strings.Contains(output, "Difficulty:") {
		t.Error("Default output should not contain difficulty")
	}
	if strings.Contains(output, "Effort Capacity") {
		t.Error("Default output should not contain effort capacity")
	}
	if strings.Contains(output, "Total Difficulty") {
		t.Error("Default output should not contain total difficulty")
	}
}

func TestPrintDistribution_VerboseMode(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Alice",
			EffortCapacity:  10,
			Chores:          []models.Chore{{Name: "Kitchen", Difficulty: 6, Earned: 5}},
			TotalDifficulty: 6,
			TotalEarned:     5,
		},
	}

	var buf bytes.Buffer
	PrintDistribution(&buf, people, PrintOptions{Verbose: true})
	output := buf.String()

	if !strings.Contains(output, "Earns: $5") {
		t.Error("Verbose output should contain earnings")
	}
	if !strings.Contains(output, "Difficulty: 6") {
		t.Error("Verbose output should contain difficulty")
	}
	if !strings.Contains(output, "Effort Capacity: 10") {
		t.Error("Verbose output should contain effort capacity")
	}
	if !strings.Contains(output, "Total Difficulty: 6 / 10") {
		t.Error("Verbose output should contain total difficulty with capacity")
	}
}

func TestPrintDistribution_VerboseNoCapacity(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Bob",
			EffortCapacity:  0, 
			Chores:          []models.Chore{{Name: "Kitchen", Difficulty: 6, Earned: 5}},
			TotalDifficulty: 6,
			TotalEarned:     5,
		},
	}

	var buf bytes.Buffer
	PrintDistribution(&buf, people, PrintOptions{Verbose: true})
	output := buf.String()

	if !strings.Contains(output, "Total Difficulty: 6") {
		t.Error("Verbose output should contain total difficulty")
	}
	if strings.Contains(output, "Effort Capacity") {
		t.Error("Should not show effort capacity when it's 0")
	}
	if strings.Contains(output, "/ 0") {
		t.Error("Should not show capacity limit when it's 0")
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
