package sms

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/faradayfan/chore-distributor/internal/models"
	"github.com/faradayfan/chore-distributor/internal/templates"
)

type Sender struct {
	DryRun       bool
	TemplatePath string
}

func NewSender(dryRun bool, templatePath string) *Sender {
	return &Sender{
		DryRun:       dryRun,
		TemplatePath: templatePath,
	}
}

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

		message, err := s.formatMessage(person, verbose)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", person.Name, err))
			continue
		}

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

func (s *Sender) formatMessage(person models.Person, verbose bool) (string, error) {
	// If a template path is provided, use it
	if s.TemplatePath != "" {
		// Check if template file exists
		if _, err := os.Stat(s.TemplatePath); err == nil {
			data := templates.BuildPersonData(person, verbose)
			return templates.LoadAndExecute(s.TemplatePath, data)
		}
		// If template path is specified but file doesn't exist, return error
		return "", fmt.Errorf("template file not found: %s", s.TemplatePath)
	}

	// Fall back to hardcoded format
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Hi %s! Here are your chores:\n\n", person.Name))

	// Add pre-assigned chores first
	for _, chore := range person.PreAssignedChores {
		if verbose {
			sb.WriteString(fmt.Sprintf("• %s (Difficulty: %d, Earns: $%d)\n",
				chore.Name, chore.Difficulty, chore.Earned))
			if chore.Description != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", chore.Description))
			}
		} else {
			sb.WriteString(fmt.Sprintf("• %s (Earns: $%d)\n",
				chore.Name, chore.Earned))
			if chore.Description != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", chore.Description))
			}
		}
	}
	// Then add distributed chores
	for _, chore := range person.Chores {
		if verbose {
			sb.WriteString(fmt.Sprintf("• %s (Difficulty: %d, Earns: $%d)\n",
				chore.Name, chore.Difficulty, chore.Earned))
			if chore.Description != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", chore.Description))
			}
		} else {
			sb.WriteString(fmt.Sprintf("• %s (Earns: $%d)\n",
				chore.Name, chore.Earned))
			if chore.Description != "" {
				sb.WriteString(fmt.Sprintf("  %s\n", chore.Description))
			}
		}
	}

	sb.WriteString(fmt.Sprintf("\nTotal: $%d", person.TotalEarned))

	if verbose && person.EffortCapacity > 0 {
		sb.WriteString(fmt.Sprintf("\nEffort: %d / %d", person.TotalDifficulty, person.EffortCapacity))
	}

	return sb.String(), nil
}

func sendViaMessages(contact, message string) error {
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

func IsSupported() bool {
	return runtime.GOOS == "darwin"
}
