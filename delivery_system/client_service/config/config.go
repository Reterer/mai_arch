package config

import (
	"encoding/json"
	"os"
)

type Api struct {
}

type Repository struct {
}

type Service struct {
}

type Config struct {
	Port string `json:"port"`

	Api        Api        `json:"api"`
	Service    Service    `json:"service"`
	Repository Repository `json:"repository"`
}

func ParseConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
