package distributor

import (
	"fmt"
	"io"
	"math/rand/v2"
	"sort"

	"github.com/faradayfan/chore-distributor/internal/models"
)

// PrintOptions controls what information is displayed in the output
type PrintOptions struct {
	Verbose bool // Show difficulty and capacity information
}

// Distribute assigns chores to people based on earning balance and capacity limits
func Distribute(chores []models.Chore, people []models.Person) []models.Person {
	sortedChores := make([]models.Chore, len(chores))
	copy(sortedChores, chores)

	// Shuffle for randomization
	rand.Shuffle(len(sortedChores), func(i, j int) {
		sortedChores[i], sortedChores[j] = sortedChores[j], sortedChores[i]
	})

	// Sort by earned (highest first) with stable sort to preserve shuffle order for ties
	sort.SliceStable(sortedChores, func(i, j int) bool {
		return sortedChores[i].Earned > sortedChores[j].Earned
	})

	for _, chore := range sortedChores {
		var candidates []int
		minEarned := -1

		for i := 0; i < len(people); i++ {
			hasCapacity := people[i].EffortCapacity == 0 ||
				(people[i].TotalDifficulty+chore.Difficulty <= people[i].EffortCapacity)

			if !hasCapacity {
				continue
			}

			if minEarned == -1 || people[i].TotalEarned < minEarned {
				minEarned = people[i].TotalEarned
				candidates = []int{i}
			} else if people[i].TotalEarned == minEarned {
				candidates = append(candidates, i)
			}
		}

		if len(candidates) == 0 {
			fmt.Printf("Warning: Could not assign chore '%s' - no one has capacity\n", chore.Name)
			continue
		}

		minIndex := candidates[rand.IntN(len(candidates))]

		people[minIndex].Chores = append(people[minIndex].Chores, chore)
		people[minIndex].TotalDifficulty += chore.Difficulty
		people[minIndex].TotalEarned += chore.Earned
	}

	return people
}

// PrintDistribution outputs the chore distribution to the given writer
func PrintDistribution(w io.Writer, people []models.Person, opts PrintOptions) {
	fmt.Fprintf(w, "\n=== Chore Distribution ===\n\n")

	for _, person := range people {
		fmt.Fprintf(w, "%s", person.Name)
		if opts.Verbose && person.EffortCapacity > 0 {
			fmt.Fprintf(w, " (Effort Capacity: %d)", person.EffortCapacity)
		}
		fmt.Fprintln(w, ":")
		fmt.Fprintln(w, "  Chores:")
		for _, chore := range person.Chores {
			if opts.Verbose {
				fmt.Fprintf(w, "    - %s (Difficulty: %d, Earns: $%d)\n",
					chore.Name, chore.Difficulty, chore.Earned)
			} else {
				fmt.Fprintf(w, "    - %s (Earns: $%d)\n",
					chore.Name, chore.Earned)
			}
		}
		if opts.Verbose {
			fmt.Fprintf(w, "  Total Difficulty: %d", person.TotalDifficulty)
			if person.EffortCapacity > 0 {
				fmt.Fprintf(w, " / %d", person.EffortCapacity)
			}
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "  Total Earned: $%d\n", person.TotalEarned)
		fmt.Fprintln(w)
	}
}
