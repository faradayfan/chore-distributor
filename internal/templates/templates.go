package templates

import (
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/faradayfan/chore-distributor/internal/models"
)

// ChoreData represents a single chore for template rendering
type ChoreData struct {
	Name        string
	Difficulty  int
	Earned      float64
	Description string
}

// PersonData represents all data for a person's chore assignment
type PersonData struct {
	PersonName        string
	Contact           string
	Date              time.Time
	PreAssignedChores []ChoreData
	DistributedChores []ChoreData
	AllChores         []ChoreData
	TotalEarned       float64
	TotalDifficulty   int
	Capacity          int
	Verbose           bool
}

// BuildPersonData converts a models.Person to PersonData for template rendering
func BuildPersonData(person models.Person, verbose bool) PersonData {
	data := PersonData{
		PersonName:      person.Name,
		Contact:         person.Contact,
		Date:            time.Now(),
		TotalEarned:     float64(person.TotalEarned),
		TotalDifficulty: person.TotalDifficulty,
		Capacity:        person.EffortCapacity,
		Verbose:         verbose,
	}

	// Convert pre-assigned chores
	for _, chore := range person.PreAssignedChores {
		choreData := ChoreData{
			Name:        chore.Name,
			Difficulty:  chore.Difficulty,
			Earned:      float64(chore.Earned),
			Description: chore.Description,
		}
		data.PreAssignedChores = append(data.PreAssignedChores, choreData)
		data.AllChores = append(data.AllChores, choreData)
	}

	// Convert distributed chores
	for _, chore := range person.Chores {
		choreData := ChoreData{
			Name:        chore.Name,
			Difficulty:  chore.Difficulty,
			Earned:      float64(chore.Earned),
			Description: chore.Description,
		}
		data.DistributedChores = append(data.DistributedChores, choreData)
		data.AllChores = append(data.AllChores, choreData)
	}

	return data
}

// HelperFuncs returns the template helper functions
func HelperFuncs() template.FuncMap {
	return template.FuncMap{
		"currency": func(amount float64) string {
			return fmt.Sprintf("$%.2f", amount)
		},
		"date": func(format string, t time.Time) string {
			return t.Format(format)
		},
		"pluralize": func(count int, singular, plural string) string {
			if count == 1 {
				return singular
			}
			return plural
		},
	}
}

// LoadAndExecute loads a template from a file path and executes it with the given data
func LoadAndExecute(templatePath string, data PersonData) (string, error) {
	// Read template file
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	// Parse template with helper functions
	tmpl, err := template.New("message").Funcs(HelperFuncs()).Parse(string(templateContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	// Execute template
	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return output.String(), nil
}
