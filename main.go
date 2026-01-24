package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"sort"
)

type Chore struct {
	Name       string
	Difficulty int
	Earned     int
}

type Person struct {
	Name            string
	EffortCapacity  int // 0 means no capacity limit
	Chores          []Chore
	TotalDifficulty int
	TotalEarned     int
}

type Config struct {
	Chores []Chore  `json:"chores"`
	People []Person `json:"people"`
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for i := range config.People {
		if config.People[i].Chores == nil {
			config.People[i].Chores = []Chore{}
		}
	}

	return &config, nil
}

func distributeChores(chores []Chore, people []Person) []Person {
	sortedChores := make([]Chore, len(chores))
	copy(sortedChores, chores)

	rand.Shuffle(len(sortedChores), func(i, j int) {
		sortedChores[i], sortedChores[j] = sortedChores[j], sortedChores[i]
	})

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

func printDistribution(people []Person) {
	fmt.Printf("\n=== Chore Distribution (Balanced by earned) ===\n\n")

	for _, person := range people {
		fmt.Printf("%s", person.Name)
		if person.EffortCapacity > 0 {
			fmt.Printf(" (Effort Capacity: %d)", person.EffortCapacity)
		}
		fmt.Println(":")
		fmt.Println("  Chores:")
		for _, chore := range person.Chores {
			fmt.Printf("    - %s (Difficulty: %d, Earns: $%d)\n",
				chore.Name, chore.Difficulty, chore.Earned)
		}
		fmt.Printf("  Total Difficulty: %d", person.TotalDifficulty)
		if person.EffortCapacity > 0 {
			fmt.Printf(" / %d", person.EffortCapacity)
		}
		fmt.Println()
		fmt.Printf("  Total Earned: $%d\n", person.TotalEarned)
		fmt.Println()
	}
}

func main() {
	configPath := flag.String("config", "chores_config.json", "Path to the JSON configuration file")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	config.People = distributeChores(config.Chores, config.People)
	printDistribution(config.People)
}
