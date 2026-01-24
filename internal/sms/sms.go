package sms

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/faradayfan/chore-distributor/internal/models"
)

// Sender handles sending messages via iMessage
type Sender struct {
	DryRun bool // If true, print messages instead of sending
}

// NewSender creates a new message sender
func NewSender(dryRun bool) *Sender {
	return &Sender{DryRun: dryRun}
}

// SendChoreAssignments sends each person their chore assignments via iMessage
func (s *Sender) SendChoreAssignments(people []models.Person, verbose bool) error {
	if runtime.GOOS != "darwin" && !s.DryRun {
		return fmt.Errorf("iMessage is only supported on macOS")
	}

	var errs []string
	for _, person := range people {
		if person.Contact == "" {
			fmt.Printf("Skipping %s: no contact configured\n", person.Name)
			continue
		}

		message := formatMessage(person, verbose)

		if s.DryRun {
			fmt.Printf("\n--- Would send to %s (%s) ---\n%s\n", person.Name, person.Contact, message)
			continue
		}

		if err := sendViaMessages(person.Contact, message); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", person.Name, err))
			continue
		}

		fmt.Printf("✓ Sent chores to %s (%s)\n", person.Name, person.Contact)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send some messages: %s", strings.Join(errs, "; "))
	}

	return nil
}

// formatMessage creates the iMessage content for a person
func formatMessage(person models.Person, verbose bool) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Hi %s! Here are your chores:\n\n", person.Name))

	for _, chore := range person.Chores {
		if verbose {
			sb.WriteString(fmt.Sprintf("• %s (Difficulty: %d, Earns: $%d)\n",
				chore.Name, chore.Difficulty, chore.Earned))
		} else {
			sb.WriteString(fmt.Sprintf("• %s (Earns: $%d)\n",
				chore.Name, chore.Earned))
		}
	}

	sb.WriteString(fmt.Sprintf("\nTotal: $%d", person.TotalEarned))

	if verbose && person.EffortCapacity > 0 {
		sb.WriteString(fmt.Sprintf("\nEffort: %d / %d", person.TotalDifficulty, person.EffortCapacity))
	}

	return sb.String()
}

// sendViaMessages sends a message using macOS Messages app via AppleScript
// The contact can be a phone number (e.g., "+15551234567") or an Apple ID email (e.g., "user@icloud.com")
func sendViaMessages(contact, message string) error {
	// Escape special characters for AppleScript
	escapedMessage := strings.ReplaceAll(message, `\`, `\\`)
	escapedMessage = strings.ReplaceAll(escapedMessage, `"`, `\"`)

	script := fmt.Sprintf(`
tell application "Messages"
	set targetService to 1st service whose service type = iMessage
	set targetBuddy to buddy "%s" of targetService
	send "%s" to targetBuddy
end tell
`, contact, escapedMessage)

	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("AppleScript error: %v, output: %s", err, string(output))
	}

	return nil
}

// IsSupported checks if iMessage sending is supported on this platform
func IsSupported() bool {
	return runtime.GOOS == "darwin"
}
