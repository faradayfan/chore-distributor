package notes

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/faradayfan/chore-distributor/internal/models"
	"github.com/faradayfan/chore-distributor/internal/templates"
)

type Writer struct {
	DryRun       bool
	NoteName     string
	TemplatePath string
}

func NewWriter(noteName string, dryRun bool, templatePath string) *Writer {
	return &Writer{
		NoteName:     noteName,
		DryRun:       dryRun,
		TemplatePath: templatePath,
	}
}

func (w *Writer) PrependChoreList(people []models.Person, verbose bool) error {
	if runtime.GOOS != "darwin" && !w.DryRun {
		return fmt.Errorf("Apple Notes is only supported on macOS")
	}

	content, plainContent, err := w.formatNoteContent(people, verbose)
	if err != nil {
		return err
	}

	if w.DryRun {
		fmt.Printf("\n--- Would insert into note '%s' ---\n%s\n", w.NoteName, plainContent)
		return nil
	}

	if err := updateNote(w.NoteName, content); err != nil {
		return err
	}

	fmt.Printf("✓ Added chore list to note '%s'\n", w.NoteName)
	return nil
}

// formatNoteContent returns both HTML content (for Notes) and plain text (for dry-run)
func (w *Writer) formatNoteContent(people []models.Person, verbose bool) (htmlContent string, plainContent string, err error) {
	// If a template path is provided, use it
	if w.TemplatePath != "" {
		// Check if template file exists
		if _, statErr := os.Stat(w.TemplatePath); statErr == nil {
			// For template-based content, we need to combine all people into one output
			htmlContent, plainContent, err = w.formatWithTemplate(people, verbose)
			return
		}
		// If template path is specified but file doesn't exist, return error
		err = fmt.Errorf("template file not found: %s", w.TemplatePath)
		return
	}

	// Fall back to hardcoded format
	htmlContent = formatNoteContentHTML(people, verbose)
	plainContent = formatNoteContentPlain(people, verbose)
	return
}

// formatWithTemplate processes all people using the template
func (w *Writer) formatWithTemplate(people []models.Person, verbose bool) (htmlContent string, plainContent string, err error) {
	var htmlBuilder, plainBuilder strings.Builder

	for _, person := range people {
		data := templates.BuildPersonData(person, verbose)
		content, templateErr := templates.LoadAndExecute(w.TemplatePath, data)
		if templateErr != nil {
			err = templateErr
			return
		}
		// Use the same content for both HTML and plain (template determines format)
		htmlBuilder.WriteString(content)
		plainBuilder.WriteString(content)
	}

	htmlContent = htmlBuilder.String()
	plainContent = plainBuilder.String()
	return
}

func formatNoteContentHTML(people []models.Person, verbose bool) string {
	var sb strings.Builder

	dateStr := time.Now().Format("Monday, January 2, 2006")
	sb.WriteString(fmt.Sprintf("<div><b>%s</b></div>", dateStr))
	sb.WriteString("<div><br></div>")

	for _, person := range people {
		if verbose && person.EffortCapacity > 0 {
			sb.WriteString(fmt.Sprintf("<div><b>%s</b> (Capacity: %d)</div>", person.Name, person.EffortCapacity))
		} else {
			sb.WriteString(fmt.Sprintf("<div><b>%s</b></div>", person.Name))
		}

		// Add pre-assigned chores first
		for _, chore := range person.PreAssignedChores {
			if verbose {
				sb.WriteString(fmt.Sprintf("<div>• %s (Difficulty: %d, Earns: $%d)</div>",
					chore.Name, chore.Difficulty, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("<div style=\"padding-left: 20px; color: #666;\">%s</div>",
						chore.Description))
				}
			} else {
				sb.WriteString(fmt.Sprintf("<div>• %s — $%d</div>",
					chore.Name, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("<div style=\"padding-left: 20px; color: #666;\">%s</div>",
						chore.Description))
				}
			}
		}
		// Then add distributed chores
		for _, chore := range person.Chores {
			if verbose {
				sb.WriteString(fmt.Sprintf("<div>• %s (Difficulty: %d, Earns: $%d)</div>",
					chore.Name, chore.Difficulty, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("<div style=\"padding-left: 20px; color: #666;\">%s</div>",
						chore.Description))
				}
			} else {
				sb.WriteString(fmt.Sprintf("<div>• %s — $%d</div>",
					chore.Name, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("<div style=\"padding-left: 20px; color: #666;\">%s</div>",
						chore.Description))
				}
			}
		}

		if verbose && person.EffortCapacity > 0 {
			sb.WriteString(fmt.Sprintf("<div>Total: $%d | Effort: %d / %d</div>",
				person.TotalEarned, person.TotalDifficulty, person.EffortCapacity))
		} else {
			sb.WriteString(fmt.Sprintf("<div>Total: $%d</div>", person.TotalEarned))
		}
		sb.WriteString("<div><br></div>")
	}

	sb.WriteString("<div>─────────────────────</div>")
	sb.WriteString("<div><br></div>")

	return sb.String()
}

func formatNoteContentPlain(people []models.Person, verbose bool) string {
	var sb strings.Builder

	dateStr := time.Now().Format("Monday, January 2, 2006")
	sb.WriteString(fmt.Sprintf("═══ %s ═══\n\n", dateStr))

	for _, person := range people {
		if verbose && person.EffortCapacity > 0 {
			sb.WriteString(fmt.Sprintf("%s (Capacity: %d)\n", person.Name, person.EffortCapacity))
		} else {
			sb.WriteString(fmt.Sprintf("%s\n", person.Name))
		}

		// Add pre-assigned chores first
		for _, chore := range person.PreAssignedChores {
			if verbose {
				sb.WriteString(fmt.Sprintf("  • %s (Difficulty: %d, Earns: $%d)\n",
					chore.Name, chore.Difficulty, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("    %s\n", chore.Description))
				}
			} else {
				sb.WriteString(fmt.Sprintf("  • %s — $%d\n",
					chore.Name, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("    %s\n", chore.Description))
				}
			}
		}
		// Then add distributed chores
		for _, chore := range person.Chores {
			if verbose {
				sb.WriteString(fmt.Sprintf("  • %s (Difficulty: %d, Earns: $%d)\n",
					chore.Name, chore.Difficulty, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("    %s\n", chore.Description))
				}
			} else {
				sb.WriteString(fmt.Sprintf("  • %s — $%d\n",
					chore.Name, chore.Earned))
				if chore.Description != "" {
					sb.WriteString(fmt.Sprintf("    %s\n", chore.Description))
				}
			}
		}

		if verbose && person.EffortCapacity > 0 {
			sb.WriteString(fmt.Sprintf("  Total: $%d | Effort: %d / %d\n\n",
				person.TotalEarned, person.TotalDifficulty, person.EffortCapacity))
		} else {
			sb.WriteString(fmt.Sprintf("  Total: $%d\n\n", person.TotalEarned))
		}
	}

	sb.WriteString("────────────────────────\n")

	return sb.String()
}

func updateNote(noteName, newContent string) error {
	escapedContent := strings.ReplaceAll(newContent, `\`, `\\`)
	escapedContent = strings.ReplaceAll(escapedContent, `"`, `\"`)
	escapedNoteName := strings.ReplaceAll(noteName, `"`, `\"`)

	script := fmt.Sprintf(`
tell application "Notes"
	set noteName to "%s"
	set newContent to "%s"
	set titleHTML to "<div>" & noteName & "</div>"
	
	-- Try to find existing note
	set noteExists to false
	try
		set targetNote to note noteName of default account
		set noteExists to true
	end try
	
	if noteExists then
		-- Get current body
		set currentBody to body of targetNote
		
		-- Remove the first line (title) by finding the first </div> and taking everything after it
		set oldContent to ""
		try
			set divEnd to offset of "</div>" in currentBody
			if divEnd > 0 then
				set oldContent to text (divEnd + 6) thru -1 of currentBody
			end if
		end try
		
		-- Rebuild: title + new chores + old content (without old title)
		set body of targetNote to titleHTML & newContent & oldContent
	else
		-- Create new note with title + content
		tell default account
			make new note with properties {body:(titleHTML & newContent)}
		end tell
	end if
end tell
`, escapedNoteName, escapedContent)

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
