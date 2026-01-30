package models

type Chore struct {
	Name        string `json:"Name"`
	Difficulty  int    `json:"Difficulty"`
	Earned      int    `json:"Earned"`
	Description string `json:"Description,omitempty"`
}

type Person struct {
	Name              string  `json:"Name"`
	Contact           string  `json:"Contact,omitempty"`
	EffortCapacity    int     `json:"EffortCapacity"`
	PreAssignedChores []Chore `json:"PreAssignedChores,omitempty"`
	Chores            []Chore `json:"-"`
	TotalDifficulty   int     `json:"-"`
	TotalEarned       int     `json:"-"`
}

type Config struct {
	Chores            []Chore `json:"chores"`
	People            []Person `json:"people"`
	SMSTemplatePath   string  `json:"smsTemplatePath,omitempty"`
	NotesTemplatePath string  `json:"notesTemplatePath,omitempty"`
}
