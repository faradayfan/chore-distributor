package sms

import (
	"strings"
	"testing"

	"github.com/faradayfan/chore-distributor/internal/models"
)

func TestFormatMessage_Default(t *testing.T) {
	person := models.Person{
		Name:            "Alice",
		Contact:         "+1234567890",
		EffortCapacity:  10,
		TotalDifficulty: 6,
		TotalEarned:     5,
		Chores: []models.Chore{
			{Name: "Kitchen", Difficulty: 6, Earned: 5},
		},
	}

	message := formatMessage(person, false)

	if !strings.Contains(message, "Hi Alice!") {
		t.Error("Message should contain greeting with name")
	}
	if !strings.Contains(message, "Kitchen") {
		t.Error("Message should contain chore name")
	}
	if !strings.Contains(message, "Earns: $5") {
		t.Error("Message should contain earnings")
	}
	if !strings.Contains(message, "Total: $5") {
		t.Error("Message should contain total earned")
	}

	if strings.Contains(message, "Difficulty:") {
		t.Error("Default message should not contain difficulty")
	}
	if strings.Contains(message, "Effort:") {
		t.Error("Default message should not contain effort")
	}
}

func TestFormatMessage_Verbose(t *testing.T) {
	person := models.Person{
		Name:            "Bob",
		Contact:         "bob@icloud.com",
		EffortCapacity:  15,
		TotalDifficulty: 10,
		TotalEarned:     8,
		Chores: []models.Chore{
			{Name: "Kitchen", Difficulty: 6, Earned: 5},
			{Name: "Bathroom", Difficulty: 4, Earned: 3},
		},
	}

	message := formatMessage(person, true)

	if !strings.Contains(message, "Difficulty: 6") {
		t.Error("Verbose message should contain difficulty")
	}
	if !strings.Contains(message, "Effort: 10 / 15") {
		t.Error("Verbose message should contain effort with capacity")
	}
}

func TestFormatMessage_VerboseNoCapacity(t *testing.T) {
	person := models.Person{
		Name:            "Charlie",
		Contact:         "charlie@gmail.com",
		EffortCapacity:  0, 
		TotalDifficulty: 10,
		TotalEarned:     8,
		Chores: []models.Chore{
			{Name: "Kitchen", Difficulty: 6, Earned: 5},
		},
	}

	message := formatMessage(person, true)

	if !strings.Contains(message, "Difficulty:") {
		t.Error("Verbose message should contain difficulty")
	}
	if strings.Contains(message, "Effort:") {
		t.Error("Should not show effort line when capacity is 0")
	}
}

func TestFormatMessage_MultipleChores(t *testing.T) {
	person := models.Person{
		Name:        "Alice",
		Contact:     "+1234567890",
		TotalEarned: 12,
		Chores: []models.Chore{
			{Name: "Kitchen", Difficulty: 6, Earned: 5},
			{Name: "Bathroom", Difficulty: 4, Earned: 4},
			{Name: "Living Room", Difficulty: 3, Earned: 3},
		},
	}

	message := formatMessage(person, false)

	if !strings.Contains(message, "Kitchen") {
		t.Error("Message should contain Kitchen")
	}
	if !strings.Contains(message, "Bathroom") {
		t.Error("Message should contain Bathroom")
	}
	if !strings.Contains(message, "Living Room") {
		t.Error("Message should contain Living Room")
	}
	if !strings.Contains(message, "Total: $12") {
		t.Error("Message should contain correct total")
	}
}

func TestNewSender(t *testing.T) {
	sender := NewSender(true)
	if !sender.DryRun {
		t.Error("Sender should have DryRun set to true")
	}

	sender = NewSender(false)
	if sender.DryRun {
		t.Error("Sender should have DryRun set to false")
	}
}

func TestFormatMessage_EmailContact(t *testing.T) {
	person := models.Person{
		Name:        "Dad",
		Contact:     "dad@icloud.com",
		TotalEarned: 5,
		Chores: []models.Chore{
			{Name: "Kitchen", Difficulty: 6, Earned: 5},
		},
	}

	message := formatMessage(person, false)

	if !strings.Contains(message, "Hi Dad!") {
		t.Error("Message should contain greeting")
	}
	if !strings.Contains(message, "Kitchen") {
		t.Error("Message should contain chore")
	}
}

func TestFormatMessage_WithDescription(t *testing.T) {
	person := models.Person{
		Name:            "Alice",
		Contact:         "+1234567890",
		TotalEarned:     9,
		Chores: []models.Chore{
			{Name: "Kitchen", Difficulty: 6, Earned: 5, Description: "Clean counters, sink, and floors"},
			{Name: "Bathroom", Difficulty: 5, Earned: 4, Description: ""},
		},
	}

	message := formatMessage(person, false)

	if !strings.Contains(message, "Kitchen") {
		t.Error("Message should contain chore name")
	}
	if !strings.Contains(message, "Clean counters, sink, and floors") {
		t.Error("Message should contain description")
	}
	if !strings.Contains(message, "Bathroom") {
		t.Error("Message should contain second chore")
	}
}

func TestFormatMessage_VerboseWithDescription(t *testing.T) {
	person := models.Person{
		Name:            "Bob",
		Contact:         "bob@icloud.com",
		EffortCapacity:  15,
		TotalDifficulty: 4,
		TotalEarned:     3,
		Chores: []models.Chore{
			{Name: "Living Room", Difficulty: 4, Earned: 3, Description: "Vacuum and dust"},
		},
	}

	message := formatMessage(person, true)

	if !strings.Contains(message, "Living Room") {
		t.Error("Message should contain chore name")
	}
	if !strings.Contains(message, "Difficulty: 4") {
		t.Error("Verbose message should contain difficulty")
	}
	if !strings.Contains(message, "Vacuum and dust") {
		t.Error("Message should contain description")
	}
}
