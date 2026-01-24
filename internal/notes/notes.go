package notes

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/faradayfan/chore-distributor/internal/models"
)

// Writer handles writing to Apple Notes
type Writer struct {
	DryRun   bool
	NoteName string
}

// NewWriter creates a new Apple Notes writer
func NewWriter(noteName string, dryRun bool) *Writer {
	return &Writer{
		NoteName: noteName,
		DryRun:   dryRun,
	}
}

// PrependChoreList adds the chore distribution to the specified Apple Note
func (w *Writer) PrependChoreList(people []models.Person, verbose bool) error {
	if runtime.GOOS != "darwin" && !w.DryRun {
		return fmt.Errorf("Apple Notes is only supported on macOS")
	}

	content := formatNoteContentHTML(people, verbose)

	if w.DryRun {
		fmt.Printf("\n--- Would insert into note '%s' ---\n%s\n", w.NoteName, formatNoteContentPlain(people, verbose))
		return nil
	}

	if err := updateNote(w.NoteName, content); err != nil {
		return err
	}

	fmt.Printf("✓ Added chore list to note '%s'\n", w.NoteName)
	return nil
}

// formatNoteContentHTML creates HTML content for Apple Notes
func formatNoteContentHTML(people []models.Person, verbose bool) string {
	var sb strings.Builder

	// Date header
	dateStr := time.Now().Format("Monday, January 2, 2006")
	sb.WriteString(fmt.Sprintf("<div><b>%s</b></div>", dateStr))
	sb.WriteString("<div><br></div>")

	for _, person := range people {
		// Person name
		if verbose && person.EffortCapacity > 0 {
			sb.WriteString(fmt.Sprintf("<div><b>%s</b> (Capacity: %d)</div>", person.Name, person.EffortCapacity))
		} else {
			sb.WriteString(fmt.Sprintf("<div><b>%s</b></div>", person.Name))
		}

		// Chores
		for _, chore := range person.Chores {
			if verbose {
				sb.WriteString(fmt.Sprintf("<div>• %s (Difficulty: %d, Earns: $%d)</div>",
					chore.Name, chore.Difficulty, chore.Earned))
			} else {
				sb.WriteString(fmt.Sprintf("<div>• %s — $%d</div>",
					chore.Name, chore.Earned))
			}
		}

		// Totals
		if verbose && person.EffortCapacity > 0 {
			sb.WriteString(fmt.Sprintf("<div>Total: $%d | Effort: %d / %d</div>",
				person.TotalEarned, person.TotalDifficulty, person.EffortCapacity))
		} else {
			sb.WriteString(fmt.Sprintf("<div>Total: $%d</div>", person.TotalEarned))
		}
		sb.WriteString("<div><br></div>")
	}

	// Separator
	sb.WriteString("<div>─────────────────────</div>")
	sb.WriteString("<div><br></div>")

	return sb.String()
}

// formatNoteContentPlain creates plain text content for dry-run preview
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

		for _, chore := range person.Chores {
			if verbose {
				sb.WriteString(fmt.Sprintf("  • %s (Difficulty: %d, Earns: $%d)\n",
					chore.Name, chore.Difficulty, chore.Earned))
			} else {
				sb.WriteString(fmt.Sprintf("  • %s — $%d\n",
					chore.Name, chore.Earned))
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

// updateNote finds the note, removes first line, then rebuilds with: title + new content + old content
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

// IsSupported checks if Apple Notes is supported on this platform
func IsSupported() bool {
	return runtime.GOOS == "darwin"
}
