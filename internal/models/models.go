package models

// Chore represents a household chore with its properties
type Chore struct {
	Name       string `json:"Name"`
	Difficulty int    `json:"Difficulty"`
	Earned     int    `json:"Earned"`
}

// Person represents a family member who can be assigned chores
type Person struct {
	Name            string  `json:"Name"`
	Contact         string  `json:"Contact,omitempty"` // Phone number or Apple ID email for iMessage
	EffortCapacity  int     `json:"EffortCapacity"`    // 0 means no capacity limit
	Chores          []Chore `json:"-"`
	TotalDifficulty int     `json:"-"`
	TotalEarned     int     `json:"-"`
}

// Config holds the configuration loaded from JSON
type Config struct {
	Chores []Chore  `json:"chores"`
	People []Person `json:"people"`
}
