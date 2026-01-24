package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/faradayfan/chore-distributor/internal/models"
)

func Load(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var config models.Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	for i := range config.People {
		if config.People[i].Chores == nil {
			config.People[i].Chores = []models.Chore{}
		}
	}

	return &config, nil
}
