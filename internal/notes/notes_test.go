package notes

import (
	"strings"
	"testing"
	"time"

	"github.com/faradayfan/chore-distributor/internal/models"
)

func TestNewWriter(t *testing.T) {
	writer := NewWriter("Test Note", true)

	if writer.NoteName != "Test Note" {
		t.Errorf("Expected note name 'Test Note', got '%s'", writer.NoteName)
	}
	if !writer.DryRun {
		t.Error("Writer should have DryRun set to true")
	}

	writer = NewWriter("Chores", false)
	if writer.NoteName != "Chores" {
		t.Errorf("Expected note name 'Chores', got '%s'", writer.NoteName)
	}
	if writer.DryRun {
		t.Error("Writer should have DryRun set to false")
	}
}

func TestFormatNoteContentHTML_Default(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Alice",
			EffortCapacity:  10,
			TotalDifficulty: 6,
			TotalEarned:     5,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
			},
		},
	}

	content := formatNoteContentHTML(people, false)

	today := time.Now().Format("January 2, 2006")
	if !strings.Contains(content, today) {
		t.Errorf("Content should contain today's date (%s)", today)
	}

	if !strings.Contains(content, "<div>") {
		t.Error("Content should use div tags")
	}

	if !strings.Contains(content, "<b>Alice</b>") {
		t.Error("Content should contain person name in bold")
	}

	if !strings.Contains(content, "• Kitchen") {
		t.Error("Content should contain chore with bullet")
	}
	if !strings.Contains(content, "$5") {
		t.Error("Content should contain earnings")
	}

	if strings.Contains(content, "Difficulty:") {
		t.Error("Default content should not contain difficulty")
	}
	if strings.Contains(content, "Capacity:") {
		t.Error("Default content should not contain capacity")
	}

	if !strings.Contains(content, "─────") {
		t.Error("Content should contain separator line")
	}
}

func TestFormatNoteContentHTML_Verbose(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Bob",
			EffortCapacity:  15,
			TotalDifficulty: 10,
			TotalEarned:     8,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
				{Name: "Bathroom", Difficulty: 4, Earned: 3},
			},
		},
	}

	content := formatNoteContentHTML(people, true)

	if !strings.Contains(content, "(Capacity: 15)") {
		t.Error("Verbose content should contain capacity")
	}

	if !strings.Contains(content, "Difficulty: 6") {
		t.Error("Verbose content should contain difficulty")
	}

	if !strings.Contains(content, "Effort: 10 / 15") {
		t.Error("Verbose content should contain effort ratio")
	}
}

func TestFormatNoteContentHTML_VerboseNoCapacity(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Charlie",
			EffortCapacity:  0, 
			TotalDifficulty: 6,
			TotalEarned:     5,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
			},
		},
	}

	content := formatNoteContentHTML(people, true)

	if !strings.Contains(content, "Difficulty: 6") {
		t.Error("Verbose content should contain difficulty")
	}

	if strings.Contains(content, "Capacity:") {
		t.Error("Should not show capacity when it's 0")
	}

	if strings.Contains(content, "Effort:") {
		t.Error("Should not show effort when capacity is 0")
	}
}

func TestFormatNoteContentHTML_MultiplePeople(t *testing.T) {
	people := []models.Person{
		{
			Name:        "Alice",
			TotalEarned: 5,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
			},
		},
		{
			Name:        "Bob",
			TotalEarned: 4,
			Chores: []models.Chore{
				{Name: "Bathroom", Difficulty: 5, Earned: 4},
			},
		},
	}

	content := formatNoteContentHTML(people, false)

	if !strings.Contains(content, "<b>Alice</b>") {
		t.Error("Content should contain Alice")
	}
	if !strings.Contains(content, "<b>Bob</b>") {
		t.Error("Content should contain Bob")
	}

	if !strings.Contains(content, "Kitchen") {
		t.Error("Content should contain Kitchen")
	}
	if !strings.Contains(content, "Bathroom") {
		t.Error("Content should contain Bathroom")
	}

	if !strings.Contains(content, "Total: $5") {
		t.Error("Content should contain Alice's total")
	}
	if !strings.Contains(content, "Total: $4") {
		t.Error("Content should contain Bob's total")
	}
}

func TestFormatNoteContentHTML_MultipleChores(t *testing.T) {
	people := []models.Person{
		{
			Name:        "Alice",
			TotalEarned: 12,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
				{Name: "Bathroom", Difficulty: 4, Earned: 4},
				{Name: "Living Room", Difficulty: 3, Earned: 3},
			},
		},
	}

	content := formatNoteContentHTML(people, false)

	if !strings.Contains(content, "• Kitchen") {
		t.Error("Content should contain Kitchen")
	}
	if !strings.Contains(content, "• Bathroom") {
		t.Error("Content should contain Bathroom")
	}
	if !strings.Contains(content, "• Living Room") {
		t.Error("Content should contain Living Room")
	}
}

func TestFormatNoteContentPlain_Default(t *testing.T) {
	people := []models.Person{
		{
			Name:        "Alice",
			TotalEarned: 5,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
			},
		},
	}

	content := formatNoteContentPlain(people, false)

	today := time.Now().Format("January 2, 2006")
	if !strings.Contains(content, today) {
		t.Errorf("Content should contain today's date (%s)", today)
	}

	if !strings.Contains(content, "•") {
		t.Error("Plain content should contain bullet symbol")
	}

	if !strings.Contains(content, "Alice") {
		t.Error("Plain content should contain person name")
	}
	if !strings.Contains(content, "Kitchen") {
		t.Error("Plain content should contain chore name")
	}

	if !strings.Contains(content, "Total: $5") {
		t.Error("Plain content should contain total")
	}

	if !strings.Contains(content, "────") {
		t.Error("Plain content should contain separator")
	}
}

func TestFormatNoteContentPlain_Verbose(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Bob",
			EffortCapacity:  15,
			TotalDifficulty: 10,
			TotalEarned:     8,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
			},
		},
	}

	content := formatNoteContentPlain(people, true)

	if !strings.Contains(content, "(Capacity: 15)") {
		t.Error("Verbose plain content should contain capacity")
	}

	if !strings.Contains(content, "Difficulty: 6") {
		t.Error("Verbose plain content should contain difficulty")
	}

	if !strings.Contains(content, "Effort: 10 / 15") {
		t.Error("Verbose plain content should contain effort")
	}
}

func TestFormatNoteContentPlain_VerboseNoCapacity(t *testing.T) {
	people := []models.Person{
		{
			Name:            "Charlie",
			EffortCapacity:  0,
			TotalDifficulty: 6,
			TotalEarned:     5,
			Chores: []models.Chore{
				{Name: "Kitchen", Difficulty: 6, Earned: 5},
			},
		},
	}

	content := formatNoteContentPlain(people, true)

	if !strings.Contains(content, "Difficulty: 6") {
		t.Error("Verbose plain content should contain difficulty")
	}

	if strings.Contains(content, "Capacity:") {
		t.Error("Should not show capacity when it's 0")
	}
	if strings.Contains(content, "Effort:") {
		t.Error("Should not show effort when capacity is 0")
	}
}

func TestFormatNoteContentHTML_EmptyChores(t *testing.T) {
	people := []models.Person{
		{
			Name:        "Alice",
			TotalEarned: 0,
			Chores:      []models.Chore{},
		},
	}

	content := formatNoteContentHTML(people, false)

	if !strings.Contains(content, "<b>Alice</b>") {
		t.Error("Content should contain person name even with no chores")
	}

	if !strings.Contains(content, "Total: $0") {
		t.Error("Content should contain total of $0")
	}
}

func TestIsSupported(t *testing.T) {
	_ = IsSupported()
}
